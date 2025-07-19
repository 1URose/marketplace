package postgresql

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	entityAF "github.com/1URose/marketplace/internal/announcement/domain/ad_filter/entity"
	"github.com/1URose/marketplace/internal/common/db/postgresql"
	"github.com/jackc/pgx/v5"
	"strings"
)

type AdRepository struct {
	Connection *postgresql.Client
}

func NewAdRepository(connection *postgresql.Client) *AdRepository {
	return &AdRepository{Connection: connection}
}

func (ar *AdRepository) CreateAd(ctx context.Context, ad *entity.Ad) (*entity.Ad, error) {
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
		return nil, err
	}

	return ad, nil
}

func (ar *AdRepository) GetAllAds(ctx context.Context, filter *entityAF.AdFilter) ([]*entity.Ad, error) {
	sql, args := ar.buildGetAllAdsQuery(filter)

	rows, err := ar.Connection.GetPool().Query(ctx, sql, args...)
	if err != nil {
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
			&a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetAllAds scan: %w", err)
		}
		ads = append(ads, a)
	}
	return ads, nil
}

func (ar *AdRepository) buildGetAllAdsQuery(adFilter *entityAF.AdFilter) (string, []interface{}) {
	args := make([]interface{}, 0)
	filters := make([]string, 0)

	// Фильтрация по цене
	if adFilter.MinPrice != nil {
		args = append(args, *adFilter.MinPrice)
		filters = append(filters, fmt.Sprintf("price >= $%d", len(args)))
	}
	if adFilter.MaxPrice != nil {
		args = append(args, *adFilter.MaxPrice)
		filters = append(filters, fmt.Sprintf("price <= $%d", len(args)))
	}

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("SELECT id, title, description, image_url, price, author_id, created_at FROM ads")

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

	return sqlBuilder.String(), args

}

func (ar *AdRepository) GetAdByID(ctx context.Context, id int) (*entity.Ad, error) {
	const q = `
        SELECT id, title, description, image_url, price, author_id, created_at
        FROM ads
        WHERE id = $1
    `
	a := new(entity.Ad)
	err := ar.Connection.GetPool().
		QueryRow(ctx, q, id).
		Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.ImageURL,
			&a.Price,
			&a.AuthorID,
			&a.CreatedAt,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetAdByID scan: %w", err)
	}
	return a, nil
}
