package module

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2/google"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var secretManagerInfo *SecretManagerInfo

type SecretManagerInfo struct {
	projectID string
	secretID  string
}

func InitSecretManagerInfo(projectID string, secretID string) {
	secretManagerInfo = &SecretManagerInfo{
		projectID: projectID,
		secretID:  secretID,
	}
}

func getDownloadURL(w http.ResponseWriter, r *http.Request) {
	mod := Module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
		Version:   chi.URLParam(r, "version"),
	}
	fmt.Println("Enters the getDownloadURL function")
	ctx := context.Background()
	fmt.Println("gcs bucket : ", gcsBucket)
	downloadURL, err := gcsBucket.getModule(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("X-Terraform-Get", downloadURL)
	fmt.Println("DownloadURL", downloadURL)
	w.WriteHeader(http.StatusNoContent)
}

func getServiceAccountFromSecretManager() []byte {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	projectID := secretManagerInfo.projectID
	secretID := secretManagerInfo.secretID
	version := "latest"

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectID + "/secrets/" + secretID + "/versions/" + version,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return result.Payload.Data
}

func (b *GCSBackend) getModule(mod Module, ctx context.Context) (string, error) {
	fmt.Println("mod :", modPath(mod))
	fmt.Println("context : ", ctx)

	saKeyFile := getServiceAccountFromSecretManager()

	cfg, err := google.JWTConfigFromJSON(saKeyFile)
	if err != nil {
		log.Fatalln(err)
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey:     cfg.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(2 * time.Minute),
	}

	signedUrl, err := b.client.Bucket(b.bucket).SignedURL(modPath(mod), opts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(signedUrl), nil
}
