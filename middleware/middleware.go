package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/nikhilpratapgit/storex/database/dbHelper"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
)

type userContextKeyType struct{}

var userContextKey = userContextKeyType{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "missing token")
			return
		}
		token, parseErr := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if parseErr != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}
		sessionID := claimValues["sessionID"].(string)
		fmt.Println(sessionID)
		userID, err := dbHelper.ValidateSession(sessionID)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, err, "invalid user")
			return
		}
		//userDetail, err := dbHelper.GetUserByID(userID)
		//if err != nil {
		//	utils.RespondError(w, http.StatusUnauthorized, err, "user not found")
		//	return
		//}

		user := &models.UserCtx{
			UserID:    userID,
			SessionID: sessionID,
			Role:      claimValues["role"].(string),
		}
		fmt.Println(user.Role)
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
func UserContext(r *http.Request) *models.UserCtx {
	user, _ := r.Context().Value(userContextKey).(*models.UserCtx)
	return user
}
func RoleMiddleware(Roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx := UserContext(r)
			if userCtx == nil {
				utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
				return
			}
			UserRole := userCtx.Role

			for _, role := range Roles {
				if UserRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.RespondError(w, http.StatusForbidden, nil, "not-authorised")
		})
	}
}
