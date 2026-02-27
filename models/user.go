package models

import (
	"time"
)

type RegisterUser struct {
	Name        string `json:"name" db:"name" validate:"required,min=3,max=50"`
	Email       string `json:"email" db:"email" validate:"required,email"`
	Role        string `json:"role" db:"role" validate:"required,oneof=admin employee project_manager asset_manager employee_manager"`
	Type        string `json:"type" db:"type" validate:"required,oneof=full_time intern freelancer" `
	PhoneNumber string `json:"phoneNumber" db:"phone_number" validate:"required,len=10"`
	Password    string `json:"password" db:"password" validate:"required,min=8,max=20"`
}

type LoginUser struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=8,max=20"`
}

//	type UserAuth struct {
//		ID       string `json:"id" db:"id"`
//		Password string `json:"password" db:"password"`
//	}
type UserCtx struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
	Role      string `json:"role"`
}
type User struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Email       string     `json:"email" db:"email"`
	PhoneNumber string     `json:"phoneNumber" db:"phone_number"`
	Role        string     `json:"role" db:"role"`
	Employment  string     `json:"employment" db:"employment"`
	Password    string     `json:"password" db:"password"`
	CreatedAt   *time.Time `db:"created_at"`
	ArchivedAt  *time.Time `db:"archived_at"`
}
type Asset struct {
	Brand         string    `json:"brand" db:"brand" validate:"required"`
	Model         string    `json:"model" db:"model" validate:"required"`
	SerialNumber  string    `json:"serialNumber" db:"serial_number" validate:"required"`
	AssetType     string    `json:"assetType" db:"type" validate:"required,oneof=laptop keyboard mouse mobile"`
	Status        string    `json:"status" db:"status" validate:"required,oneof=available assigned in_service for_repair damaged"`
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
	AssetType    string    `json:"type" db:"type"`
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
	AssignedTo string `json:"assignedTo" db:"assigned_to"`
}
type ServiceAsset struct {
	ServiceStart time.Time `json:"serviceStart" db:"service_start"`
	ServiceEnd   time.Time `json:"serviceEnd" db:"service_end"`
	ReturnedOn   time.Time `json:"returnedOn" db:"returned_on"`
}
type DeleteAsset struct {
	ArchivedBy string `json:"archivedBy" db:"archived_by"`
}
type DashboardSummary struct {
	Total            int `json:"totalAssets" db:"total"`
	Available        int `json:"available" db:"available"`
	Assigned         int `json:"assigned" db:"assigned"`
	WaitingForRepair int `json:"waitingForRepair" db:"waiting_for_repair"`
	InService        int `json:"inService" db:"in_service"`
	Damaged          int `json:"damaged" db:"damaged"`
}

type DashboardData struct {
	Summary DashboardSummary
	Assets  []AssetInfo
}
type UserInfoRequest struct {
	ID           string             `json:"id" db:"id"`
	Name         string             `json:"name" db:"name" validate:"required,min=3,max=50"`
	Email        string             `json:"email" db:"email" validate:"required,email"`
	PhoneNumber  string             `json:"phoneNumber" db:"phone_no" validate:"required,len=10"`
	Role         string             `json:"role" db:"role" validate:"required"`
	Employment   string             `json:"type" db:"type" validate:"required"`
	CreatedAt    string             `json:"createdAt" db:"created_at" validate:"required"`
	AssetDetails []AssetInfoRequest `json:"assetDetails"`
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
	SerialNo      string    `json:"serialNo" validate:"required"`
	Type          string    `json:"type" validate:"required" validate:"required,oneof=laptop keyboard mouse mobile"`
	Owner         string    `json:"owner" validate:"required" validate:"required,oneof=client remotestate"`
	WarrantyStart time.Time `json:"warrantyStart" validate:"required" validate:"required"`
	WarrantyEnd   time.Time `json:"warrantyEnd" validate:"required" validate:"required"`

	Laptop   *LaptopSpecs   `json:"laptop,omitempty"`
	Mouse    *MouseSpecs    `json:"mouse,omitempty"`
	Keyboard *KeyboardSpecs `json:"keyboard,omitempty"`
	Mobile   *MobileSpecs   `json:"mobile,omitempty"`
}
