package repository

import (
	"regexp"
	"strings"

	"gorm.io/gorm"

	"md-note/backend/internal/models"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) List() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) Create(name string) (*models.Tag, error) {
	tag := models.Tag{Name: name, Slug: slugify(name)}
	if err := r.db.Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) Update(id uint64, name string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	tag.Name = name
	tag.Slug = slugify(name)
	if err := r.db.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
