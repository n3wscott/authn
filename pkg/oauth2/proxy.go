package oauth2

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/exec"
	"path"
	"strings"
)

// note: this needs oauth2-proxy in the base image.

func StartOAuth2Proxy(ctx context.Context) error {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		return err
	}

	cmd := env.OAuth2Proxy
	args := []string{"-config=" + path.Join(env.ConfigPath, OAuth2ProxyCfg)}

	fmt.Println(strings.Join(args, " "))

	ex := exec.CommandContext(ctx, cmd, args...)
	ex.Stderr = os.Stderr
	ex.Stdout = os.Stdout

	if err := ex.Run(); err != nil {
		return err
	}
	return nil
}