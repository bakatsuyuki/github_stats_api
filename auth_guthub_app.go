package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	jwt "github.com/dgrijalva/jwt-go"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func init() {
	functions.HTTP("HelloHTTP", authGitHubApp)
}

func authGitHubApp(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken("189639")
	if err != nil {
		return
	}

	client := new(http.Client)
	req, _ := http.NewRequest("POST", "https://api.github.com/app/installations/24886519/access_tokens", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	w.Write(body)
	return
}

func generateToken(appID string) (string, error) {
	/*
		b, err := ioutil.ReadFile("github_app.pem")
		if err != nil {
			return "", err
		}*/
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.GetSecretRequest{
		Name: "AppPrivateKey",
	}

	// Call the API.
	result, err := client.GetSecret(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %v", err)
	}
	replication := result.Replication
	b := []byte(replication.String())
	c := &jwt.StandardClaims{
		Issuer:    appID,
		ExpiresAt: time.Now().Unix() + 60,
		IssuedAt:  time.Now().Unix() - 10,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), c)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return "", err
	}
	t, err := token.SignedString(key)
	return t, err
}
