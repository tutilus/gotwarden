package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowDevices shows all devices
func (ctx *WardenCtx) ShowDevices(c *gin.Context) {
	devices, err := ctx.Db.AllDevices()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"devices": devices})
	}
}

// ShowUsers shows all the users data
func (ctx *WardenCtx) ShowUsers(c *gin.Context) {
	users, err := ctx.Db.AllUsers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// ShowFolders shows all folders data
func (ctx *WardenCtx) ShowFolders(c *gin.Context) {
	folders, err := ctx.Db.AllFolders()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"folders": folders})
	}
}

// ShowCiphers shows all folders data
func (ctx *WardenCtx) ShowCiphers(c *gin.Context) {
	ciphers, err := ctx.Db.AllCiphers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{"folders": ciphers})
	}
}
