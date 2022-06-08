package http

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fvbock/endless"
)

func NewHttpServer() {
	var err error
	if os.Getenv("CLOUDGOBRRR_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	if os.Getenv("TRUSTED_PROXIES") == "nil" {
		err = router.SetTrustedProxies(nil)
	} else {
		err = router.SetTrustedProxies(strings.Split(os.Getenv("TRUSTED_PROXIES"), " "))
	}
	if err != nil {
		panic(err)
	}

	newRouter(router)

	err = endless.ListenAndServe(":"+os.Getenv("HTTP_PORT"), router)
	if err != nil {
		fmt.Println(err.Error())
	}
}
