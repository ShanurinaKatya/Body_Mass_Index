package handler

import (
	"net/http"

	"patients_categories/internal/app/repository"

	"github.com/gin-gonic/gin"
)

const CreatorUserID uint = 1

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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/feed")
	})
	router.GET("/feed", h.FeedHandler)
	router.GET("/add", h.AddHandler)
	router.GET("/cards", h.CardsHandler)
	router.POST("/add", h.CreateHandler)
	router.POST("/publish", h.PublishHandler)
	router.POST("/delete", h.DeleteHandler)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, err error) {
	ctx.String(http.StatusNotFound, err.Error())
}

func (h *Handler) defaultImageURL() string {
	return "/static/image/default.jpg"
}

func (h *Handler) defaultVideoURL() string {
	return "/static/image/Default.mp4"
}
