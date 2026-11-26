package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/dto"
)

func writeAPIError(c *gin.Context, status int, body dto.APIErrorBody) {
	c.JSON(status, dto.APIErrorEnvelope{Error: body})
}

func writeBadRequestMessage(c *gin.Context, message string) {
	writeAPIError(c, http.StatusBadRequest, dto.APIErrorBody{Code: "BAD_REQUEST", Message: message})
}