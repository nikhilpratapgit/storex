package dbHelper

import (
	"errors"
	"time"

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

func CreateUser(name, email, userRole, userType, phoneNumber, password string) (string, error) {
	SQL := `INSERT INTO users (name,email,role,type,phone_no,password)
			VALUES ($1,LOWER(TRIM($2)),$3,$4,$5,$6)
			RETURNING id`
	var userID string
	err := database.Store.Get(&userID, SQL, name, email, userRole, userType, phoneNumber, password)
	return userID, err
}

func GetUserByEmail(email string) (models.User, error) {
	SQL := `SELECT id, password,email
			FROM USERS 
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
func VaidateSession(sessionID string) (string, error) {
	SQL := `SELECT user_id 
			FROM user_session
			WHERE id =$1 AND
			archived_at IS NULL
			`
	var userID string
	err := database.Store.Get(&userID, SQL, sessionID)
	if err != nil {
		return "", err
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

//func GetUserByID(userID string) (*models.User, error) {
//	SQL := `SELECT name ,email ,role ,type
//			FROM users
//			WHERE id=$1
//			AND archived_at IS NULL `
//	var user models.User
//	err := database.Store.Get(&user, SQL, userID)
//	if err != nil {
//		return nil, err
//	}
//	return &user, nil
//}

func CreateAsset(tx *sqlx.Tx, assetRequest models.Asset) (string, error) {
	SQL := `INSERT INTO assets (brand, model, serial_number ,type ,status ,owner ,warranty_start ,warranty_end)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			RETURNING id
			`
	var assetID string
	args := []interface{}{
		assetRequest.Brand,
		assetRequest.Model,
		assetRequest.SerialNumber,
		assetRequest.AssetType,
		assetRequest.Status,
		assetRequest.Owner,
		assetRequest.WarrantyStart,
		assetRequest.WarrantyEnd,
	}
	err := tx.Get(&assetID, SQL, args...)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	return assetID, nil
}

// FETCH ASSETS
func ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr string, limit, offset int) (models.DashboardData, error) {
	SQL := `SELECT brand ,model ,type ,serial_number ,status ,owner ,created_at
			FROM assets
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
			ORDER BY created_at
			LIMIT $7 OFFSET $8
			`

	assets := make([]models.AssetInfo, 0)
	var summary models.DashboardSummary
	// make different func

	Sql := `SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'available') AS available,
			COUNT(*) FILTER (WHERE status = 'assigned') AS assigned,
			COUNT(*) FILTER (WHERE status = 'for_repair') AS waiting_for_repair,
			COUNT(*) FILTER (WHERE status = 'in_service') AS in_service,
			COUNT(*) FILTER (WHERE status = 'damaged') AS damaged
		FROM assets
		WHERE archived_at IS NULL`

	var res models.DashboardData
	DashboardErr := database.Store.Get(&summary, Sql)
	if DashboardErr != nil {
		return res, DashboardErr
	}

	err := database.Store.Select(&assets, SQL, brandStr, modelStr, serialNumberStr, typeStr, statusStr, ownerStr, limit, offset)
	if err != nil {
		return res, err
	}
	return models.DashboardData{
		Summary: summary,
		Assets:  assets,
	}, nil
}
func CreateLaptop(tx *sqlx.Tx, assetID string, assetRequest models.LaptopSpecs) error {
	SQL := `INSERT INTO laptop (asset_id,processor,ram,storage,operating_system,charger,device_password)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`
	args := []interface{}{
		assetID,
		assetRequest.Processor,
		assetRequest.Ram,
		assetRequest.Storage,
		assetRequest.OperatingSystem,
		assetRequest.Charger,
		assetRequest.DevicePassword,
	}
	_, err := tx.Exec(SQL, args...)
	if err != nil {
		return err
	}
	return nil
}
func CreateKeyboard(tx *sqlx.Tx, assetID string, assetRequest models.KeyboardSpecs) error {
	SQL := `INSERT INTO keyboard(asset_id,layout,connectivity)
			VALUES($1,$2,$3)`
	args := []interface{}{
		assetID,
		assetRequest.Layout,
		assetRequest.Connectivity,
	}
	_, err := tx.Exec(SQL, args...)
	if err != nil {
		return err
	}
	return nil
}
func CreateMouse(tx *sqlx.Tx, assetID string, assetRequest models.MouseSpecs) error {
	SQL := `INSERT INTO mouse (asset_id,dpi,connectivity)
			VALUES($1,$2,$3)`
	args := []interface{}{
		assetID,
		assetRequest.Connectivity,
		assetRequest.Dpi,
	}
	_, err := tx.Exec(SQL, args...)
	if err != nil {
		return err
	}
	return nil
}
func CreateMobile(tx *sqlx.Tx, assetID string, assetRequest models.MobileSpecs) error {
	SQL := `INSERT INTO mobile (asset_id,operating_system,ram,storage,charger,device_password)
			VALUES ($1,$2,$3,$4,$5,$6)`

	args := []interface{}{
		assetID,
		assetRequest.OperatingSystem,
		assetRequest.Ram,
		assetRequest.Storage,
		assetRequest.Charger,
		assetRequest.DevicePassword,
	}
	_, err := tx.Exec(SQL, args...)
	if err != nil {
		return err
	}
	return nil
}
func AssignedAssets(id, assignedById, assignedTo string) error {
	SQL := `UPDATE assets
			SET assigned_by_id =$1,
			    assigned_to=$2,
			    assigned_on=NOW(),
			    status='assigned',
			    updated_at=NOW()
			WHERE id=$3
			AND archived_at IS NULL 
			    `
	_, err := database.Store.Exec(SQL, assignedById, assignedTo, id)
	if err != nil {
		return err
	}
	return nil
}
func ServiceAssets(assetID string, serviceStart, serviceEnd, returnedOn time.Time) error {
	SQL := `UPDATE assets
			SET service_start=$1
			service_end=$2
			returned_on=$3
			WHERE id=$4
			    AND status !='assigned'
			AND archived_at IS NULL 
			`
	_, err := database.Store.Exec(SQL, serviceStart, serviceEnd, returnedOn, assetID)
	if err != nil {
		return err
	}
	return nil
}
func DeleteAsset(archivedBy, assetID string) error {
	SQL := `UPDATE assets
			SET archived_at=NOW()
			archived_by=$1
			WHERE
			assetID=$2
			AND archived_at IS NULL
			`
	_, err := database.Store.Exec(SQL, archivedBy, assetID)
	if err != nil {
		return err
	}
	return nil
}
func GetAssetInfo(userID, assetStatus string) ([]models.AssetInfoRequest, error) {
	SQL := `
		SELECT id, brand, model, status, asset_type
		FROM assets
		WHERE assigned_to=$1
		AND ($2 = '' OR status::TEXT=$2)
	`
	assetDetails := make([]models.AssetInfoRequest, 0)
	err := database.Store.Select(&assetDetails, SQL, userID, assetStatus)
	return assetDetails, err
}
func GetUserInfo(name, role, employment, assetStatus string) ([]models.UserInfoRequest, error) {
	SQL := `
		SELECT id, name, email, phone_number, role, employment, created_at
		FROM users
		WHERE ($1 = '' OR name LIKE '%' || $1 || '%')
		AND ($2 = '' OR role::TEXT=$2)
		AND ($3 = '' OR employment::TEXT=$3)
	`
	users := make([]models.UserInfoRequest, 0)
	err := database.Store.Select(&users, SQL, name, role, employment)
	if err != nil {
		return users, err
	}

	filteredUsers := make([]models.UserInfoRequest, 0)
	for _, user := range users {
		assetDetails, err := GetAssetInfo(user.ID, assetStatus)
		if err != nil {
			return users, err
		}
		if assetStatus != "available" && len(assetDetails) == 0 {
			continue
		}
		user.AssetDetails = assetDetails
		filteredUsers = append(filteredUsers, user)
	}
	return filteredUsers, err

}
func FetchUser(userID string) (models.UserInfoRequest, error) {
	SQL := `
		SELECT id, name, email, phone_no, role, type, created_at
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
		return user, nil
	}
	user.AssetDetails = assetDetails
	//filterUser = append(filterUser.AssetDetails, user)
	//}
	return user, err

}
func GetAsset(userID string) ([]models.AssetInfoRequest, error) {
	SQL := `
		SELECT id, brand, model, status, type
		FROM assets
		WHERE assigned_to=$1
		AND archived_at IS NULL 
	`
	assetDetails := make([]models.AssetInfoRequest, 0)
	err := database.Store.Select(&assetDetails, SQL, userID)
	return assetDetails, err
}

func UpdateAsset(tx *sqlx.Tx, assetID, brand, model, serialNo, assetType, owner string, warrantyStart, warrantyEnd time.Time) error {
	query := `UPDATE assets
            set brand = $2, model = $3, serial_no = $4, type=$5,owner=$6,warranty_start = $7,warranty_end=$8, updated_at =now()
            where id= $1 and archived_at is null `
	_, err := tx.Exec(query, assetID, brand, model, serialNo, assetType, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil

}
func UpdateLaptop(tx *sqlx.Tx, assetID string, laptop *models.LaptopSpecs) error {
	query := `
    UPDATE laptop
    SET
        processor = $2,
        ram = $3,
        storage = $4,
        os = $5,
        charger = $6,
        password = $7
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query,
		assetID,
		laptop.Processor,
		laptop.Ram,
		laptop.Storage,
		laptop.OperatingSystem,
		laptop.Charger,
		laptop.DevicePassword,
	)

	return err
}
func UpdateMouse(tx *sqlx.Tx, assetID string, mouse *models.MouseSpecs) error {
	query := `
    UPDATE mouse
    SET
        dpi = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query, assetID, mouse.Dpi, mouse.Connectivity)
	return err
}
func UpdateKeyboard(tx *sqlx.Tx, assetID string, keyboard *models.KeyboardSpecs) error {
	query := `
    UPDATE keyboard
    SET
        layout = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(query, assetID, keyboard.Layout, keyboard.Connectivity)
	return err
}
func UpdateMobile(tx *sqlx.Tx, assetID string, mobile *models.MobileSpecs) error {
	query := `
    UPDATE mobile
    SET
        os = $2,
        ram = $3,
        storage = $4,
        charger = $5,
        password = $6
    WHERE asset_id = $1
    `

	_, err := tx.Exec(
		query,
		assetID,
		mobile.OperatingSystem,
		mobile.Ram,
		mobile.Storage,
		mobile.Charger,
		mobile.DevicePassword,
	)

	return err
}
