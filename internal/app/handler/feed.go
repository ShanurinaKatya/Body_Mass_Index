package handler

import (
	"net/http"
	"strconv"

	"patients_categories/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FeedHandler(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "1")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		id = 1
	}

	var calc ds.PatientsCategorie
	if ctx.Query("next") == "true" {
		calc, err = h.Repository.GetNext(uint(id))
	} else {
		calc, err = h.Repository.GetPublishedByID(uint(id))
	}
	if err != nil {
		h.errorHandler(ctx, err)
		return
	}

	likes, err := h.Repository.GetLikesCount(calc.ID)
	if err != nil {
		likes = 0
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"service":  calc,
		"minio":    h.MinioURL,
		"likes":    likes,
		"defImage": h.defaultImageURL(),
		"defVideo": h.defaultVideoURL(),
	})
}
