package ad

import (
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"github.com/1URose/marketplace/internal/announcement/use_cases"
	dtoErr "github.com/1URose/marketplace/internal/common/transport/rest/dto"
	"github.com/1URose/marketplace/internal/common/validator"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	service   *use_cases.AdService
	validator *validator.AdAllowedValues
}

func NewHandler(service *use_cases.AdService, validator *validator.AdAllowedValues) *Handler {
	log.Println("[handler:ad] NewHandler initialized")
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// CreateAd godoc
// @Summary      Создать новое объявление
// @Description  Создаёт объявление от имени текущего пользователя
// @Tags         ads
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer <access_token>"
// @Param        ad             body      dto.CreateAdRequest       true  "Данные для создания объявления"
// @Success      201            {object}  dto.CreateAdResponse     "Созданное объявление"
// @Failure      400            {object}  dto.ErrorResponse        "Неверные данные запроса"
// @Failure      401            {object}  dto.ErrorResponse        "Неавторизован"
// @Failure      500            {object}  dto.ErrorResponse        "Внутренняя ошибка"
// @Router       /ad [post]
func (h *Handler) CreateAd(ctx *gin.Context) {
	log.Println("[handler:ad] CreateAd called")

	userId := ctx.GetInt("userId")
	userEmail, _ := ctx.Get("userEmail")
	emailStr, _ := userEmail.(string)

	var req dto.CreateAdRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("[handler:ad][ERROR] bind body: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validator.ValidateCreateAd(req); err != nil {
		log.Println("[handler:ad][ERROR] ValidateCreateAd: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	ad, err := h.service.CreateAd(ctx, userId, &req)
	ad.AuthorEmail = emailStr

	if err != nil {
		log.Println("[handler:ad][ERROR] CreateAd: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to create ad",
		})
		return
	}

	createdAd := dto.NewCreateAdResponse(ad)

	log.Println("[handler:ad] CreateAd succeeded: createAdResponse=", createdAd)

	ctx.JSON(http.StatusCreated, createdAd)
}

// GetAllAds godoc
// @Summary      Получить список объявлений
// @Description  Возвращает постраничный, сортируемый и фильтруемый список объявлений
// @Tags         ads
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Номер страницы"                   default(1)
// @Param        sort_by     query     string  false  "Сортировать по полю"              Enums(created_at,price) default(created_at)
// @Param        sort_order  query     string  false  "Порядок сортировки"               Enums(desc,asc)         default(desc)
// @Param        min_price   query     int     false  "Минимальная цена фильтрации"      minimum(0)
// @Param        max_price   query     int     false  "Максимальная цена фильтрации"     minimum(0)
// @Success      200         {array}   dto.GetAllAdsResponse      "Список объявлений и количество страниц"
// @Failure      400         {object}  dto.ErrorResponse          "Неверные параметры запроса"
// @Failure      500         {object}  dto.ErrorResponse          "Внутренняя ошибка сервера"
// @Router       /ads [get]
func (h *Handler) GetAllAds(ctx *gin.Context) {
	log.Println("[handler:ad] GetAllAds called")

	var req dto.GetAllAdsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Println("[handler:ad][ERROR] bind query: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error: "Invalid request query",
		})
		return
	}

	if err := h.validator.ValidateGetAllAdsRequest(&req); err != nil {
		log.Println("[handler:ad][ERROR] ValidateGetAllAdsRequest: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ads, countPages, err := h.service.GetAllAds(ctx, &req)
	if err != nil {
		log.Println("[handler:ad][ERROR] GetAllAds: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to get ads",
		})
		return
	}

	var userId int
	if ctx.GetBool("isAuthenticated") {
		userId = ctx.GetInt("userId")
		log.Printf("[handler:ad] GetAllAds: authenticated userId=%d", userId)
	} else {
		log.Println("[handler:ad] GetAllAds: guest access")
	}

	adsResponse := dto.NewGetAllAdsResponse(ads, userId, countPages)

	log.Println("[handler:ad] GetAllAds succeeded: adsResponse=", adsResponse)

	ctx.JSON(http.StatusOK, adsResponse)

}
