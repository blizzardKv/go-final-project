package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	PwdHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

func pwdHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func signToken(password string) (string, error) {
	hash := pwdHash(password)
	c := claims{
		PwdHash: hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(hash))
}

func validateToken(tokenStr, password string) bool {
	hash := pwdHash(password)
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		return []byte(hash), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	c, ok := token.Claims.(*claims)
	return ok && c.PwdHash == hash
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			var tokenStr string
			if cookie, err := r.Cookie("token"); err == nil {
				tokenStr = cookie.Value
			}
			if !validateToken(tokenStr, pass) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}
	token, err := signToken(pass)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{"token": token})
}
