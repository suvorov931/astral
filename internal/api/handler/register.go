package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type user struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

func Register(logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var u user

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			logger.Error("Register: cannot decode request", zap.Error(err))
		}

		fmt.Println(u)
	}
}
