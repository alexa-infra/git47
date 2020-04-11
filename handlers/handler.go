package handlers

import (
	"log"
	"net/http"
)

func makeHandler(fn func(*Env, http.ResponseWriter, *http.Request) error, env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(env, w, r)
		if err != nil {
			status := http.StatusInternalServerError
			switch e := err.(type) {
			case Error:
				status = e.Status()
			}
			log.Printf("HTTP %d - %s", status, err)
			http.Error(w, err.Error(), status)
		}
	}
}
