package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/database/dbHelper"
	"github.com/nikhilpratapgit/storex/middleware"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
)

var validate = validator.New()

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUser models.RegisterUser
	var userID string
	var sessionID string
	if parseErr := utils.ParseBody(r.Body, &registerUser); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}

	validateErr := validate.Struct(&registerUser)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
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

	TxErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(registerUser.Name, registerUser.Email, registerUser.Role, registerUser.Type, registerUser.PhoneNumber, hashPassword)
		if saveErr != nil {
			utils.RespondError(w, http.StatusNotFound, saveErr, "failed to create user")
			return err
		}
		//SESSION
		sessionID, err = dbHelper.CreateUserSession(userID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user session")
			return err
		}
		return nil
	})
	if TxErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
	}

	token, err := utils.GenerateJWT(userID, sessionID, registerUser.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to generate token")
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
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
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
	sessionID, err := dbHelper.CreateUserSession(user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user session")
		return
	}

	token, err := utils.GenerateJWT(user.ID, sessionID, user.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to generate token")
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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var assetRequest models.Asset

	if parseErr := utils.ParseBody(r.Body, &assetRequest); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}
	validateErr := validate.Struct(&assetRequest)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	if assetRequest.WarrantyEnd.Before(assetRequest.WarrantyStart) {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid warranty range")
		return
	}

	Txerr := database.Tx(func(tx *sqlx.Tx) error {
		assetID, err := dbHelper.CreateAsset(tx, assetRequest)
		if err != nil {
			return fmt.Errorf("failed to create asset: %w", err)
		}

		switch assetRequest.AssetType {
		case "laptop":
			if err := dbHelper.CreateLaptop(tx, assetID, assetRequest.Laptop); err != nil {
				return fmt.Errorf("failed to create laptop: %w", err)
			}
		case "keyboard":
			if err := dbHelper.CreateKeyboard(tx, assetID, assetRequest.Keyboard); err != nil {
				return fmt.Errorf("failed to create keyboard: %w", err)
			}
		case "mouse":
			if err := dbHelper.CreateMouse(tx, assetID, assetRequest.Mouse); err != nil {
				return fmt.Errorf("failed to create mouse: %w", err)
			}
		case "mobile":
			if err := dbHelper.CreateMobile(tx, assetID, assetRequest.Mobile); err != nil {
				return fmt.Errorf("failed to create mobile: %w", err)
			}
		default:
			return fmt.Errorf("unsupported asset type: %s", assetRequest.AssetType)
		}

		return nil
	})

	if Txerr != nil {
		utils.RespondError(w, http.StatusInternalServerError, Txerr, "failed to create asset")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, "asset created successfully")
}

func ShowAssets(w http.ResponseWriter, r *http.Request) {
	//var DashboardInfo models.DashboardData
	//allAsset := make([]models.AssetInfo, 0)
	typeStr := r.URL.Query().Get("type")
	statusStr := r.URL.Query().Get("status")
	ownerStr := r.URL.Query().Get("owner")
	brandStr := r.URL.Query().Get("brand")
	modelStr := r.URL.Query().Get("model")
	serialNumberStr := r.URL.Query().Get("serialNumber")

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

	Assets, err := dbHelper.ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr, limit, offset)
	if err != nil {
		//log.Println(err)
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Assets models.DashboardData `json:"assets"`
	}{Assets: Assets})
}
func AssignedAssets(w http.ResponseWriter, r *http.Request) {
	var assignedAsset models.AssignedAsset
	assetID := chi.URLParam(r, "id")

	if parseErr := utils.ParseBody(r.Body, &assignedAsset); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed parsing body")
		return
	}
	validateErr := validate.Struct(&assignedAsset)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	err := dbHelper.AssignedAssets(assetID, userID, assignedAsset.AssignedTo)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to assigned assets")
	}
	utils.RespondJSON(w, http.StatusOK, "successfully assigned")
}
func ServiceAssets(w http.ResponseWriter, r *http.Request) {
	var serviceAsset models.ServiceAsset
	assetID := chi.URLParam(r, "id")

	if err := utils.ParseBody(r.Body, &serviceAsset); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	validateErr := validate.Struct(&serviceAsset)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	err := dbHelper.ServiceAssets(assetID, serviceAsset.ServiceStart, serviceAsset.ServiceEnd, serviceAsset.ReturnedOn)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "asset service failed")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset service complete")
}
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	var deleteAsset models.DeleteAsset
	assetID := chi.URLParam(r, "id")

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	if err := utils.ParseBody(r.Body, &deleteAsset); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	validateErr := validate.Struct(&deleteAsset)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	err := dbHelper.DeleteAsset(userID, assetID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to delete asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset deleted successfully")
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	role := query.Get("role")
	employment := query.Get("employment")
	assetStatus := query.Get("status")

	userDetails, err := dbHelper.GetUserInfo(name, role, employment, assetStatus)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users": userDetails,
	})
}
func FetchUser(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query()
	userID := chi.URLParam(r, "id")

	userDetails, err := dbHelper.FetchUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users": userDetails,
	})
}

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	assetId := chi.URLParam(r, "id")
	if assetId == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid assetID")
		return
	}
	var body models.UpdateAssetRequest
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid body")
		return
	}
	validateErr := validate.Struct(&body)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	if body.WarrantyEnd.Before(body.WarrantyStart) {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid warranty range")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbHelper.UpdateAsset(tx, assetId, body.Brand, body.Model, body.SerialNo, body.Type, body.Owner, body.WarrantyStart, body.WarrantyEnd)
		if err != nil {
			return err
		}
		switch body.Type {

		case "laptop":
			if body.Laptop == nil {
				return fmt.Errorf("laptop details required")
			}
			return dbHelper.UpdateLaptop(tx, assetId, body.Laptop)

		case "mouse":
			if body.Mouse == nil {
				return fmt.Errorf("mouse details required")
			}
			return dbHelper.UpdateMouse(tx, assetId, body.Mouse)
		case "keyboard":
			if body.Keyboard == nil {
				return fmt.Errorf("keyboard details required")
			}
			return dbHelper.UpdateKeyboard(tx, assetId, body.Keyboard)
		case "mobile":
			if body.Mobile == nil {
				return fmt.Errorf("mobile details required")
			}
			return dbHelper.UpdateMobile(tx, assetId, body.Mobile)

		default:
			return fmt.Errorf("unsupported asset type")
		}
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusBadRequest, txErr, "fail to update asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "asset updated",
	})
}
