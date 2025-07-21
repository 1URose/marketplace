package validator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"github.com/1URose/marketplace/internal/common/config/ad_limits"
)

type AdAllowedValues struct {
	AllowedSortFields map[string]struct{}
	AllowedSortOrders map[string]struct{}

	MinTitleLen       int
	MaxTitleLen       int
	MinDescriptionLen int
	MaxDescriptionLen int
	MinPrice          int
	MaxPrice          int
	MaxImageFileSize  int64
	AllowedImageTypes map[string]struct{}
}

func NewAllowedValues(cfg *ad_limits.AdConfig) *AdAllowedValues {
	log.Printf("[validator:ad] NewAllowedValues called: cfg=%+v", cfg)

	av := &AdAllowedValues{
		AllowedSortFields: make(map[string]struct{}),
		AllowedSortOrders: make(map[string]struct{}),
		MinTitleLen:       cfg.MinTitleLen,
		MaxTitleLen:       cfg.MaxTitleLen,
		MinDescriptionLen: cfg.MinDescriptionLen,
		MaxDescriptionLen: cfg.MaxDescriptionLen,
		MinPrice:          cfg.MinPrice,
		MaxPrice:          cfg.MaxPrice,
		MaxImageFileSize:  cfg.MaxImageFileSize,
		AllowedImageTypes: make(map[string]struct{}),
	}

	for _, f := range strings.Split(cfg.AllowedSortFields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			av.AllowedSortFields[f] = struct{}{}
		}
	}
	for _, o := range strings.Split(cfg.AllowedSortOrders, ",") {
		o = strings.ToLower(strings.TrimSpace(o))
		if o != "" {
			av.AllowedSortOrders[o] = struct{}{}
		}
	}
	for _, ext := range strings.Split(cfg.AllowedImageTypes, ",") {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		clean := strings.TrimPrefix(ext, "image/")
		mimeType := "image/" + clean
		av.AllowedImageTypes[mimeType] = struct{}{}
	}

	log.Printf(
		"[validator:ad] NewAllowedValues initialized: sortFields=%d sortOrders=%d imageTypes=%d titleLen=[%d,%d] descLen=[%d,%d] price=[%d,%d] maxImgSize=%d",
		len(av.AllowedSortFields), len(av.AllowedSortOrders), len(av.AllowedImageTypes),
		av.MinTitleLen, av.MaxTitleLen, av.MinDescriptionLen, av.MaxDescriptionLen,
		av.MinPrice, av.MaxPrice, av.MaxImageFileSize,
	)
	return av
}

func (av *AdAllowedValues) ValidateGetAllAdsRequest(req *dto.GetAllAdsRequest) error {
	log.Printf(
		"[validator:ad] ValidateGetAllAdsRequest called: page=%d sortBy=%q sortOrder=%q minPrice=%v maxPrice=%v",
		req.Page, req.SortBy, req.SortOrder, req.MinPrice, req.MaxPrice,
	)

	if req.Page < 1 {
		err := fmt.Errorf("page must be ≥ 1, got %d", req.Page)
		log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
		return err
	}
	if _, ok := av.AllowedSortFields[req.SortBy]; !ok {
		err := fmt.Errorf("unsupported sort_by: %q", req.SortBy)
		log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
		return err
	}
	if _, ok := av.AllowedSortOrders[req.SortOrder]; !ok {
		err := fmt.Errorf("unsupported sort_order: %q", req.SortOrder)
		log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
		return err
	}
	if req.MinPrice != nil {
		if *req.MinPrice < 0 {
			err := fmt.Errorf("min_price must be ≥ 0, got %d", *req.MinPrice)
			log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
			return err
		}
	}
	if req.MaxPrice != nil {
		if *req.MaxPrice < 0 {
			err := fmt.Errorf("max_price must be ≥ 0, got %d", *req.MaxPrice)
			log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
			return err
		}
	}
	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			err := fmt.Errorf("min_price (%d) cannot be greater than max_price (%d)", *req.MinPrice, *req.MaxPrice)
			log.Printf("[validator:ad][ERROR] ValidateGetAllAdsRequest: %v", err)
			return err
		}
	}

	log.Println("[validator:ad] ValidateGetAllAdsRequest succeeded")
	return nil
}

func (av *AdAllowedValues) ValidateCreateAd(req dto.CreateAdRequest) error {
	log.Printf(
		"[validator:ad] ValidateCreateAd called: titleLen=%d descriptionLen=%d price=%d imageURL=%q",
		len(req.Title), len(req.Description), req.Price, req.ImageURL,
	)

	if err := av.validateTitle(req.Title); err != nil {
		return err
	}
	if err := av.validateDescription(req.Description); err != nil {
		return err
	}
	if err := av.validatePrice(req.Price); err != nil {
		return err
	}
	if err := av.validateImageURL(req.ImageURL); err != nil {
		return err
	}

	log.Println("[validator:ad] ValidateCreateAd succeeded")
	return nil
}

func (av *AdAllowedValues) validateTitle(title string) error {
	ln := len(title)
	log.Printf("[validator:ad] validateTitle: title=%q length=%d", title, ln)

	if ln < av.MinTitleLen {
		err := fmt.Errorf("title too short: %d < %d", ln, av.MinTitleLen)
		log.Printf("[validator:ad][ERROR] validateTitle: %v", err)
		return err
	}
	if ln > av.MaxTitleLen {
		err := fmt.Errorf("title too long: %d > %d", ln, av.MaxTitleLen)
		log.Printf("[validator:ad][ERROR] validateTitle: %v", err)
		return err
	}

	log.Println("[validator:ad] validateTitle succeeded")
	return nil
}

func (av *AdAllowedValues) validateDescription(desc string) error {
	ln := len(desc)
	log.Printf("[validator:ad] validateDescription: length=%d", ln)

	if ln < av.MinDescriptionLen {
		err := fmt.Errorf("description too short: %d < %d", ln, av.MinDescriptionLen)
		log.Printf("[validator:ad][ERROR] validateDescription: %v", err)
		return err
	}
	if ln > av.MaxDescriptionLen {
		err := fmt.Errorf("description too long: %d > %d", ln, av.MaxDescriptionLen)
		log.Printf("[validator:ad][ERROR] validateDescription: %v", err)
		return err
	}

	log.Println("[validator:ad] validateDescription succeeded")
	return nil
}

func (av *AdAllowedValues) validatePrice(price int) error {
	log.Printf("[validator:ad] validatePrice: price=%d", price)

	if price < av.MinPrice {
		err := fmt.Errorf("price must be at least %d", av.MinPrice)
		log.Printf("[validator:ad][ERROR] validatePrice: %v", err)
		return err
	}
	if price > av.MaxPrice {
		err := fmt.Errorf("price must be at most %d", av.MaxPrice)
		log.Printf("[validator:ad][ERROR] validatePrice: %v", err)
		return err
	}

	log.Println("[validator:ad] validatePrice succeeded")
	return nil
}

func (av *AdAllowedValues) validateImageURL(url string) error {
	log.Printf("[validator:ad] validateImageURL called: url=%q", url)

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		err := errors.New("image url must start with http:// or https://")
		log.Printf("[validator:ad][ERROR] validateImageURL: %v", err)
		return err
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		log.Printf("[validator:ad][ERROR] validateImageURL HEAD failed: %v", err)
		return fmt.Errorf("cannot HEAD image url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		log.Printf("[validator:ad][ERROR] validateImageURL: %v", err)
		return err
	}

	ct := resp.Header.Get("Content-Type")
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}
	if _, ok := av.AllowedImageTypes[ct]; !ok {
		err := fmt.Errorf("unsupported content type: %s", ct)
		log.Printf("[validator:ad][ERROR] validateImageURL: %v", err)
		return err
	}

	size := resp.ContentLength
	if size <= 0 {
		err := errors.New("content length is missing or zero")
		log.Printf("[validator:ad][ERROR] validateImageURL: %v", err)
		return err
	}
	if size > av.MaxImageFileSize {
		err := fmt.Errorf("image too large: %d > %d", size, av.MaxImageFileSize)
		log.Printf("[validator:ad][ERROR] validateImageURL: %v", err)
		return err
	}

	log.Printf("[validator:ad] validateImageURL succeeded: contentType=%s size=%d", ct, size)
	return nil
}
