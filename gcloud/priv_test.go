package gcloud

import (
	"context"
	"log"
	"testing"
)

func Test_getPid(t *testing.T) {
	pid := getPid(context.Background())
	if len(pid) == 0 {
		t.Fatal(pid)
	}
	if pid != devProjectId {
		t.Fatal("pid different")
	}
	log.Println("ProjectID: ", pid)
}
