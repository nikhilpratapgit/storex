package dbHelper

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/models"
)

func IsUserExist(email string) (bool, error) {
	SQL := `SELECT count(*)>0
			FROM users
			WHERE email=TRIM(LOWER($1))
			AND archived_at IS NULL
			`
	var exists bool
	err := database.Store.Get(&exists, SQL, email)
	return exists, err
}

func CreateUser(tx *sqlx.Tx, name, email, userType, phoneNumber, password string) (string, error) {
	SQL := `INSERT INTO users (name,email,type,phone_number,password)
			VALUES ($1,LOWER(TRIM($2)),$3,$4,$5)
			RETURNING id`
	var userID string
	err := tx.Get(&userID, SQL, name, email, userType, phoneNumber, password)
	return userID, err
}

func GetUserByEmail(email string) (models.User, error) {
	SQL := `SELECT id, password,email,role
			FROM users 
			WHERE email=TRIM(LOWER($1)) 
			AND 
			archived_at IS NULL
			`
	var User models.User
	err := database.Store.Get(&User, SQL, email)

	if err != nil {
		return User, err
	}

	return User, nil
}
func CreateRegisterUserSession(tx *sqlx.Tx, id string) (string, error) {
	SQL := `INSERT INTO user_session (user_id)
			VALUES ($1) RETURNING id
			`
	var sessionID string
	err := tx.Get(&sessionID, SQL, id)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
func CreateLoginUserSession(id string) (string, error) {
	SQL := `INSERT INTO user_session (user_id)
			VALUES ($1) RETURNING id
			`
	var sessionID string
	err := database.Store.Get(&sessionID, SQL, id)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
func ValidateSession(sessionID string) (string, error) {
	SQL := `SELECT user_id 
			FROM user_session
			WHERE id =$1 AND
			archived_at IS NULL
			`
	var userID string
	err := database.Store.Get(&userID, SQL, sessionID)
	return userID, err
}

func DeleteSessionByToken(sessionID string) error {
	SQL := `UPDATE user_session
			SET archived_at= NOW()
			WHERE id=$1
			AND archived_at IS NULL 
			`
	result, err := database.Store.Exec(SQL, sessionID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("invalid session")
	}
	return nil
}

func GetUserInfo(name, role, userType, assetStatus string, limit, offset int) ([]models.UserInfoRequest, error) {
	SQL := `SELECT u.id,
                  u.name,
                  u.email,
                  u.phone_number,
                  u.role,
                  u.type,
                  u.created_at,
                  a.id   AS asset_id,
                  a.brand,
                  a.model,
                  a.status,
                  a.type AS asset_type
           FROM users u
                  LEFT JOIN assets a
                           ON a.assigned_to = u.id
                          AND a.archived_at IS NULL
                          AND ($4 = '' OR a.status::TEXT = $4)
                  WHERE ($1 = '' OR u.name ILIKE '%' || $1 || '%')
                  AND ($2 = '' OR u.role::TEXT = $2)
                  AND ($3 = '' OR u.type::TEXT = $3)
                  AND u.archived_at IS NULL
                  LIMIT $5 OFFSET $6
	`
	userAssets := make([]models.UserAssetDetail, 0)
	err := database.Store.Select(&userAssets, SQL, name, role, userType, assetStatus, limit, offset)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*models.UserInfoRequest)

	for _, userAssetInfo := range userAssets {
		if _, ok := userMap[userAssetInfo.ID]; !ok {
			userMap[userAssetInfo.ID] = &models.UserInfoRequest{
				ID:           userAssetInfo.ID,
				Name:         userAssetInfo.Name,
				Email:        userAssetInfo.Email,
				PhoneNumber:  userAssetInfo.PhoneNumber,
				Role:         userAssetInfo.Role,
				Type:         userAssetInfo.Type,
				CreatedAt:    userAssetInfo.CreatedAt,
				AssetDetails: []models.AssetInfoRequest{},
			}
		}
		if userAssetInfo.AssetID != nil {
			asset := models.AssetInfoRequest{
				ID:     *userAssetInfo.AssetID,
				Brand:  *userAssetInfo.Brand,
				Model:  *userAssetInfo.Model,
				Status: *userAssetInfo.Status,
				Type:   *userAssetInfo.AssetType,
			}
			userMap[userAssetInfo.ID].AssetDetails = append(userMap[userAssetInfo.ID].AssetDetails, asset)
		}
	}

	users := make([]models.UserInfoRequest, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, *user)
	}
	return users, nil
}

func FetchUser(userID string) (models.UserInfoRequest, error) {
	SQL := `
		SELECT id, name, email, phone_number, role, type, created_at
		FROM users
		WHERE archived_at IS NULL 
		AND id=$1
	`
	var user models.UserInfoRequest
	err := database.Store.Get(&user, SQL, userID)
	if err != nil {
		return user, err
	}

	//var filterUser models.UserInfoRequest
	assetDetails, err := GetAsset(user.ID)
	if err != nil {
		return user, err
	}
	if len(assetDetails) == 0 {
		user.AssetDetails = assetDetails
		return user, nil
	}
	user.AssetDetails = assetDetails
	return user, err

}

func ExtractUserRole(userID string) (string, error) {
	SQL := `SELECT role
			FROM users
			WHERE id=$1 
			AND archived_at IS NULL `
	var role string
	err := database.Store.Get(&role, SQL, userID)
	if err != nil {
		return "", err
	}
	return role, nil
}
func AssignedRole(userID, role string) error {
	SQL := `UPDATE users
			SET role =$1
			WHERE id=$2
			AND archived_at IS NULL `
	_, err := database.Store.Exec(SQL, role, userID)
	if err != nil {
		return err
	}
	return nil
}
func DeleteUser(userID string) error {
	SQL := `UPDATE users
			SET archived_at=NOW()
			WHERE id=$1
			AND archived_at IS NULL 
			`
	_, err := database.Store.Exec(SQL, userID)
	if err != nil {
		return err
	}
	return nil
}
func GetUserIDByEmail(email string) (string, error) {
	SQL := `SELECT id FROM users WHERE email=$1 AND archived_at IS NULL`

	var userID string
	err := database.Store.Get(&userID, SQL, email)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func CompleteUserProfile(userID, phoneNumber, password string) error {
	SQL := `UPDATE users
			SET 
			    phone_number=$1,
			    password=$2
			WHERE id=$3
			AND archived_At IS NULL`
	_, err := database.Store.Exec(SQL, phoneNumber, password, userID)
	if err != nil {
		return err
	}
	return nil
}
