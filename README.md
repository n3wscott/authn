你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# authn

_work in progress_

The goal of this project is to enable easy integration with OAuth2 for
applications that are targeting Kubernetes and Knative.

There is an explicit assumption that the binary for
[`OAuth2-Proxy`](https://github.com/oauth2-proxy/oauth2-proxy) is in a base
layer of the container your app will run in. The easiest way to do this is to
use [`ko`](https://github.com/google/ko) with a `.ko.yaml` config like:

```yaml
baseImageOverrides:
  go.path.of/your.app/: quay.io/oauth2-proxy/oauth2-proxy
```

# OAuth2

Targeting the Knative runtime contact means the OAuth2 Proxy must run on `$PORT`
and your application is going to run on `$APP_PORT`. `Authn` will default
`$APP_PORT` to `8181` if not set.

The resulting application will look like this:

```
inbound http --> [:PORT (oauth2_proxy via authn)] --> [:$APP_PORT your custom app]
```

Only authenticated requests will reach `$APP_PORT`.

### Setup

1. Fill in [`oauth2_proxy.cfg`](./config/secrets/oauth2_proxy.cfg) with the
   correct settings.
1. Fill in [`oidc_client_id`](./config/secrets/oidc_client_id) and
   [`oidc_issuer`](./config/secrets/oidc_issuer) based on the provider selected.
1. Make a secret from these files, like:

   ```shell
   kubectl create secret generic whoami-proxy-config --from-file=./config/secrets/oauth2_proxy.cfg --from-file=./config/secrets/oidc_client_id --from-file=./config/secrets/oidc_issuer
   ```

1. Confirm the base image contains `quay.io/oauth2-proxy/oauth2-proxy` as
   mentioned above.
1. Deploy your application, here is an example for the `whoami` app:

```shell
ko apply -f config/whoami.yaml
```

### Additional Settings

Please do not use `$PORT`. This is reserved for the proxy.

`$APP_PORT` - this is the port your app should run on.

If you need to change where the secret is mounted, set env var:`CONFIG_ROOT`, it
defaults to "/etc/proxy-config/"`.

If you need to change the OAuth2 Proxy binary, set `OAUTH_PROXY_PATH`, it
defaults from "/bin/oauth2_proxy".

The internal parts of the secret mounted are expected to be: `oauth2_proxy.cfg`,
`oidc_issuer`, `oidc_client_id`.
