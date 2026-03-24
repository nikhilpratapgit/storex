package handler

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/jmoiron/sqlx"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/database/dbHelper"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
)

func FirebaseRegister(w http.ResponseWriter, r *http.Request) {

	var body models.RegisterUser

	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid body request")
		return
	}

	ctx := context.Background()

	params := (&auth.UserToCreate{}).
		Email(body.Email).
		Password(body.Password).
		DisplayName(body.Name)

	userRecord, err := utils.FirebaseAuth.CreateUser(ctx, params)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user in firebase")
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "password hashing failed")
		return
	}
	firebaseUID := userRecord.UID
	err = dbHelper.CreateUserFirebase(body, hashedPassword)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create firebase user")
		delErr := utils.FirebaseAuth.DeleteUser(ctx, firebaseUID)
		if delErr != nil {
			log.Println("failed to delete firebase user:", delErr)
		}
		return
	}
	claims := map[string]interface{}{
		"role": "employee",
	}

	err = utils.FirebaseAuth.SetCustomUserClaims(ctx, userRecord.UID, claims)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to set claims")
		return
	}

	utils.RespondJSON(w, http.StatusOK, "user register successfully")
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body models.GoogleLoginRequest

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid body")
		return
	}

	if err := validate.Struct(&body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed")
		return
	}
	token, err := utils.FirebaseAuth.VerifyIDToken(ctx, body.IdToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid firebase token")
		return
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	if email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email not found in token")
		return
	}

	isUserExist, err := dbHelper.IsUserExist(email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to check user")
		return
	}
	var isProfileCompleted bool = false
	var userID string
	if !isUserExist {

		txErr := database.Tx(func(tx *sqlx.Tx) error {

			id, err := dbHelper.CreateUser(
				tx,
				name,
				email,
				"full_time",
				"",
				"",
			)

			if err != nil {
				return err
			}

			userID = id
			return nil
		})

		if txErr != nil {
			utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
			return
		}

	} else {
		isProfileCompleted = true
		userID, err = dbHelper.GetUserIDByEmail(email)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch user")
			return
		}
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":            "google login successful",
		"isProfileCompleted": isProfileCompleted,
		"userID":             userID,
		"token":              body.IdToken,
	})
}
