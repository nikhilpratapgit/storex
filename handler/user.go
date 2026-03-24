package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/database/dbHelper"
	"github.com/nikhilpratapgit/storex/middleware"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUser models.RegisterUser

	var userID string
	var sessionID string
	var createUserErr error

	if parseErr := utils.ParseBody(r.Body, &registerUser); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}

	validateErr := validate.Struct(&registerUser)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "failed to validate body")
		return
	}

	isUserExist, existErr := dbHelper.IsUserExist(registerUser.Email)
	if existErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existErr, "failed to check user existence")
		return
	}
	if isUserExist {
		utils.RespondError(w, http.StatusInternalServerError, nil, "user already exist")
		return
	}
	hashPassword, err := utils.HashPassword(registerUser.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed while hashing password")
		return
	}

	TransactionErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, createUserErr = dbHelper.CreateUser(tx, registerUser.Name, registerUser.Email, registerUser.Type, registerUser.PhoneNumber, hashPassword)
		if createUserErr != nil {
			utils.RespondError(w, http.StatusNotFound, createUserErr, "failed to create user")
			return createUserErr
		}

		sessionID, err = dbHelper.CreateRegisterUserSession(tx, userID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user session")
			return err
		}
		return nil
	})
	if TransactionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}

	userRole, err := dbHelper.ExtractUserRole(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to extract role of user")
	}
	//todo: will do this later
	//fixme: will fix it later
	token, err := utils.GenerateJWT(userID, sessionID, userRole)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{Token: token})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var UserLogin models.LoginUser

	parseErr := utils.ParseBody(r.Body, &UserLogin)
	if parseErr != nil {
		utils.RespondError(w, http.StatusNotFound, parseErr, "invalid request body")
		return
	}
	validateErr := validate.Struct(&UserLogin)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "failed to validate body")
		return
	}
	user, err := dbHelper.GetUserByEmail(UserLogin.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "failed to find user")
		return
	}

	if err := utils.CheckPassword(user.Password, UserLogin.Password); err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	sessionID, err := dbHelper.CreateLoginUserSession(user.ID)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user session")
		return
	}

	token, err := utils.GenerateJWT(user.ID, sessionID, user.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "some internal server error")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{Token: token})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	sessionID := userCtx.SessionID

	if err := dbHelper.DeleteSessionByToken(sessionID); err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid user")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "user logout successfully"})
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	role := query.Get("role")
	userType := query.Get("type")
	assetStatus := query.Get("status")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}
	offset := (page - 1) * limit

	userDetails, err := dbHelper.GetUserInfo(name, role, userType, assetStatus, limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users": userDetails,
	})
}
func FetchUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need userID")
		return
	}

	userDetails, err := dbHelper.FetchUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users": userDetails,
	})
}

func AssignedRole(w http.ResponseWriter, r *http.Request) {
	var role models.AssignedRole
	userID := chi.URLParam(r, "id")

	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need userID")
		return
	}

	if err := utils.ParseBody(r.Body, &role); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	validateErr := validate.Struct(&role)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}
	err := dbHelper.AssignedRole(userID, role.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to assigned role")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "assigned role successfully")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "required userID")
		return
	}

	TxErr := database.Tx(func(tx *sqlx.Tx) error {
		assetErr := dbHelper.UnassignedAssets(userID)
		if assetErr != nil {
			return assetErr
		}
		sessionErr := dbHelper.DeleteSession(userID)
		if sessionErr != nil {
			return sessionErr
		}
		userErr := dbHelper.DeleteUser(userID)
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if TxErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, TxErr, "failed to delete user")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "user deleted successfully")
}

func CompleteUserProfile(w http.ResponseWriter, r *http.Request) {
	var userProfile models.UserProfile

	ctx := r.Context()

	if err := utils.ParseBody(r.Body, &userProfile); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	validateErr := validate.Struct(&userProfile)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	hashPassword, err := utils.HashPassword(userProfile.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed while hashing password")
		return
	}

	token, err := utils.FirebaseAuth.VerifyIDToken(ctx, userProfile.Token)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid firebase token")
		return
	}

	email, _ := token.Claims["email"].(string)

	userID, err := dbHelper.GetUserIDByEmail(email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get userid")
	}

	completeProfileErr := dbHelper.CompleteUserProfile(userID, userProfile.PhoneNumber, hashPassword)
	if completeProfileErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to complete user profile")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "user profile completed")
}
