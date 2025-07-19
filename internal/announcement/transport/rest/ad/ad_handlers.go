package ad

import (
	"errors"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"github.com/1URose/marketplace/internal/announcement/use_cases"
	"github.com/1URose/marketplace/internal/common/validator"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *use_cases.AdService
}

func NewHandler(service *use_cases.AdService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateAd(ctx *gin.Context) {
	log.Println("[handler:ad] CreateAd called")

	userId, err := h.getUserIdFromCtx(ctx)

	if err != nil {
		log.Println("[handler:ad][ERROR] getUserIdFromCtx: ", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}

	var req dto.CreateAdRequest

	if err := ctx.BindJSON(&req); err != nil {
		log.Println("[handler:ad][ERROR] bind body: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validator.ValidateCreateAd(req); err != nil {
		log.Println("[handler:ad][ERROR] ValidateCreateAd: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ad, err := h.service.CreateAd(ctx, userId, &req)

	if err != nil {
		log.Println("[handler:ad][ERROR] CreateAd: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create ad",
		})
		return
	}

	createdAd := dto.NewCreateAdResponse(ad)

	ctx.JSON(http.StatusCreated, createdAd)
}

func (h *Handler) GetAllAds(ctx *gin.Context) {
	log.Println("[handler:ad] GetAllAds called")

	var req dto.GetAllAdsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Println("[handler:ad][ERROR] bind query: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validator.ValidateGetAllAdsRequest(&req); err != nil {
		log.Println("[handler:ad][ERROR] ValidateGetAllAdsRequest: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ads, err := h.service.GetAllAds(ctx, &req)
	if err != nil {
		log.Println("[handler:ad][ERROR] GetAllAds: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get ads",
		})
	}

	userId, _ := h.getUserIdFromCtx(ctx)
	adsResponse := make([]*dto.AdResponse, 0, len(ads))

	for _, ad := range ads {
		adsResponse = append(adsResponse, dto.NewAdResponse(ad, userId))
	}

	ctx.JSON(http.StatusOK, adsResponse)

}

func (h *Handler) GetAdByID(ctx *gin.Context) {
	log.Println("[handler:ad] GetAdByID called")

	adId := ctx.Param("id")

	adIdInt, err := strconv.Atoi(adId)

	if err != nil {
		log.Println("[handler:ad][ERROR] strconv.Atoi: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ad id",
		})
		return
	}

	ad, err := h.service.GetAdByID(ctx, adIdInt)

	if err != nil {
		log.Println("[handler:ad][ERROR] GetAdByID: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get ad",
		})
		return
	}

	if ad == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Ad not found",
		})
		return
	}

	userId, _ := h.getUserIdFromCtx(ctx)

	adResponse := dto.NewAdResponse(ad, userId)

	ctx.JSON(http.StatusOK, adResponse)

}

func (h *Handler) getUserIdFromCtx(ctx *gin.Context) (int, error) {
	userIdStr, ok := ctx.Get("UserId")
	if !ok {
		return 0, errors.New("no user in context")
	}
	id, ok := userIdStr.(int)
	if !ok {
		return 0, errors.New("userID has wrong type")
	}
	return id, nil
}
