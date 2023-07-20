package main

import (
	"log"
	"net/http"

	"oidc-demo/resource-client/api"
	"oidc-demo/resource-client/environment"
)

func main() {
	env := environment.Load()

	http.HandleFunc("/add-repo", api.AddRepo)
	http.HandleFunc("/my-repo", api.MyRepo)
	http.HandleFunc("/all-repo", api.AllRepo)

	log.Printf("%s listening on %s", "resouce", env.ListenAddress)
	http.ListenAndServe(env.ListenAddress, nil)
}
