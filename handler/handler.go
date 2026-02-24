package handler

import (
	"log"
	"net/http"

	"github.com/nikhilpratapgit/storex/database/dbHelper"
	"github.com/nikhilpratapgit/storex/middleware"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUser models.RegisterUser

	if parseErr := utils.ParseBody(r.Body, &registerUser); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}

	exist, existErr := dbHelper.IsUserExist(registerUser.Email)
	if existErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existErr, "failed to check user existence")
		return
	}
	if exist {
		utils.RespondError(w, http.StatusBadRequest, nil, "user exist")
		return
	}
	hashPassword, err := utils.HashPassword(registerUser.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed while hashing password")
		return
	}

	_, saveErr := dbHelper.CreateUser(registerUser.Name, registerUser.Email, registerUser.Role, registerUser.Type, registerUser.PhoneNumber, hashPassword)
	if saveErr != nil {
		utils.RespondError(w, http.StatusNotFound, saveErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{Message: "user created successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var UserLogin models.LoginUser

	parseErr := utils.ParseBody(r.Body, &UserLogin)
	if parseErr != nil {
		utils.RespondError(w, http.StatusNotFound, parseErr, "invalid request body")
		return
	}
	userID, err := dbHelper.GetUserByEmail(UserLogin.Email, UserLogin.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "failed to find user")
		return
	}
	sessionID, err := dbHelper.CreateUserSession(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		SessionID string `json:"sessionID"`
	}{SessionID: sessionID})
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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var assetRequest models.Asset
	//userCtx := middleware.UserContext(r)

	if parseErr := utils.ParseBody(r.Body, &assetRequest); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}
	err := dbHelper.CreateAsset(assetRequest.Brand,
		assetRequest.Model,
		assetRequest.SerialNumber,
		assetRequest.AssetType,
		assetRequest.Status,
		assetRequest.Owner,
		assetRequest.WarrantyStart,
		assetRequest.WarrantyEnd)
	if err != nil {
		log.Println(err)
		utils.RespondError(w, http.StatusBadRequest, err, "failed to create asset")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, assetRequest)
}

func ShowAssets(w http.ResponseWriter, r *http.Request) {
	//var assetInfo models.AssetInfo
	//allAsset := make([]models.AssetInfo, 0)
	typeStr := r.URL.Query().Get("type")
	statusStr := r.URL.Query().Get("status")
	ownerStr := r.URL.Query().Get("owner")
	brandStr := r.URL.Query().Get("brand")
	modelStr := r.URL.Query().Get("model")
	serialNumberStr := r.URL.Query().Get("serialNumber")

	Assets, err := dbHelper.ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr)
	if err != nil {
		log.Println(err)
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Assets []models.AssetInfo `json:"assets"`
	}{Assets: Assets})
}
