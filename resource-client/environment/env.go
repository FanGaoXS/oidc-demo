package environment

import (
	"log"

	"github.com/joho/godotenv"
)

type env struct {
	OidcIssuer    string
	ListenAddress string
}

func Load() *env {
	envMap, err := godotenv.Read(".env.resource")
	if err != nil {
		log.Fatalln("load env from .env failed: ", err)
	}

	return &env{
		OidcIssuer:    envMap["oidc_issuer"],
		ListenAddress: envMap["listen_address"],
	}
}
