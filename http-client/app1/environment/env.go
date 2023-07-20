package environment

import (
	"log"

	"github.com/joho/godotenv"
)

type env struct {
	OidcIssuer   string
	ClientID     string
	ClientSecret string
	RedirectURL  string

	ListenAddress  string
	ResourceMyBook string
}

func Load() *env {
	envMap, err := godotenv.Read(".env.gitee")
	if err != nil {
		log.Fatalln("load env from .env failed: ", err)
	}

	return &env{
		OidcIssuer:     envMap["oidc_issuer"],
		ClientID:       envMap["client_id"],
		ClientSecret:   envMap["client_secret"],
		RedirectURL:    envMap["redirect_url"],
		ListenAddress:  envMap["listen_address"],
		ResourceMyBook: envMap["resource_my_book"],
	}
}
