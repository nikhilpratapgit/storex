//package utils
//
//import (
//	"context"
//
//	"google.golang.org/api/option"
//)

package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() {

	opt := option.WithCredentialsFile("firebase-service-account.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	FirebaseAuth = client
}
