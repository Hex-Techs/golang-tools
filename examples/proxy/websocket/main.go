package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	websocketproxy "github.com/ymichaelson/golang-tools/proxy/websocket"
	"github.com/ymichaelson/klog"
)

func ExampleWebsocketProxy(c *gin.Context) {
	target := "ws://127.0.0.1:8888"

	remote, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("parse target url failed: %s", target),
		})
	}

	proxy := websocketproxy.ProxyHandlerInsecure(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	router := gin.New()

	routerGroup := router.Group("")
	{
		routerGroup.GET("/examples", ExampleWebsocketProxy)
		routerGroup.POST("/examples", ExampleWebsocketProxy)
		routerGroup.PUT("/examples", ExampleWebsocketProxy)
		routerGroup.DELETE("/examples", ExampleWebsocketProxy)
		routerGroup.PATCH("/examples", ExampleWebsocketProxy)
	}

	err := router.Run(":8080")
	if err != nil {
		klog.Errorf("Start websocket proxy server failed with error: %v", err)
		os.Exit(1)
	}
}
