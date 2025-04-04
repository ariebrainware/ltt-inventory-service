package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var JWTSecret = getEnv("JWTSECRET", "")
var JWTSecretByte = []byte(getEnv("JWTSECRET", ""))

type APIResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type APIErrorParams struct {
	Msg string
	Err error
}

type APISuccessParams struct {
	Msg  string
	Data interface{}
}

// Contains function is to check item whether is exist or not in a list and will return bool
func Contains(d string, dl []string) bool {
	for _, v := range dl {
		if v == d {
			return true
		}
	}
	return false
}

// CallErrorNotFound is for return API response not found
func CallErrorNotFound(c *gin.Context, params APIErrorParams) {
	response := APIResponse{
		Success: false,
		Error:   params.Err.Error(),
		Msg:     params.Msg,
		Data:    map[string]interface{}{},
	}
	c.JSON(http.StatusNotFound, response)
}

// CallUserError is for return error from user side
func CallUserError(c *gin.Context, params APIErrorParams) {
	response := APIResponse{
		Success: false,
		Error:   params.Err.Error(),
		Msg:     params.Msg,
		Data:    map[string]interface{}{},
	}
	c.JSON(http.StatusBadRequest, response)
}

// CallServerError is for return API response server error
func CallServerError(c *gin.Context, params APIErrorParams) {
	response := APIResponse{
		Success: false,
		Error:   params.Err.Error(),
		Msg:     params.Msg,
		Data:    map[string]interface{}{},
	}
	c.JSON(http.StatusInternalServerError, response)
}

// CallSuccessOK is for return API response with status code 200, you need to specify msg, and data as function parameter
func CallSuccessOK(c *gin.Context, params APISuccessParams) {
	response := APIResponse{
		Success: true,
		Error:   "",
		Msg:     params.Msg,
		Data:    params.Data,
	}
	c.JSON(http.StatusOK, response)
}

// CallUserFound is for return API response with status code 307 means its redirected
func CallUserFound(c *gin.Context, params APISuccessParams) {
	response := APIResponse{
		Success: true,
		Error:   "",
		Msg:     params.Msg,
		Data:    params.Data,
	}
	c.JSON(http.StatusTemporaryRedirect, response)
}

// CallUserNotAuthorized is for return API response with status code 403, you need to specify msg, and data as function paramenter
func CallUserNotAuthorized(c *gin.Context, params APIErrorParams) {
	response := APIResponse{
		Success: false,
		Error:   params.Err.Error(),
		Msg:     params.Msg,
	}
	c.JSON(http.StatusUnauthorized, response)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func HashPassword(password string) (hashedPassword string) {
	h := hmac.New(sha256.New, JWTSecretByte)
	h.Write([]byte(password))
	hashedPassword = hex.EncodeToString(h.Sum(nil))
	return
}
