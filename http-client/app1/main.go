package main

import (
	"log"
	"net/http"

	"oidc-demo/http-client/app1/api"
	"oidc-demo/http-client/app1/auth"
	"oidc-demo/http-client/app1/environment"
	"oidc-demo/http-client/app1/user"
)

func main() {
	env := environment.Load()

	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/logout", auth.Logout)
	http.HandleFunc("/callback", auth.LoginCallback)
	http.HandleFunc("/read", api.ReadMyRepo)
	http.HandleFunc("/users", user.Users)

	log.Printf("%s listening on %s", env.ClientID, env.ListenAddress)
	http.ListenAndServe(env.ListenAddress, nil)
}
