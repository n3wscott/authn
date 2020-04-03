// Package oauth2 holds utilities for bootstrapping an OAuth2 proxy
// from the metadata injected by a OAuth2 secret. Within an app,
// users can write:
//    middleware, err := oauth2.New(ctx, handler)
// or
//    http.ListenAndServe(":8181", oauth2.NewOrDie(ctx, handler)))
// This is modeled after the Bindings pattern.
package oauth2
