package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUser struct {
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Role        string `json:"role" db:"role"`
	Type        string `json:"type" db:"type"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
	Password    string `json:"password" db:"password"`
}

type LoginUser struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserAuth struct {
	ID       string `json:"id" db:"id"`
	Password string `json:"password" db:"password"`
}
type UserCtx struct {
	UserID    uuid.UUID `json:"userID"`
	SessionID string    `json:"sessionID"`
	Role      string    `json:"role"`
}
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Type  string `json:"type"`
}
type Asset struct {
	Brand         string    `json:"brand" db:"brand"`
	Model         string    `json:"model" db:"model"`
	SerialNumber  string    `json:"serialNumber" db:"serial_number"`
	AssetType     string    `json:"AssetType" db:"type"`
	Status        string    `json:"status" db:"status"`
	Owner         string    `json:"owner" db:"owner"`
	WarrantyStart time.Time `json:"warrantyStart" db:"warranty_start"`
	WarrantyEnd   time.Time `json:"warrantyEnd" db:"warranty_end"`
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
