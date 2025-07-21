package ad_limits

import (
	"github.com/1URose/marketplace/internal/common/settings"
	"log"
)

type AdConfig struct {
	AllowedSortFields string
	AllowedSortOrders string

	PageSize int

	MinTitleLen       int
	MaxTitleLen       int
	MinDescriptionLen int
	MaxDescriptionLen int
	MinPrice          int
	MaxPrice          int
	MaxImageFileSize  int64
	AllowedImageTypes string
}

func NewAdConfig(
	sortFields, sortOrders string,
	pageSize int,
	minTitle, maxTitle, minDesc, maxDesc, minPrice, maxPrice int,
	maxImgSize64 int64,
	imgTypes string,
) *AdConfig {
	return &AdConfig{
		AllowedSortFields: sortFields,
		AllowedSortOrders: sortOrders,
		PageSize:          pageSize,
		MinTitleLen:       minTitle,
		MaxTitleLen:       maxTitle,
		MinDescriptionLen: minDesc,
		MaxDescriptionLen: maxDesc,
		MinPrice:          minPrice,
		MaxPrice:          maxPrice,
		MaxImageFileSize:  maxImgSize64,
		AllowedImageTypes: imgTypes,
	}
}

func LoadAdConfigFromEnv() *AdConfig {
	const (
		envSortFields      = "ADS_ALLOWED_SORT_FIELDS"
		envSortOrders      = "ADS_ALLOWED_SORT_ORDERS"
		envPageSize        = "ADS_PAGE_SIZE"
		envMinTitleLen     = "ADS_MIN_TITLE_LEN"
		envMaxTitleLen     = "ADS_MAX_TITLE_LEN"
		envMinDescription  = "ADS_MIN_DESC_LEN"
		envMaxDescription  = "ADS_MAX_DESC_LEN"
		envMinPrice        = "ADS_MIN_PRICE"
		envMaxPrice        = "ADS_MAX_PRICE"
		envMaxImageSize    = "ADS_MAX_IMAGE_SIZE"
		envAllowedImgTypes = "ADS_ALLOWED_IMAGE_TYPES"
	)

	sortFields := settings.GetEnvSrt(envSortFields)
	sortOrders := settings.GetEnvSrt(envSortOrders)
	imgTypes := settings.GetEnvSrt(envAllowedImgTypes)

	pageSize, err := settings.GetEnvInt(envPageSize)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envPageSize, err)
	}

	minTitle, err := settings.GetEnvInt(envMinTitleLen)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMinTitleLen, err)
	}
	maxTitle, err := settings.GetEnvInt(envMaxTitleLen)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMaxTitleLen, err)
	}

	minDesc, err := settings.GetEnvInt(envMinDescription)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMinDescription, err)
	}
	maxDesc, err := settings.GetEnvInt(envMaxDescription)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMaxDescription, err)
	}

	minPrice, err := settings.GetEnvInt(envMinPrice)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMinPrice, err)
	}
	maxPrice, err := settings.GetEnvInt(envMaxPrice)
	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMaxPrice, err)
	}

	maxImgSize, err := settings.GetEnvInt(envMaxImageSize)

	if err != nil {
		log.Panicf("[ad_limits][FATAL] invalid %s: %v", envMaxImageSize, err)
	}
	maxImgSize64 := int64(maxImgSize)

	ac := NewAdConfig(
		sortFields,
		sortOrders,
		pageSize,
		minTitle,
		maxTitle,
		minDesc,
		maxDesc,
		minPrice,
		maxPrice,
		maxImgSize64,
		imgTypes,
	)

	log.Printf(
		"[ad_limits:config] loaded: sortFields=%s sortOrders=%s pageSize=%d minTitle=%d maxTitle=%d minDesc=%d maxDesc=%d minPrice=%d maxPrice=%d maxImgSize=%d imgTypes=%s",
		sortFields,
		sortOrders,
		pageSize,
		minTitle,
		maxTitle,
		minDesc,
		maxDesc,
		minPrice,
		maxPrice,
		maxImgSize,
		imgTypes,
	)

	return ac
}
