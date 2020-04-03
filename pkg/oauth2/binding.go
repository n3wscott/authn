package oauth2

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coreos/go-oidc"
	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAppPort = "8181"
)

func init() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		_ = os.Setenv("APP_PORT", defaultAppPort)
	}
}

type envConfig struct {
	Port        string `envconfig:"PORT" default:"8080"` // TODO: need to fix up the config with this port if it is not 8080
	ConfigPath  string `envconfig:"CONFIG_ROOT" default:"/etc/proxy-config/"`
	OAuth2Proxy string `envconfig:"OAUTH_PROXY_PATH" default:"/bin/oauth2_proxy"`
}

const (
	OAuth2ProxyCfg = "oauth2_proxy.cfg"
	IssuerName     = "oicd_issuer"
	ClientIDName   = "oicd_client_id"
)

func ReadIssuer(root string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(root, IssuerName))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadClientID(root string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(root, ClientIDName))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewOrDie(ctx context.Context, handler http.Handler) http.Handler {
	provider, err := New(ctx, handler)
	if err != nil {
		panic(err)
	}

	return provider
}

func New(ctx context.Context, handler http.Handler) (http.Handler, error) {
	provider, err := NewProvider(ctx, handler)
	if err != nil {
		return handler, err
	}

	// Start the oauth proxy:
	go func() {
		if err := StartOAuth2Proxy(ctx); err != nil {
			panic(err)
		}
	}()

	return provider, nil
}

func NewProvider(ctx context.Context, handler http.Handler) (http.Handler, error) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		return handler, err
	}

	issuer, err := ReadIssuer(env.ConfigPath)
	if err != nil {
		return handler, err
	}

	clientID, err := ReadClientID(env.ConfigPath)
	if err != nil {
		return handler, err
	}
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return handler, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("/", Middleware(ctx, handler, verifier))

	return mux, nil
}
