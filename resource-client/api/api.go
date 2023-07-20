package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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
	name := r.FormValue("name")
	if name = strings.TrimSpace(name); name == "" {
		http.Error(w, fmt.Sprintf("invalid repo name: empty repo name"), http.StatusBadRequest)
		return
	}

	typ, accessToken, err := tokenFromHeader(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("get token from header failed: %s", err), http.StatusUnauthorized)
		return
	}
	ui, err := getUserInfo(typ, accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("get userinfo failed: %s", err), http.StatusUnauthorized)
		return
	}

	ok := s.AddRepo(name, ui.Subject)
	fmt.Fprintf(w, "%t", ok)
}

func MyRepo(w http.ResponseWriter, r *http.Request) {
	typ, accessToken, err := tokenFromHeader(r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("get token from header failed: %s", err), http.StatusUnauthorized)
		return
	}
	ui, err := getUserInfo(typ, accessToken)
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

func tokenFromHeader(header http.Header) (typ string, token string, err error) {
	token = header.Get("Authorization")
	splits := strings.SplitN(token, " ", 2)
	if len(splits) < 2 {
		return "", "", fmt.Errorf("invalid authorization: empty authorization")
	}

	typ = splits[0]
	token = splits[1]
	if typ != "Bearer" && typ != "bearer" {
		return "", "", fmt.Errorf("invalid authorization type: %s", typ)
	}

	return typ, token, nil
}

func getUserInfo(typ, accessToken string) (*userinfo.Userinfo, error) {
	client := http.DefaultClient

	// get the configurations from {issuer}/.well-known/openid-configuration
	u, _ := url.Parse(OidcIssuer)
	u.Path = filepath.Join(u.Path, "/.well-known/openid-configuration")
	res, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("list openid-configuration from %s failed: %s", u.String(), err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("list openid-configuration from %s failed: %s", u.String(), msg)
	}

	var configurations map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&configurations); err != nil {
		return nil, fmt.Errorf("parse configurations from %s failed: %s", u.String(), err)
	}

	// get the userinfo from {issuer}/{userinfo_endpoint}
	userinfoEndpoint := configurations["userinfo_endpoint"].(string)
	req, _ := http.NewRequest("GET", userinfoEndpoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", typ, accessToken))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get userinfo from %s failed: %s", u.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("get userinfo from %s failed: %s", u.String(), msg)
	}

	var ui userinfo.Userinfo
	if err = json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return nil, fmt.Errorf("parse userinfo from %s failed: %s", userinfoEndpoint, err)
	}
	return &ui, nil
}
