package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

func UiServer(r *gin.Engine) {
	assetsBox, _ := rice.FindBox("assets")
	if assetsBox != nil {
		fmt.Printf("assetsBox is not found\n")
		r.StaticFS("/ui", assetsBox.HTTPBox())
	} else {
		r.Static("/ui", "./assets")
	}
}
