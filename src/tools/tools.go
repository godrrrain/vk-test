package tools

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
)

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(
			"method", r.Method,
			"path", r.URL.EscapedPath(),
		)
		next(w, r)
	}
}

func RequestAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(
			"method", r.Method,
			"path", r.URL.EscapedPath(),
		)

		username, password, ok := r.BasicAuth()
		if ok {
			cred := username + password
			credHash := sha256.Sum256([]byte(cred))

			//username = "abc", password = "123"
			rightHash := [32]byte{108, 161, 61, 82, 202, 112, 200, 131, 224, 240, 187, 16, 30, 66, 90,
				137, 232, 98, 77, 229, 29, 178, 210, 57, 37, 147, 175, 106, 132, 17, 128, 144}
			usernameMatch := (subtle.ConstantTimeCompare(credHash[:], rightHash[:]) == 1)

			// if username == "abc" && password == "123" {
			// 	next(w, r)
			// 	return
			// }
			if usernameMatch {
				next(w, r)
				return
			}
		}
		log.Println("Unauthorized")

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
