package repository

import (
	"gorm.io/gorm"

	"md-note/backend/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Order("id asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateRole(id uint64, role string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}

func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}
