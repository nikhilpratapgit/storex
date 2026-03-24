package models

type RegisterUser struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Type        string `json:"type"  validate:"required,oneof=full_time intern freelancer" `
	PhoneNumber string `json:"phoneNumber"  validate:"required,numeric,len=10"`
	Password    string `json:"password"  validate:"required,min=8,max=20"`
}

type LoginUser struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required,min=2,max=20"`
}

type UserCtx struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
	Role      string `json:"role"`
}
type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
	Password string `json:"password" db:"password"`
}
type UserInfoRequest struct {
	ID           string             `json:"id" db:"id"`
	Name         string             `json:"name" db:"name" validate:"required,min=3,max=50"`
	Email        string             `json:"email" db:"email" validate:"required,email"`
	PhoneNumber  string             `json:"phoneNumber" db:"phone_number" validate:"required,len=10"`
	Role         string             `json:"role" db:"role" validate:"required"`
	Type         string             `json:"type" db:"type" validate:"required"`
	CreatedAt    string             `json:"createdAt" db:"created_at" validate:"required"`
	AssetDetails []AssetInfoRequest `json:"assetDetails"`
}
type FirebaseLoginRequest struct {
	IdToken string `json:"idToken" validate:"required"`
}
type AssignedRole struct {
	Role string `json:"role" validate:"required"`
}

//	type DeleteUser struct {
//		UserID string `json:"userID" db:"user_id"`
//	}
type GoogleLoginRequest struct {
	IdToken string `json:"idToken" validate:"required"`
}

type UserProfile struct {
	Token string `json:"token" validate:"required"`
	//Type        string `json:"userType" validate:"required"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number" validate:"required,len=10"`
	Password    string `json:"password" db:"password"`
}
