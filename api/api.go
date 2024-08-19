package api

import (
	"github.com/gin-gonic/gin"
	"main/db"
)

func Setup(router *gin.Engine) {
	db.SetupDB()
	SetupAuthenAPI(router)
	SetupProductAPI(router)
	SetupTransactionAPI(router)
}