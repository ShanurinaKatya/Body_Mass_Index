package handler

import (
	"net/http"
	"strconv"

	"patients_categories/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CardsHandler(ctx *gin.Context) {
	ageStr := ctx.Query("age")
	var services []ds.PatientsCategorie
	var err error

	if ageStr != "" {
		age, convErr := strconv.Atoi(ageStr)
		if convErr == nil && age > 0 {
			services, err = h.Repository.GetPublishedByAge(age)
		} else {
			services, err = h.Repository.GetPublished()
		}
	} else {
		services, err = h.Repository.GetPublished()
	}
	if err != nil {
		services = []ds.PatientsCategorie{}
	}

	likes, err := h.Repository.GetLikesCounts()
	if err != nil {
		likes = map[int]int{}
	}

	ctx.HTML(http.StatusOK, "cards.html", gin.H{
		"services": services,
		"minio":    h.MinioURL,
		"age":      ageStr,
		"likes":    likes,
		"defImage": h.defaultImageURL(),
	})
}
