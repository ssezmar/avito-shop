package auth

import (
    "net/http"
    "strings"
    "time"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
)

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func GenerateJWT(email, secret string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func JWTMiddleware(secret string) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, "Bearer ")
            if len(parts) != 2 {
                http.Error(w, "Invalid token format", http.StatusUnauthorized)
                return
            }

            tokenStr := parts[1]
            claims := &Claims{}

            token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })

            if err != nil {
                if err == jwt.ErrSignatureInvalid {
                    http.Error(w, "Invalid token signature", http.StatusUnauthorized)
                    return
                }
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            if !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            context.Set(r, "user", claims.Email)
            next.ServeHTTP(w, r)
        })
    }
}