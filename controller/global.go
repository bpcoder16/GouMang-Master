package controller

import (
	"github.com/gin-gonic/gin"
)

type Global struct{}

func (g *Global) Config(ctx *gin.Context) {
	returnSuccessJson(ctx, gin.H{})
}
