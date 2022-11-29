package main

import (
	"ekyc-app/source/api/portal"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	portal.Reg(router)
	router.Run(":8080")
}
