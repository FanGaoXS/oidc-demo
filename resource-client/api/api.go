package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/http"
	"strings"

	"oidc-demo/resource-client/environment"
	"oidc-demo/resource-client/storage"
	"oidc-demo/resource-client/userinfo"
)

var (
	OidcIssuer string

	s *storage.Storage
)

func init() {
	env := environment.Load()
	OidcIssuer = env.OidcIssuer

	s = storage.New()
}

func AddRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.FormValue("name")
	if name = strings.TrimSpace(name); name == "" {
		http.Error(w, fmt.Sprintf("invalid repo name: empty repo name"), http.StatusBadRequest)
		return
	}

	token, err := tokenFromHeader(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("get token from header failed: %s", err), http.StatusUnauthorized)
		return
	}
	ui, err := Userinfo(ctx, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("get userinfo failed: %s", err), http.StatusUnauthorized)
		return
	}

	ok := s.AddRepo(name, ui.Subject, ui.Audience)
	fmt.Fprintf(w, "%t", ok)
}

func MyRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := tokenFromHeader(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("get token from header failed: %s", err), http.StatusUnauthorized)
		return
	}
	ui, err := Userinfo(ctx, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("get userinfo failed: %s", err), http.StatusUnauthorized)
		return
	}

	books := s.GetRepoBySubject(ui.Subject)
	bytes, _ := json.Marshal(books)
	w.Write(bytes)
}

func AllRepo(w http.ResponseWriter, r *http.Request) {
	// for admin

	books := s.AllRepo()
	bytes, _ := json.Marshal(books)
	w.Write(bytes)
}

func tokenFromHeader(header http.Header) (token string, err error) {
	token = header.Get("Authorization")
	splits := strings.SplitN(token, " ", 2)
	if len(splits) < 2 {
		return "", fmt.Errorf("invalid authorization: empty authorization")
	}

	typ := splits[0]
	token = splits[1]
	if typ != "Bearer" && typ != "bearer" {
		return "", fmt.Errorf("invalid authorization type: %s", typ)
	}

	return token, nil
}

func Userinfo(ctx context.Context, accessToken string) (*userinfo.Userinfo, error) {
	provider, err := oidc.NewProvider(ctx, OidcIssuer)
	if err != nil {
		return nil, fmt.Errorf("init oidc provider failed: %s", err)
	}

	// get userinfo from userinfo endpoint with access_token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	providerUi, err := provider.UserInfo(ctx, ts)
	if err != nil {
		return nil, fmt.Errorf("get userinfo from %s failed: %s", provider.UserInfoEndpoint(), err)
	}

	// parse provider userinfo into internal userinfo
	var ui userinfo.Userinfo
	if err = providerUi.Claims(&ui); err != nil {
		return nil, fmt.Errorf("parse userinfo from %s failed: %s", provider.UserInfoEndpoint(), err)
	}

	return &ui, nil
}
