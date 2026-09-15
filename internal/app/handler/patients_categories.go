package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"patients_categories/internal/app/ds"
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

func (h *Handler) AddHandler(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftByUser(CreatorUserID)
	if err != nil {
		h.errorHandler(ctx, err)
		return
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"draft":    draft,
		"minio":    h.MinioURL,
		"defImage": h.defaultImageURL(),
		"defVideo": h.defaultVideoURL(),
	})
}

func (h *Handler) CreateHandler(ctx *gin.Context) {
	name := strings.TrimSpace(ctx.PostForm("name"))
	if name == "" {
		ctx.Redirect(http.StatusFound, "/add")
		return
	}

	existing, err := h.Repository.GetDraftByUser(CreatorUserID)
	if err != nil {
		h.errorHandler(ctx, err)
		return
	}
	if existing == nil {
		draft, err := h.Repository.CreateCategory(name, CreatorUserID)
		if err != nil {
			h.errorHandler(ctx, err)
			return
		}

		src, srcErr := h.Repository.GetPublishedByName(name)
		if srcErr != nil {
			h.errorHandler(ctx, srcErr)
			return
		}
		if src != nil {
			if err := h.Repository.PrefillDraft(draft.ID, *src); err != nil {
				h.errorHandler(ctx, err)
				return
			}
		}
	}

	ctx.Redirect(http.StatusFound, "/add")
}

func (h *Handler) PublishHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.PostForm("calculation_id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, fmt.Errorf("invalid calculation id"))
		return
	}

	name := strings.TrimSpace(ctx.PostForm("name"))
	if name == "" {
		ctx.Redirect(http.StatusFound, "/add")
		return
	}
	description := ctx.PostForm("description")

	age, err := strconv.Atoi(ctx.PostForm("age"))
	if err != nil || age < 1 {
		age = 0
	}
	gender := ctx.PostForm("gender")
	if gender != "Male" && gender != "Female" {
		gender = "Male"
	}
	weight, err := strconv.ParseFloat(ctx.PostForm("weight"), 64)
	if err != nil {
		weight = 0
	}
	height, err := strconv.Atoi(ctx.PostForm("height"))
	if err != nil {
		height = 0
	}

	if err := h.Repository.PublishCategory(uint(id), name, description, age, gender, weight, height); err != nil {
		h.errorHandler(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/feed?id=%d", id))
}

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