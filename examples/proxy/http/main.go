package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	httpproxy "github.com/ymichaelson/golang-tools/proxy/http"
	"github.com/ymichaelson/klog"
)

func ExampleHttpProxy(c *gin.Context) {
	target := "https://127.0.0.1:8888"

	remote, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("parse target url failed: %s", target),
		})
	}

	proxy := httpproxy.NewHttpProxy(remote, nil)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	router := gin.New()

	routerGroup := router.Group("")
	{
		routerGroup.GET("/examples", ExampleHttpProxy)
		routerGroup.POST("/examples", ExampleHttpProxy)
		routerGroup.PUT("/examples", ExampleHttpProxy)
		routerGroup.DELETE("/examples", ExampleHttpProxy)
		routerGroup.PATCH("/examples", ExampleHttpProxy)
	}

	err := router.Run(":8080")
	if err != nil {
		klog.Errorf("Start http proxy server failed with error: %v", err)
		os.Exit(1)
	}
}
