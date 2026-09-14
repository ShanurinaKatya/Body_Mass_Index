package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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
		if _, err := h.Repository.CreateCategory(name, CreatorUserID); err != nil {
			h.errorHandler(ctx, err)
			return
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
