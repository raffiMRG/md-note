package repository

import (
	"md-note/backend/internal/models"

	"gorm.io/gorm"
)

type CORSRepository struct {
	db *gorm.DB
}

func NewCORSRepository(db *gorm.DB) *CORSRepository {
	return &CORSRepository{db: db}
}

func (r *CORSRepository) List() ([]models.CORSOrigin, error) {
	var rows []models.CORSOrigin
	return rows, r.db.Order("created_at asc").Find(&rows).Error
}

func (r *CORSRepository) Create(origin string) (*models.CORSOrigin, error) {
	row := &models.CORSOrigin{Origin: origin}
	return row, r.db.Create(row).Error
}

func (r *CORSRepository) Delete(id uint64) error {
	return r.db.Delete(&models.CORSOrigin{}, id).Error
}

func (r *CORSRepository) AllOrigins() ([]string, error) {
	var rows []models.CORSOrigin
	if err := r.db.Select("origin").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.Origin
	}
	return out, nil
}
