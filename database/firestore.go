package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client
var FirebaseAuthClient *auth.Client
var FirebaseApp *firebase.App

func buildFirestoreAuthJson() []byte {
	projectId := os.Getenv("FIREBASE_PROJECT_ID")
	privateKeyId := os.Getenv("FIREBASE_PRIVATE_KEY_ID")
	privateKey := strings.ReplaceAll(os.Getenv("FIREBASE_PRIVATE_KEY"), "\\n", "\n")
	clientEmail := os.Getenv("FIREBASE_CLIENT_EMAIL")
	clientId := os.Getenv("FIREBASE_CLIENT_ID")

	type serviceAccount struct {
		Type           string `json:"type"`
		ProjectID      string `json:"project_id"`
		PrivateKeyID   string `json:"private_key_id"`
		PrivateKey     string `json:"private_key"`
		ClientEmail    string `json:"client_email"`
		ClientID       string `json:"client_id"`
		AuthURI        string `json:"auth_uri"`
		TokenURI       string `json:"token_uri"`
		AuthProvider   string `json:"auth_provider_x509_cert_url"`
		ClientCertURL  string `json:"client_x509_cert_url"`
		UniverseDomain string `json:"universe_domain"`
	}

	account := serviceAccount{
		Type:           "service_account",
		ProjectID:      projectId,
		PrivateKeyID:   privateKeyId,
		PrivateKey:     privateKey,
		ClientEmail:    clientEmail,
		ClientID:       clientId,
		AuthURI:        "https://accounts.google.com/o/oauth2/auth",
		TokenURI:       "https://oauth2.googleapis.com/token",
		AuthProvider:   "https://www.googleapis.com/oauth2/v1/certs",
		ClientCertURL:  "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-oweum%40secret-427917.iam.gserviceaccount.com",
		UniverseDomain: "googleapis.com",
	}

	json, err := json.Marshal(account)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	return json
}

func InitFirestore() {
	var err error

	ctx := context.Background()
	opt := option.WithCredentialsJSON(buildFirestoreAuthJson())
	projectId := os.Getenv("FIREBASE_PROJECT_ID")
	FirestoreClient, err = firestore.NewClient(ctx, projectId, opt)
	if err != nil {
		fmt.Println("firebase json : " + string(buildFirestoreAuthJson()))
		fmt.Println("firebase secret key : " + os.Getenv("FIREBASE_PRIVATE_KEY"))
		log.Fatalf("Failed to create a Firestore client: %v", err)
	}

	log.Println("Firestore client created")
}

func InitFirebaseAuth() {
	var err error

	opt := option.WithCredentialsJSON(buildFirestoreAuthJson())
	FirebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create a Firebase app: %v", err)
	}

	FirebaseAuthClient, err = FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create a Firebase auth client: %v", err)
	}

	log.Println("Firebase auth client created")
}
