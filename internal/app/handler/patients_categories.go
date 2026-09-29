package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"patients_categories/internal/app/ds"
	"patients_categories/internal/app/repository"
	"patients_categories/internal/app/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	// ===== Домен услуги: patients_categories =====
	api.GET("/patients_categories", h.APIGetPatientsCategories)        // список published + фильтр ?age=
	api.POST("/patients_categories", h.APICreatePatientsCategorie)     // добавление + файлы (multipart)
	api.GET("/patients_categories/draft", h.APIGetDraft)               // черновик текущего пользователя
	api.GET("/patients_categories/feed", h.APIFeed)                    // лента: без id или ?id=X&next=true
	api.PUT("/patients_categories/:id/publish", h.APIPublish)          // draft -> published
	api.DELETE("/patients_categories/:id", h.APIDelete)                // soft delete
	api.POST("/patients_categories/:id/like", h.APILike)               // {"like": 1|0}

	// ===== Домен пользователь =====
	api.POST("/users/register", h.APIRegister)
	api.POST("/users/auth", h.APIAuth)
	api.POST("/users/logout", h.APILogout)
}

func defaultGender(g string) string {
	if g != "Male" && g != "Female" {
		return "Male"
	}
	return g
}

// =========================================================
// API: Домен услуги
// =========================================================

// GET /api/patients_categories — список опубликованных, фильтр ?age=
func (h *Handler) APIGetPatientsCategories(ctx *gin.Context) {
	ageStr := ctx.Query("age")

	var services []ds.PatientsCategorie
	var err error

	if ageStr != "" {
		age, convErr := strconv.Atoi(ageStr)
		if convErr != nil || age < 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}
		services, err = h.Repository.GetPublishedFiltered(age)
	} else {
		services, err = h.Repository.GetPublished()
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	likes, err := h.Repository.GetLikesCounts()
	if err != nil {
		likes = map[int]int{}
	}

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategories(services, likes))
}

// POST /api/patients_categories — multipart: поля + файлы image/video -> MinIO, статус draft
func (h *Handler) APICreatePatientsCategorie(ctx *gin.Context) {
	creatorID := user.Get().ID

	if err := ctx.Request.ParseMultipartForm(8 << 20); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(ctx.Request.FormValue("name"))
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	existing, err := h.Repository.GetDraftByUser(creatorID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if existing != nil {
		ctx.Status(http.StatusConflict)
		return
	}

	description := ctx.Request.FormValue("description")

	ageStr := ctx.Request.FormValue("age")
	age, ageErr := strconv.Atoi(ageStr)
	if ageErr != nil || age < 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}
	gender := defaultGender(ctx.Request.FormValue("gender"))

	weightStr := ctx.Request.FormValue("weight")
	weight, weightErr := strconv.ParseFloat(weightStr, 64)
	if weightErr != nil || weight < 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	heightStr := ctx.Request.FormValue("height")
	height, heightErr := strconv.Atoi(heightStr)
	if heightErr != nil || height < 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Валидация опциональных файлов
	var imgHeader, videoHeader *multipart.FileHeader
	imgHeader, err = ctx.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		ctx.Status(http.StatusBadRequest)
		return
	}
	videoHeader, err = ctx.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if imgHeader != nil {
		if code, vErr := validateImageUpload(imgHeader); vErr != nil {
			ctx.Status(code)
			return
		}
	}
	if videoHeader != nil {
		if code, vErr := validateVideoUpload(videoHeader); vErr != nil {
			ctx.Status(code)
			return
		}
	}

	pc, err := h.Repository.CreateCategoryDraft(ds.PatientsCategorie{
		Name:        name,
		Description: description,
		Age:         age,
		Gender:      gender,
		Weight:      weight,
		Height:      height,
		UserID:      creatorID,
	})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if imgHeader != nil {
		if err := h.Repository.UploadCategoryImage(pc.ID, imgHeader); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	if videoHeader != nil {
		if err := h.Repository.UploadCategoryVideo(pc.ID, videoHeader); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	updated, err := h.Repository.GetDraftByUser(creatorID)
	if err != nil || updated == nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, ds.SerializePatientsCategorie(*updated, 0))
}

// GET /api/patients_categories/draft — черновик текущего пользователя (id не указывается)
func (h *Handler) APIGetDraft(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftByUser(user.Get().ID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if draft == nil {
		ctx.JSON(http.StatusOK, ds.SerializePatientsCategorie(ds.PatientsCategorie{}, 0))
		return
	}

	likes, err := h.Repository.GetLikesCount(draft.ID)
	if err != nil {
		likes = 0
	}

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategorie(*draft, likes))
}

// GET /api/patients_categories/feed — лента: без id (первая) или ?id=X&next=true
func (h *Handler) APIFeed(ctx *gin.Context) {
	var calc ds.PatientsCategorie
	var err error

	idStr := ctx.Query("id")
	next := ctx.Query("next") == "true"

	switch {
	case next && idStr != "":
		id, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		calc, err = h.Repository.GetNext(uint(id))
	case idStr != "":
		id, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		calc, err = h.Repository.GetPublishedByID(uint(id))
	default:
		calc, err = h.Repository.GetFirstPublished()
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	likes, err := h.Repository.GetLikesCount(calc.ID)
	if err != nil {
		likes = 0
	}

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategorie(calc, likes))
}

// PUT /api/patients_categories/:id/publish — смена статуса draft -> published
func (h *Handler) APIPublish(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.PublishDraft(uint(id)); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	pc, err := h.Repository.GetPublishedByID(uint(id))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	likes, err := h.Repository.GetLikesCount(pc.ID)
	if err != nil {
		likes = 0
	}

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategorie(pc, likes))
}

// DELETE /api/patients_categories/:id — soft delete
func (h *Handler) APIDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.DeleteCategory(uint(id)); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

// POST /api/patients_categories/:id/like — {"like": 1} ставит, {"like": 0} снимает
func (h *Handler) APILike(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var req struct {
		Like int `json:"like"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if req.Like != 0 && req.Like != 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.SetLike(uint(id), user.Get().ID, req.Like == 1); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	count, err := h.Repository.GetLikesCount(uint(id))
	if err != nil {
		count = 0
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id, "likes": count})
}

// =========================================================
// API: Домен пользователь
// =========================================================

// POST /api/users/register — регистрация нового пользователя
func (h *Handler) APIRegister(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	existing, err := h.Repository.GetUserByUsername(req.Username)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if existing != nil {
		ctx.Status(http.StatusConflict)
		return
	}

	u, err := h.Repository.CreateUser(req.Username, req.Password)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, ds.SerializeUser(u))
}

// POST /api/users/auth — заглушка для 4-й лабораторной
func (h *Handler) APIAuth(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// POST /api/users/logout — заглушка для 4-й лабораторной
func (h *Handler) APILogout(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// =========================================================
// Валидация файлов
// =========================================================

func isImage(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}
	for _, t := range imageTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func isVideo(contentType string) bool {
	videoTypes := []string{
		"video/mp4",
		"video/webm",
		"video/quicktime",
	}
	for _, t := range videoTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func validateImageUpload(header *multipart.FileHeader) (int, error) {
	file, err := header.Open()
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось получить файл")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}
	contentType := http.DetectContentType(buffer)
	if !isImage(contentType) {
		return http.StatusBadRequest, fmt.Errorf("файл должен быть изображением")
	}
	if _, err = file.Seek(0, 0); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при обработке файла")
	}
	return 0, nil
}

func validateVideoUpload(header *multipart.FileHeader) (int, error) {
	file, err := header.Open()
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось получить файл")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}
	contentType := http.DetectContentType(buffer)
	if !isVideo(contentType) {
		return http.StatusBadRequest, fmt.Errorf("файл должен быть видео")
	}
	if _, err = file.Seek(0, 0); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при обработке файла")
	}
	return 0, nil
}