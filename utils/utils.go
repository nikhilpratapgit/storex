package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	StatusCode    int    `json:"statusCode"`
	Error         string `json:"error"`
	MessageToUser string `json:"message_to_user"`
}

func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}
func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			fmt.Printf("failed to respond JSON with error:%v", err)
		}
	}
}
func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string) {
	w.WriteHeader(statusCode)
	var errString string
	if err != nil {
		errString = err.Error()
	}
	newError := Error{
		StatusCode:    statusCode,
		Error:         errString,
		MessageToUser: messageToUser,
	}
	if err := json.NewEncoder(w).Encode(newError); err != nil {
		fmt.Printf("failed to send error %v", err)
	}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
