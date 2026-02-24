package middleware

import (
	"context"
	"fmt"
	"net/http"

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
			//http.Error(w, "missing token", http.StatusUnauthorized)
			utils.RespondError(w, http.StatusUnauthorized, nil, "missing token")
			return
		}
		userID, err := dbHelper.VaidateSession(tokenStr)
		fmt.Println(userID)
		if err != nil {
			//http.Error(w, "invalid user", http.StatusUnauthorized)
			utils.RespondError(w, http.StatusUnauthorized, err, "invalid user")
			return
		}
		userDetail, err := dbHelper.GetUserByID(userID.String())
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, err, "user not found")
			return
		}

		user := &models.UserCtx{
			UserID:    userID,
			SessionID: tokenStr,
			Role:      userDetail.Role,
		}
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
