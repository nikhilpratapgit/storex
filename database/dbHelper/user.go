package dbHelper

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/models"
	"github.com/nikhilpratapgit/storex/utils"
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

func CreateUser(name, email, userRole, userType, phoneNumber, password string) (string, error) {
	SQL := `INSERT INTO users (name,email,role,type,phone_no,password)
			VALUES ($1,LOWER(TRIM($2)),$3,$4,$5,$6)
			RETURNING id`
	var userID string
	err := database.Store.Get(&userID, SQL, name, email, userRole, userType, phoneNumber, password)
	return userID, err
}

func GetUserByEmail(email, password string) (string, error) {
	SQL := `SELECT id, password
			FROM USERS 
			WHERE email=$1 AND 
			archived_at IS NULL
			`
	var User models.UserAuth
	err := database.Store.Get(&User, SQL, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("no user exist")
		}
		return "", err
	}
	if err := utils.CheckPassword(User.Password, password); err != nil {
		return "", err
	}
	return User.ID, nil
}
func CreateUserSession(id string) (string, error) {
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
func VaidateSession(sessionID string) (uuid.UUID, error) {
	SQL := `SELECT user_id 
			FROM user_session
			WHERE id =$1 AND
			archived_at IS NULL
			`
	var userID uuid.UUID
	err := database.Store.Get(&userID, SQL, sessionID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func DeleteSessionByToken(sessionID string) error {
	SQL := `UPDATE user_session
			SET archived_at= NOW()
			WHERE id=$1
			AND archiced_at IS NULL 
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
		return errors.New("Invalid Session")
	}
	return nil
}
func GetUserByID(userID string) (*models.User, error) {
	SQL := `SELECT name ,email ,role ,type
			FROM users 
			WHERE id=$1
			AND archived_at IS NULL `
	var user models.User
	err := database.Store.Get(&user, SQL, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CreateAsset(brand, model, serialNumber, assetType, status, owner string, warrantyStart, warrantyEnd time.Time) error {
	SQL := `INSERT INTO assets (brand ,model ,serial_number ,type ,status ,owner ,warranty_start ,warranty_end)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := database.Store.Exec(SQL, brand, model, serialNumber, assetType, status, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil
}
func ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr string) ([]models.AssetInfo, error) {
	SQL := `SELECT brand ,model ,type ,serial_number ,status ,owner ,created_at
			from assets
			WHERE archived_at IS NULL 
			AND (
			    $1= '' or brand LIKE'%'||$1||'%'
			)
			AND(
			    $2 ='' or model LIKE'%'||$2||'%'
			)
			AND(
			    $3 ='' or serial_number LIKE'%'||$3||'%'
			)
			AND (
			    $4='' or type::text LIKE'%'||$4||'%'
			)
			AND(
			    $5='' or status::text LIKE '%'||$5||'%'
			)
			AND(
			    $6=''or owner::text LIKE '%'||$6||'%'
			)
			order by created_at
			`
	assets := make([]models.AssetInfo, 0)

	err := database.Store.Select(&assets, SQL, brandStr, modelStr, serialNumberStr, typeStr, statusStr, ownerStr)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
