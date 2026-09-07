package handler

import (
	"net/http"
	"patients_categories/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
	MinioURL   string
}

func NewHandler(r *repository.Repository, minioURL string) *Handler {
	return &Handler{
		Repository: r,
		MinioURL:   minioURL,
	}
}

func (h *Handler) FeedHandler(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "1")
	next := ctx.Query("next")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 1
	}

	var service repository.Service
	if next == "true" {
		service, err = h.Repository.GetNextService(id)
	} else {
		service, err = h.Repository.GetServiceByID(id)
	}
	if err != nil {
		ctx.String(http.StatusNotFound, "Service not found")
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"service": service,
		"minio":   h.MinioURL,
		"likes":   len(service.Likes),
	})
}

func (h *Handler) AddHandler(ctx *gin.Context) {
	draft, err := h.Repository.GetDraft()
	if err != nil {
		ctx.String(http.StatusNotFound, "Draft not found")
		return
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"service": draft,
		"minio":   h.MinioURL,
	})
}

func (h *Handler) CardsHandler(ctx *gin.Context) {
	ageStr := ctx.Query("age")
	var services []repository.Service
	var err error

	if ageStr != "" {
		age, e := strconv.Atoi(ageStr)
		if e == nil && age > 0 {
			services, err = h.Repository.GetServicesByAge(age)
		} else {
			services, err = h.Repository.GetServices()
		}
	} else {
		services, err = h.Repository.GetServices()
	}

	if err != nil {
		services = []repository.Service{}
	}

	ctx.HTML(http.StatusOK, "cards.html", gin.H{
		"services": services,
		"minio":    h.MinioURL,
		"age":      ageStr,
	})
}
