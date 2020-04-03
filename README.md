# WhoAmI

Demo using OAuth.

Fill out `./config/oauth2_proxy.cfg`, `./config/oicd_client_id` and `./config/oicd_issuer` and make a secret:

```shell
kubectl create secret generic whoami-proxy-config --from-file=./config/oauth2_proxy.cfg --from-file=./config/oicd_client_id --from-file=./config/oicd_issuer
```

