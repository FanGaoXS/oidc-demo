package user

import (
	"encoding/json"
	"net/http"

	"oidc-demo/http-client/app2/storage"
)

var (
	s *storage.Storage
)

func init() {
	s = storage.New()
}

func Users(w http.ResponseWriter, r *http.Request) {
	users := s.AllUser()

	bytes, _ := json.Marshal(users)
	w.Write(bytes)
}
