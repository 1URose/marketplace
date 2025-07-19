package validator

import (
	"errors"
	"fmt"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"net/http"
	"strings"
	"time"
)

const (
	MinTitleLen       = 5
	MaxTitleLen       = 100
	MinDescriptionLen = 10
	MaxDescriptionLen = 1000
	MinPrice          = 1
	MaxPrice          = 1_000_000_00
	MaxImageFileSize  = 5 * 1024 * 1024 // 5 МБ
)

func ValidateCreateAd(req dto.CreateAdRequest) error {
	if err := ValidateTitle(req.Title); err != nil {
		return err
	}
	if err := ValidateDescription(req.Description); err != nil {
		return err
	}
	if err := ValidatePrice(req.Price); err != nil {
		return err
	}
	if err := ValidateImageURL(req.ImageURL); err != nil {
		return err
	}
	return nil
}

//TODO переделать константы и допускаемые типы

var allowedImageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/jpg":  {},
}

func ValidateTitle(title string) error {
	ln := len(title)
	if ln < MinTitleLen {
		return fmt.Errorf("title too short: %d < %d", ln, MinTitleLen)
	}
	if ln > MaxTitleLen {
		return fmt.Errorf("title too long: %d > %d", ln, MaxTitleLen)
	}
	return nil
}

func ValidateDescription(desc string) error {
	ln := len(desc)
	if ln < MinDescriptionLen {
		return fmt.Errorf("description too short: %d < %d", ln, MinDescriptionLen)
	}
	if ln > MaxDescriptionLen {
		return fmt.Errorf("description too long: %d > %d", ln, MaxDescriptionLen)
	}
	return nil
}

func ValidatePrice(price int) error {
	if price < MinPrice {
		return fmt.Errorf("price must be at least %d", MinPrice)
	}
	if price > MaxPrice {
		return fmt.Errorf("price must be at most %d", MaxPrice)
	}
	return nil
}

func ValidateImageURL(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("image url must start with http:// or https://")
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		return fmt.Errorf("cannot HEAD image url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}
	if _, ok := allowedImageTypes[ct]; !ok {
		return fmt.Errorf("unsupported content type: %s", ct)
	}

	size := resp.ContentLength
	if size <= 0 {
		return errors.New("content length is missing or zero")
	}
	if size > MaxImageFileSize {
		return fmt.Errorf("image too large: %d > %d", size, MaxImageFileSize)
	}

	return nil
}
