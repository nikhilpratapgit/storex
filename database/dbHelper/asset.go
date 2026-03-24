package dbHelper

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nikhilpratapgit/storex/database"
	"github.com/nikhilpratapgit/storex/models"
)

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
func ShowAssets(typeStr, statusStr, ownerStr, brandStr, modelStr, serialNumberStr string, limit, offset int) ([]models.AssetInfo, error) {
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

	err := database.Store.Select(&assets, SQL, brandStr, modelStr, serialNumberStr, typeStr, statusStr, ownerStr, limit, offset)
	return assets, err
}
func AssetDashboard() (models.AssetsDashboard, error) {
	var assetCount models.AssetsDashboard

	SQL := `SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'available') AS available,
			COUNT(*) FILTER (WHERE status = 'assigned') AS assigned,
			COUNT(*) FILTER (WHERE status = 'waiting_for_repair') AS waiting_for_repair,
			COUNT(*) FILTER (WHERE status = 'in_service') AS in_service,
			COUNT(*) FILTER (WHERE status = 'damaged') AS damaged
		FROM assets
		WHERE archived_at IS NULL`

	//var res models.DashboardData
	DashboardErr := database.Store.Get(&assetCount, SQL)
	if DashboardErr != nil {
		return assetCount, DashboardErr
	}
	return assetCount, nil

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
	return err
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
	return err
}
func CreateMouse(tx *sqlx.Tx, assetID string, assetRequest models.MouseSpecs) error {
	SQL := `INSERT INTO mouse (asset_id,dpi,connectivity)
			VALUES($1,$2,$3)`
	args := []interface{}{
		assetID,
		assetRequest.Dpi,
		assetRequest.Connectivity,
	}
	_, err := tx.Exec(SQL, args...)
	return err
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
	return err
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
	return err
}
func ServiceAssets(assetID string) error {
	SQL := `UPDATE assets
			SET service_start=NOW(),
			    status='in_service'
			WHERE id=$1
			        AND status ='waiting_for_repair'
			AND archived_at IS NULL 
			`
	_, err := database.Store.Exec(SQL, assetID)
	return err
}
func DeleteAsset(archivedBy, assetID string) error {
	SQL := `UPDATE assets
			SET archived_at=NOW(),
			archived_by=$1
			WHERE
			id=$2
			AND archived_at IS NULL
			`
	_, err := database.Store.Exec(SQL, archivedBy, assetID)
	return err
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
	SQL := `UPDATE assets
            set brand = $2, model = $3, serial_number = $4, type=$5,owner=$6,warranty_start = $7,warranty_end=$8, updated_at =now()
            where id= $1 and archived_at is null `
	_, err := tx.Exec(SQL, assetID, brand, model, serialNo, assetType, owner, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil

}
func UpdateLaptop(tx *sqlx.Tx, assetID string, laptop *models.LaptopSpecs) error {
	SQL := `
    UPDATE laptop
    SET
        processor = $2,
        ram = $3,
        storage = $4,
        operating_system = $5,
        charger = $6,
        device_password = $7
    WHERE asset_id = $1
    `

	_, err := tx.Exec(SQL,
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
	SQL := `
    UPDATE mouse
    SET
        dpi = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(SQL, assetID, mouse.Dpi, mouse.Connectivity)
	return err
}
func UpdateKeyboard(tx *sqlx.Tx, assetID string, keyboard *models.KeyboardSpecs) error {
	SQL := `
    UPDATE keyboard
    SET
        layout = $2,
        connectivity = $3
    WHERE asset_id = $1
    `

	_, err := tx.Exec(SQL, assetID, keyboard.Layout, keyboard.Connectivity)
	return err
}
func UpdateMobile(tx *sqlx.Tx, assetID string, mobile *models.MobileSpecs) error {
	SQL := `
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
		SQL,
		assetID,
		mobile.OperatingSystem,
		mobile.Ram,
		mobile.Storage,
		mobile.Charger,
		mobile.DevicePassword,
	)

	return err
}

func CreateUserFirebase(user models.RegisterUser, hashedPassword string) error {

	SQL := `
	INSERT INTO users (name,email,type,phone_number,password)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	var userID uuid.UUID

	err := database.Store.Get(
		&userID,
		SQL,
		user.Name,
		user.Email,
		user.Type,
		user.PhoneNumber,
		hashedPassword,
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckStatus(assestId string) bool {
	// for checking whether assigned or not for differnt usage
	SQL := `SELECT status FROM assets WHERE id = $1 AND archived_at IS NULL`
	var status string
	err := database.Store.Get(&status, SQL, assestId)
	if err != nil {
		return false
	}
	if status == "available" {
		return true
	}
	return false
}

func UnassignedAssets(userID string) error {
	SQL := `UPDATE assets
            SET 
                assigned_by_id =null,
                assigned_to=null,
            	assigned_on=null,
            	updated_at=NOW(),
            	returned_on= NOW(),
            	status='available'
				WHERE assigned_to=$1
				AND archived_at IS NULL
                `
	_, err := database.Store.Exec(SQL, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSession(userID string) error {
	SQL := `UPDATE user_session
		  SET archived_at=NOW()
		  WHERE user_id=$1
		  AND archived_at IS NULL 
		  `
	_, err := database.Store.Exec(SQL, userID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAssetStatus(assetID string) error {
	SQL := `UPDATE assets
SET
    assigned_by_id = NULL,
    assigned_to = NULL,
    assigned_on = NULL,
    updated_at = NOW(),
    returned_on = CASE
        WHEN status = 'assigned' THEN NOW()
        ELSE returned_on
    END,
    status = CASE
        WHEN status IN ('available', 'assigned') THEN 'damaged'
        WHEN status = 'damaged' THEN 'waiting_for_repair'
    	ELSE status
    END
WHERE 
    id = $1
    AND archived_at IS NULL;
`
	_, err := database.Store.Exec(SQL, assetID)
	if err != nil {
		return err
	}
	return nil
}

func ServiceComplete(assetID string) error {
	SQL := `UPDATE assets
		  	SET service_end=NOW(),
		  	    updated_at=NOW(),
		  	    status='available'
			WHERE id=$1
			AND archived_at IS NULL 
			AND status='in_service'
		  	`
	_, err := database.Store.Exec(SQL, assetID)
	if err != nil {
		return err
	}
	return nil
}
