package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

func UiServer(r *gin.Engine) {
	assetsBox, _ := rice.FindBox("assets")
	if assetsBox != nil {
		fmt.Printf("assetsBox is found\n")
		r.StaticFS("/ui", assetsBox.HTTPBox())
	} else {
		//r.Static("/ui", "./assets")
		r.GET("/ui/*any", DevProxy())
	}
}

func DevProxy() gin.HandlerFunc {
	// https://stackoverflow.com/a/39009974
	target := "localhost:3005"
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			r := c.Request
			req.URL = r.URL
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

