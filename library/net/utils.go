package net

import "strings"

const ekycTime string = "ekyctime"

func UrlPrettyParse(raw string) string {
	var (
		pretty string = raw
	)
	for strings.HasSuffix(pretty, " ") {
		pretty = strings.TrimPrefix(pretty, " ")
	}
	for strings.HasSuffix(pretty, "/") || strings.HasSuffix(pretty, " ") {
		pretty = strings.TrimSuffix(pretty, "/")
		pretty = strings.TrimSuffix(pretty, " ")
	}
	if !strings.HasPrefix(pretty, "http") && !strings.HasPrefix(pretty, "ws") {
		pretty = "http://" + pretty
	}
	return pretty
}
