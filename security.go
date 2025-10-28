package main

import (
	"context" // Necesario para el contexto de la petición
	"fmt"
	"log"
	"net/http"
	"strings" // Necesario para strings.HasPrefix
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ⬇️ DEFINICIONES NECESARIAS PARA EL CONTEXTO
type ContextKey string

const ContextKeyUserID ContextKey = "userID"

type Claims struct {
	jwt.RegisteredClaims

	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

func GenerateToken(UserID int, Role string, SecretKey string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	claims := Claims{
		UserID: UserID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("error al firmar el token JWT: %w", err)
	}
	return tokenString, nil
}

func AuthMiddleware(SecretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[7:]

			tokenParsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(SecretKey), nil
			})

			if err != nil {
				// ⬇️ CORRECCIÓN: Se agrega el log para ver el error.
				log.Printf("Error de verificación JWT: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !tokenParsed.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := tokenParsed.Claims.(*Claims)
			if !ok {
				http.Error(w, "Error interno de validación", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(r *http.Request) (int, error) {
	// 1. Obtener el valor del contexto usando la clave que definimos.
	// r.Context() devuelve el contexto adjunto a la petición.
	userIDValue := r.Context().Value(ContextKeyUserID)

	// 2. Verificar si el valor existe.
	if userIDValue == nil {
		// Esto solo debería pasar si la ruta no está protegida o si el middleware falló.
		return 0, fmt.Errorf("UserID no encontrado en el contexto")
	}

	// 3. Hacer la aserción de tipo y verificar que sea un entero (int).
	userID, ok := userIDValue.(int)
	if !ok {
		// Esto indica un error de programación donde se guardó un tipo incorrecto.
		return 0, fmt.Errorf("valor de UserID en el contexto no es un entero")
	}

	// 4. Devolver el UserID.
	return userID, nil
}
