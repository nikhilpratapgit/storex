package models

import (
	"time"
)

type Asset struct {
	Brand         string    `json:"brand" db:"brand" validate:"required"`
	Model         string    `json:"model" db:"model" validate:"required"`
	SerialNumber  string    `json:"serialNumber" db:"serial_number" validate:"required"`
	AssetType     string    `json:"type" db:"type" validate:"required,oneof=laptop keyboard mouse mobile"`
	Status        string    `json:"status" db:"status" validate:"required,oneof=available assigned in_service waiting_for_repair damaged"`
	Owner         string    `json:"owner" db:"owner" validate:"required,oneof=client remotestate"`
	WarrantyStart time.Time `json:"warrantyStart" db:"warranty_start" validate:"required"`
	WarrantyEnd   time.Time `json:"warrantyEnd" db:"warranty_end" validate:"required"`

	Laptop   LaptopSpecs   `json:"laptopSpecs,omitempty"`
	Keyboard KeyboardSpecs `json:"keyboardSpecs,omitempty"`
	Mouse    MouseSpecs    `json:"mouseSpecs,omitempty"`
	Mobile   MobileSpecs   `json:"mobileSpecs,omitempty"`
}

type AssetInfo struct {
	Brand        string    `json:"brand" db:"brand"`
	Model        string    `json:"model" db:"model"`
	Type         string    `json:"type" db:"type"`
	SerialNumber string    `json:"serialNumber" db:"serial_number"`
	AssetStatus  string    `json:"assetStatus" db:"status"`
	AssignedTo   string    `json:"assignedTo" db:"assigned_to"`
	Owner        string    `json:"owner" db:"owner"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
type LaptopSpecs struct {
	AssetID         string `json:"assetID" db:"asset_id"`
	Processor       string `json:"processor" db:"processor"`
	Ram             string `json:"ram" db:"ram"`
	Storage         string `json:"storage" db:"storage"`
	OperatingSystem string `json:"operatingSystem" db:"operating_system"`
	Charger         string `json:"charger" db:"charger"`
	DevicePassword  string `json:"devicePassword" db:"device_password"`
}
type KeyboardSpecs struct {
	AssetID      string `json:"assetID" db:"asset_id"`
	Layout       string `json:"layout" db:"layout"`
	Connectivity string `json:"connectivity" db:"connectivity"`
}
type MouseSpecs struct {
	AssetID      string `json:"assetID" db:"asset_id"`
	Dpi          int    `json:"dpi" db:"dpi"`
	Connectivity string `json:"connectivity" db:"connectivity"`
}
type MobileSpecs struct {
	AssetID         string `json:"assetID" db:"asset_id"`
	OperatingSystem string `json:"operatingSystem" db:"operating_system"`
	Ram             string `json:"ram" db:"ram"`
	Storage         string `json:"storage" db:"storage"`
	Charger         string `json:"charger" db:"charger"`
	DevicePassword  string `json:"devicePassword" db:"device_password"`
}
type AssignedAsset struct {
	UserId string `json:"userId" db:"user_id"`
}
type ServiceAsset struct {
	ServiceStart time.Time `json:"serviceStart" db:"service_start" validate:"required"`
}
type DeleteAsset struct {
	ArchivedBy string `json:"archivedBy" db:"archived_by"`
}
type AssetsDashboard struct {
	Total            int `json:"totalAssets" db:"total"`
	Available        int `json:"available" db:"available"`
	Assigned         int `json:"assigned" db:"assigned"`
	WaitingForRepair int `json:"waitingForRepair" db:"waiting_for_repair"`
	InService        int `json:"inService" db:"in_service"`
	Damaged          int `json:"damaged" db:"damaged"`
}

type AssetDashboardData struct {
	AssetCount AssetsDashboard
	Assets     []AssetInfo
}

type AssetInfoRequest struct {
	ID     string `json:"id" db:"id"`
	Brand  string `json:"brand" db:"brand"`
	Model  string `json:"model" db:"model"`
	Status string `json:"status" db:"status"`
	Type   string `json:"type" db:"type"`
}
type UpdateAssetRequest struct {
	Brand         string    `json:"brand" validate:"required"`
	Model         string    `json:"model" validate:"required"`
	SerialNumber  string    `json:"serialNumber" validate:"required"`
	Type          string    `json:"type" validate:"required,oneof=laptop keyboard mouse mobile"`
	Owner         string    `json:"owner" validate:"required,oneof=client remotestate"`
	WarrantyStart time.Time `json:"warrantyStart" validate:"required"`
	WarrantyEnd   time.Time `json:"warrantyEnd" validate:"required"`

	Laptop   *LaptopSpecs   `json:"laptop,omitempty"`
	Mouse    *MouseSpecs    `json:"mouse,omitempty"`
	Keyboard *KeyboardSpecs `json:"keyboard,omitempty"`
	Mobile   *MobileSpecs   `json:"mobile,omitempty"`
}

type UserAssetDetail struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	Role        string `db:"role"`
	Type        string `db:"type"`
	CreatedAt   string `db:"created_at"`

	AssetID   *string `db:"asset_id"`
	Brand     *string `db:"brand"`
	Model     *string `db:"model"`
	Status    *string `db:"status"`
	AssetType *string `db:"asset_type"`
}
