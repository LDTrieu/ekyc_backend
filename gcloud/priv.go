package gcloud

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

const (
	devProjectId = "ekyc-services" //915240804196
)

func getPid(ctx context.Context) string {
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		return devProjectId
	}
	if len(credentials.ProjectID) == 0 {
		return devProjectId
	}
	return credentials.ProjectID
}
