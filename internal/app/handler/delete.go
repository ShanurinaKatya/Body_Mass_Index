package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.PostForm("calculation_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid calculation id"})
		return
	}

	if err := h.Repository.DeleteCategory(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/cards")
}
