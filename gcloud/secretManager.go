package gcloud

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetSecret(ctx context.Context, keyName string) (
	secret []byte, err error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return
	}
	defer client.Close()
	pid := getPid(ctx)
	name := fmt.Sprintf("projects/%v/secrets/%v/versions/2", pid, keyName)
	// Build the request.
	req := &secretmanagerpb.
		AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return
	}

	secret = result.Payload.Data

	return
}
