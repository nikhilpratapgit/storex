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

	txErr := database.Tx(func(tx *sqlx.Tx) error {
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

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create asset")
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

	assets, err := dbHelper.ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr, limit, offset)
	if err != nil {
		//log.Println(err)
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch assets")
		return
	}
	assetCount, err := dbHelper.AssetDashboard()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch dashboard data")
		return
	}
	assetData := models.AssetDashboardData{
		AssetCount: assetCount,
		Assets:     assets,
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Assets models.AssetDashboardData `json:"assets"`
	}{Assets: assetData})
}
func AssignedAssets(w http.ResponseWriter, r *http.Request) {
	var assignedAsset models.AssignedAsset
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need assetID")
		return
	}

	if parseErr := utils.ParseBody(r.Body, &assignedAsset); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed parsing body")
		return
	}
	validateErr := validate.Struct(&assignedAsset)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	if !dbHelper.CheckStatus(assetID) {
		utils.RespondError(w, http.StatusBadRequest, nil, "already asset assigned")
		return
	}

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	err := dbHelper.AssignedAssets(assetID, userID, assignedAsset.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to assigned assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "successfully assigned")
}
func ServiceAssets(w http.ResponseWriter, r *http.Request) {
	//var serviceAsset models.ServiceAsset

	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need assetID")
		return
	}
	err := dbHelper.ServiceAssets(assetID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "asset service failed")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset service complete")
}
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	var deleteAsset models.DeleteAsset
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need assetID")
		return
	}

	if err := utils.ParseBody(r.Body, &deleteAsset); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	validateErr := validate.Struct(&deleteAsset)
	if validateErr != nil {
		utils.RespondError(w, http.StatusBadRequest, validateErr, "fail to validate body")
		return
	}

	err := dbHelper.DeleteAsset(deleteAsset.ArchivedBy, assetID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to delete asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset deleted successfully")
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
		err := dbHelper.UpdateAsset(tx, assetId, body.Brand, body.Model, body.SerialNumber, body.Type, body.Owner, body.WarrantyStart, body.WarrantyEnd)
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

func UnassignedAllAsset(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "userID required")
		return
	}
	err := dbHelper.UnassignedAssets(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to unassigned assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset successfully unassigned")
}

func UpdateAssetStatus(w http.ResponseWriter, r *http.Request) {
	//var assetID models.MarkDamage
	assetID := chi.URLParam(r, "id")

	if assetID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "id cannot be empty")
		return
	}
	err := dbHelper.UpdateAssetStatus(assetID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "updated asset status successfully")

}

func ServiceComplete(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "need assetID")
		return
	}
	err := dbHelper.ServiceComplete(assetID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "asset service failed")
		return
	}
	utils.RespondJSON(w, http.StatusOK, "asset service successfully")
}
