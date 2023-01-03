package auth

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service IAuthService
}

func (controller *AuthController) GetAuthenticateToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("authorization")
	isAuthenticated, err := controller.Service.ValidateToken(authHeader)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error at GetAuthenticateToken: %s\n%s\n", err.Error(), debug.Stack())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
		})
		return
	}

	if isAuthenticated {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"authenticated": true,
			},
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"data": gin.H{
				"authenticated": false,
			},
		})
	}
}

func (controller *AuthController) PostGenerateAuthToken(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"authenticated": false,
			},
		})
		return
	}

	validUser := controller.Service.ValidateLogin(json.Email, json.Password)
	if validUser == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"data": gin.H{
				"authenticated": false,
			},
		})
		return
	}
	token, err := controller.Service.SignJWT(*validUser)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"token": token,
			},
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error at PostGenerateAuthToken: %s\n%s\n", err.Error(), debug.Stack())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    gin.H{},
		})
	}
}
