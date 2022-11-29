package cfg

import (
	"context"
	"testing"
)

func Test_Get(t *testing.T) {
	info, err := Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(info.PrePrivKey) == 0 {
		t.Fatal("prePreiKey invalid")
	}
	if len(info.LogBucket) == 0 {
		t.Fatal("log bucket invalid")
	}
	t.Fatal("OKE")
}
