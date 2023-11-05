package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kaellybot/kaelly-twitter/services"
)

type Controller interface {
	Run()
}

type Impl struct {
	r       *gin.Engine
	service services.Service
}
