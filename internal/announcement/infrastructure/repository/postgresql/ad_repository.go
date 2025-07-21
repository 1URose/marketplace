package postgresql

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	entityAF "github.com/1URose/marketplace/internal/announcement/domain/ad_filter/entity"
	"github.com/1URose/marketplace/internal/common/db/postgresql"
	"log"
	"strings"
)

type AdRepository struct {
	Connection *postgresql.Client
}

func NewAdRepository(connection *postgresql.Client) *AdRepository {
	log.Printf("[repository:ad] NewAdRepository initialized")
	return &AdRepository{Connection: connection}
}

func (ar *AdRepository) CreateAd(ctx context.Context, ad *entity.Ad) (*entity.Ad, error) {
	log.Printf("[repository:ad] CreateAd called: title=%q description=%q imageURL=%q price=%d authorID=%d",
		ad.Title, ad.Description, ad.ImageURL, ad.Price, ad.AuthorID,
	)

	const q = `
        INSERT INTO ads (title, description, image_url, price, author_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
	row := ar.Connection.GetPool().QueryRow(ctx, q,
		ad.Title,
		ad.Description,
		ad.ImageURL,
		ad.Price,
		ad.AuthorID,
	)

	if err := row.Scan(&ad.ID, &ad.CreatedAt); err != nil {
		log.Printf("[repository:ad][ERROR] CreateAd scan failed: %v", err)
		return nil, err
	}

	log.Printf("[repository:ad] CreateAd succeeded: adID=%d", ad.ID)

	return ad, nil
}

func (ar *AdRepository) GetAllAds(ctx context.Context, filter *entityAF.AdFilter) ([]*entity.Ad, error) {
	log.Printf("[repository:ad] GetAllAds called: page=%d pageSize=%d sortBy=%s sortOrder=%s minPrice=%v maxPrice=%v",
		filter.Page, filter.PageSize, filter.SortBy, filter.SortOrder, filter.MinPrice, filter.MaxPrice,
	)

	sql, args := ar.buildGetAllAdsQuery(filter)

	rows, err := ar.Connection.GetPool().Query(ctx, sql, args...)
	if err != nil {
		log.Printf("[repository:ad][ERROR] GetAllAds query failed: %v", err)
		return nil, fmt.Errorf("GetAllAds query: %w", err)
	}
	defer rows.Close()

	var ads []*entity.Ad
	for rows.Next() {
		a := new(entity.Ad)
		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.ImageURL,
			&a.Price,
			&a.AuthorID,
			&a.AuthorEmail,
			&a.CreatedAt,
		); err != nil {
			log.Printf("[repository:ad][ERROR] GetAllAds scan failed: %v", err)
			return nil, fmt.Errorf("GetAllAds scan: %w", err)
		}
		ads = append(ads, a)
	}
	log.Printf("[repository:ad] GetAllAds succeeded: returned=%d", len(ads))

	return ads, nil
}

func (ar *AdRepository) buildGetAllAdsQuery(adFilter *entityAF.AdFilter) (string, []interface{}) {
	log.Printf("[repository:ad] buildGetAllAdsQuery called")

	args := make([]interface{}, 0)
	filters := make([]string, 0)

	if adFilter.MinPrice != nil {
		args = append(args, *adFilter.MinPrice)
		filters = append(filters, fmt.Sprintf("price >= $%d", len(args)))
		log.Printf("[repository:ad] buildGetAllAdsQuery: minPrice=%v", adFilter.MinPrice)
	}
	if adFilter.MaxPrice != nil {
		args = append(args, *adFilter.MaxPrice)
		filters = append(filters, fmt.Sprintf("price <= $%d", len(args)))
	}

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString(`
        SELECT 
            a.id, a.title, a.description, a.image_url, a.price, 
            a.author_id, u.email AS author_email, a.created_at
        FROM ads a
        JOIN users u ON a.author_id = u.id
    `)

	if len(filters) > 0 {
		sqlBuilder.WriteString(" WHERE ")
		sqlBuilder.WriteString(strings.Join(filters, " AND "))
	}

	sortBy := adFilter.SortBy
	sortOrder := strings.ToUpper(adFilter.SortOrder)

	sqlBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder))

	page := adFilter.Page

	size := adFilter.PageSize

	offset := (page - 1) * size
	args = append(args, offset, size)
	sqlBuilder.WriteString(fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(args)-1, len(args)))

	log.Printf("[repository:ad] buildGetAllAdsQuery: sql=%s", sqlBuilder.String())

	return sqlBuilder.String(), args

}

func (ar *AdRepository) CountAds(ctx context.Context) (int, error) {
	log.Printf("[repository:ad] CountAds called")
	const q = `
		SELECT count(*) FROM ads
	`
	var count int
	err := ar.Connection.GetPool().
		QueryRow(ctx, q).
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountAds scan: %w", err)
	}
	return count, nil
}
