package handler

import (
	"encoding/json"
	"net/http"

	"astral/internal/api"
)

func decodeBody(w http.ResponseWriter, r *http.Request, user *api.User) error {
	r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(user); err != nil {
		return err
	}

	return nil
}
