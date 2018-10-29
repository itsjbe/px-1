package handlers

import (
	"encoding/json"
	"fmt"
	"km/phoenix/src/phoenix/torii/user"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		format := "%s - - [%s] \"%s %s %s\" %s\n"
		fmt.Printf(format, r.RemoteAddr, time.Now().Format(time.RFC1123),
			r.Method, r.URL.Path, r.Proto, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}
