package options_test

import (
	"ekyc-app/library/net/options"
	"os"
	"testing"
)

func Test_SetFormData(t *testing.T) {
	var (
		opts = options.CurlOption()
		form = make(map[string]interface{})
	)
	//
	mod, err := os.ReadFile("./curl_test.go")
	if err != nil {
		t.Fatal(err)
	}
	form["curl_test"] = options.FormFile{
		Filename: "curl_test.go",
		File:     mod,
	}
	//
	form["test"] = "test_value"

	if err := opts.SetFormData(form); err != nil {
		t.Fatal(err)
	}

	t.Logf("Header: %+v", opts.Header)
	t.Logf("Body: %+v", opts.Body)
}
