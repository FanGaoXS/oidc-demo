package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"oidc-demo/http-client/app1/auth"
	"oidc-demo/http-client/app1/environment"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	OidcIssuer     string
	ClientID       string
	ResourceMyBook string
)

func init() {
	env := environment.Load()
	OidcIssuer = env.OidcIssuer
	ClientID = env.ClientID
	ResourceMyBook = env.ResourceMyBook
}

func ReadMyRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, err := auth.GetFromCookie(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("get token from request failed: %s", err), http.StatusUnauthorized)
		return
	}
	accessToken := token.AccessToken
	tokenType := token.TokenType
	rawIDToken := token.IdToken

	provider, err := oidc.NewProvider(ctx, OidcIssuer)
	if err != nil {
		http.Error(w, fmt.Sprintf("init oidc provider failed: %s", err), http.StatusInternalServerError)
		return
	}
	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: ClientID})

	// verify rawIDToken to get idToken, if the idToken is expired, then refresh the token with refresh_token
	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil && !strings.Contains(err.Error(), "oidc: token is expired") {
		http.Error(w, fmt.Sprintf("verify id token failed: %s", err), http.StatusUnauthorized)
		return
	}
	if err != nil && strings.Contains(err.Error(), "oidc: token is expired") {
		// use refresh_token to refresh token
		config := auth.Oauth2Config(provider)
		ts := config.TokenSource(ctx, &oauth2.Token{RefreshToken: token.RefreshToken})
		newToken, err := ts.Token()
		if err != nil {
			http.Error(w, fmt.Sprintf("refresh token failed: %s", err), http.StatusInternalServerError)
			return
		}
		auth.SetIntoCookie(w, newToken) // set new token(contains access_token, refresh_token, id_token...) into cookie
		accessToken = newToken.AccessToken
		tokenType = newToken.TokenType
		rawIDToken = newToken.Extra("id_token").(string)
		idToken, _ = idTokenVerifier.Verify(ctx, rawIDToken)
	}
	if err = idToken.VerifyAccessToken(accessToken); err != nil {
		http.Error(w, fmt.Sprintf("id_token does not match access_token"), http.StatusUnauthorized)
		return
	} // check if id_token matches access_token

	client := http.DefaultClient
	req, _ := http.NewRequest("GET", ResourceMyBook, nil)
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req) // do get request with Authorization (access_token)
	if err != nil {
		http.Error(w, fmt.Sprintf("get my book from resource failed: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	w.Write(bytes)
}
