package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"patients_categories/internal/app/config"
	"patients_categories/internal/app/ds"
	redis2 "patients_categories/internal/app/redis"
	"patients_categories/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
	Redis      *redis2.Client
	Config     *config.Config
}

func NewHandler(r *repository.Repository, redisClient *redis2.Client, cfg *config.Config) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redisClient,
		Config:     cfg,
	}
}

func (h *Handler) RegisterAPI(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	// Без авторизации доступно только чтение и регистрация/аутентификация.
	api.GET("/patients_categories", h.AuthOptional, h.APIGetPatientsCategories)
	api.GET("/patients_categories/feed", h.AuthOptional, h.APIFeed)
	api.POST("/users/register", h.APIRegister)
	api.POST("/users/auth", h.APIAuth)

	// Методы создателя и все изменяющие методы — только с валидным JWT.
	auth := api.Group("", h.AuthRequired)
	auth.POST("/patients_categories", h.APICreatePatientsCategorie)
	auth.GET("/patients_categories/draft", h.APIGetDraft)
	auth.PUT("/patients_categories/:id/publish", h.APIPublish)
	auth.DELETE("/patients_categories/:id", h.APIDelete)
	auth.POST("/patients_categories/:id/like", h.APILike)
	auth.POST("/users/logout", h.APILogout)
}

// APIGetPatientsCategories godoc
// @Summary      Плитка опубликованных услуг
// @Description  Список опубликованных услуг с признаком is_creator; фильтр по возрасту ?age=N
// @Tags         patients_categories
// @Produce      json
// @Param        age query int false "возраст для фильтрации"
// @Success      200 {array} ds.PatientsCategoriesTileSerializer
// @Failure      400 {string} string "некорректный age"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories [get]
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

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategoriesTiles(services, likes, getUserID(ctx)))
}

// APICreatePatientsCategorie godoc
// @Summary      Создать черновик услуги
// @Description  Создаёт черновик с файлами изображения и видео; создатель берётся из JWT
// @Tags         patients_categories
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        name  formData string true  "название услуги"
// @Param        image formData file   false "изображение"
// @Param        video formData file   false "короткое видео"
// @Success      201 {object} ds.PatientsCategoriesSerializer
// @Failure      400 {string} string "некорректные данные или файл"
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      409 {string} string "у пользователя уже есть черновик"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories [post]
func (h *Handler) APICreatePatientsCategorie(ctx *gin.Context) {
	creatorID := getUserID(ctx)

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
		Name:   name,
		UserID: creatorID,
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

// APIGetDraft godoc
// @Summary      Получить черновик текущего пользователя
// @Description  Черновик пользователя из JWT, id не указывается (не более 1 записи)
// @Tags         patients_categories
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} ds.PatientsCategoriesSerializer
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories/draft [get]
func (h *Handler) APIGetDraft(ctx *gin.Context) {
	draft, err := h.Repository.GetDraftByUser(getUserID(ctx))
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

// APIFeed godoc
// @Summary      Лента опубликованных услуг
// @Description  Первая услуга или следующая после id при ?id=X&next=true; признак is_liked
// @Tags         patients_categories
// @Produce      json
// @Param        id   query int  false "id услуги"
// @Param        next query bool false "получить следующую услугу"
// @Success      200 {object} ds.PatientsCategoriesFeedSerializer
// @Failure      400 {string} string "некорректный id"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories/feed [get]
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

	liked, err := h.Repository.IsLiked(calc.ID, getUserID(ctx))
	isLiked := 0
	if err == nil && liked {
		isLiked = 1
	}

	ctx.JSON(http.StatusOK, ds.SerializePatientsCategorieFeed(calc, likes, isLiked))
}

type PublishRequest struct {
	Description string  `json:"description"`
	Age         int     `json:"age"`
	Weight      float64 `json:"weight"`
}

// APIPublish godoc
// @Summary      Опубликовать черновик услуги
// @Description  Заполняет description, age, weight и переводит статус draft -> published; только создатель услуги
// @Tags         patients_categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int            true "id услуги"
// @Param        request body PublishRequest true "описание и параметры"
// @Success      200 {object} ds.PatientsCategoriesSerializer
// @Failure      400 {string} string "некорректные данные или услуга не черновик"
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      403 {string} string "услуга другого пользователя"
// @Failure      404 {string} string "услуги нет"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories/{id}/publish [put]
func (h *Handler) APIPublish(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var req PublishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	age := req.Age
	if age < 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}
	weight := req.Weight
	if weight < 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.UpdateAndPublishDraft(uint(id), getUserID(ctx), req.Description, age, weight); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrForbidden) {
			ctx.Status(http.StatusForbidden)
			return
		}
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

// APIDelete godoc
// @Summary      Удалить услугу
// @Description  Soft delete: статус меняется на deleted; только услуги этого пользователя
// @Tags         patients_categories
// @Security     BearerAuth
// @Param        id path int true "id услуги"
// @Success      200 {string} string "успех"
// @Failure      400 {string} string "некорректный id"
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      403 {string} string "услуга другого пользователя"
// @Failure      404 {string} string "услуги нет"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories/{id} [delete]
func (h *Handler) APIDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.DeleteCategory(uint(id), getUserID(ctx)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrForbidden) {
			ctx.Status(http.StatusForbidden)
			return
		}
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

type LikeRequest struct {
	Like int `json:"like"`
}

type LikeResponse struct {
	ID    uint `json:"id"`
	Likes int  `json:"likes"`
}

// APILike godoc
// @Summary      Поставить/снять лайк
// @Description  Поле like: 1 ставит лайк, 0 снимает; лайк от текущего пользователя
// @Tags         patients_categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int       true "id услуги"
// @Param        request body LikeRequest true "1 или 0"
// @Success      200 {object} LikeResponse
// @Failure      400 {string} string "некорректный id или значение like"
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /patients_categories/{id}/like [post]
func (h *Handler) APILike(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var req LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if req.Like != 0 && req.Like != 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.Repository.SetLike(uint(id), getUserID(ctx), req.Like == 1); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	count, err := h.Repository.GetLikesCount(uint(id))
	if err != nil {
		count = 0
	}

	ctx.JSON(http.StatusOK, LikeResponse{ID: uint(id), Likes: count})
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

const jwtPrefix = "Bearer "
const userIDKey = "userID"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *Handler) createToken(userID uint) (string, error) {
	cfg := h.Config.JWT
	token := jwt.NewWithClaims(cfg.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "body-mass-index",
		},
		UserID: userID,
	})
	return token.SignedString([]byte(cfg.Token))
}

func (h *Handler) parseToken(tokenStr string) (*ds.JWTClaims, error) {
	claims := &ds.JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func bearerToken(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(header, jwtPrefix) {
		return ""
	}
	return header[len(jwtPrefix):]
}

// AuthRequired пускает дальше только запрос с валидным токеном,
// которого нет в чёрном списке. Иначе — 401.
func (h *Handler) AuthRequired(ctx *gin.Context) {
	tokenStr := bearerToken(ctx)
	if tokenStr == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	blacklisted, err := h.Redis.IsJWTBlacklisted(ctx.Request.Context(), tokenStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if blacklisted {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := h.parseToken(tokenStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set(userIDKey, claims.UserID)
	ctx.Next()
}

// AuthOptional — для методов чтения: если токен валидный, привязывает
// пользователя, иначе продолжает запрос как неавторизованный.
func (h *Handler) AuthOptional(ctx *gin.Context) {
	ctx.Set(userIDKey, uint(0))

	tokenStr := bearerToken(ctx)
	if tokenStr == "" {
		ctx.Next()
		return
	}

	blacklisted, err := h.Redis.IsJWTBlacklisted(ctx.Request.Context(), tokenStr)
	if err != nil || blacklisted {
		ctx.Next()
		return
	}

	claims, err := h.parseToken(tokenStr)
	if err != nil {
		ctx.Next()
		return
	}

	ctx.Set(userIDKey, claims.UserID)
	ctx.Next()
}

func getUserID(ctx *gin.Context) uint {
	v, ok := ctx.Get(userIDKey)
	if !ok {
		return 0
	}
	id, _ := v.(uint)
	return id
}

// APIRegister godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя, пароль сохраняется в виде sha1-хэша
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "логин и пароль"
// @Success      201 {object} ds.UserSerializer
// @Failure      400 {string} string "некорректные данные"
// @Failure      409 {string} string "логин уже занят"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /users/register [post]
func (h *Handler) APIRegister(ctx *gin.Context) {
	var req RegisterRequest
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

	u, err := h.Repository.CreateUser(req.Username, generateHashString(req.Password))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, ds.SerializeUser(u))
}

// APIAuth godoc
// @Summary      Аутентификация
// @Description  Проверяет логин и пароль, возвращает JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "логин и пароль"
// @Success      200 {object} LoginResponse
// @Failure      400 {string} string "некорректные данные"
// @Failure      401 {string} string "неверный логин или пароль"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /users/auth [post]
func (h *Handler) APIAuth(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	u, err := h.Repository.GetUserByUsername(req.Username)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if u == nil || u.PasswordHash != generateHashString(req.Password) {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	token, err := h.createToken(u.ID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		ExpiresIn:   int64(h.Config.JWT.ExpiresIn),
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

// APILogout godoc
// @Summary      Деавторизация
// @Description  Кладёт текущий токен в чёрный список Redis до истечения его срока
// @Tags         users
// @Security     BearerAuth
// @Success      200 {string} string "успех"
// @Failure      401 {string} string "нет или неверный токен"
// @Failure      500 {string} string "внутренняя ошибка"
// @Router       /users/logout [post]
func (h *Handler) APILogout(ctx *gin.Context) {
	tokenStr := bearerToken(ctx)
	if tokenStr == "" {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	claims, err := h.parseToken(tokenStr)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
	if err := h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), tokenStr, ttl); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
