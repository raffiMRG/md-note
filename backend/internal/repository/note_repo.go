package repository

import (
	"gorm.io/gorm"

	"md-note/backend/internal/models"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(note *models.Note, tagIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(note).Error; err != nil {
			return err
		}
		return replaceTags(tx, note, tagIDs)
	})
}

func (r *NoteRepository) FindByID(id uint64) (*models.Note, error) {
	var note models.Note
	if err := r.db.Preload("Tags").Preload("CreatedByUser").Preload("UpdatedByUser").First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) List(tagSlug string, page, limit int) ([]models.Note, int64, error) {
	base := r.db.Model(&models.Note{})
	if tagSlug != "" {
		base = base.Joins("JOIN note_tags ON note_tags.note_id = notes.id").
			Joins("JOIN tags ON tags.id = note_tags.tag_id AND tags.slug = ?", tagSlug)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notes []models.Note
	offset := (page - 1) * limit
	err := base.Preload("Tags").Preload("CreatedByUser").Preload("UpdatedByUser").
		Order("notes.updated_at DESC").Limit(limit).Offset(offset).Find(&notes).Error
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepository) Search(query string, page, limit int) ([]models.Note, int64, error) {
	matchExpr := "MATCH(notes.title, notes.content) AGAINST (? IN NATURAL LANGUAGE MODE)"

	var total int64
	if err := r.db.Model(&models.Note{}).Where(matchExpr, query).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notes []models.Note
	offset := (page - 1) * limit
	err := r.db.Preload("Tags").Preload("CreatedByUser").Preload("UpdatedByUser").
		Where(matchExpr, query).
		Order("notes.updated_at DESC").Limit(limit).Offset(offset).
		Find(&notes).Error
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepository) Update(note *models.Note, tagIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Use a column-scoped Updates instead of Save: note was loaded via FindByID with
		// CreatedByUser/UpdatedByUser preloaded, and Save() would auto-sync updated_by back
		// to the stale preloaded association's ID, clobbering the new editor's user ID.
		updates := map[string]interface{}{
			"title":   note.Title,
			"content": note.Content,
		}
		if note.UpdatedBy != nil {
			updates["updated_by"] = *note.UpdatedBy
		}
		if err := tx.Model(&models.Note{}).Where("id = ?", note.ID).Updates(updates).Error; err != nil {
			return err
		}
		return replaceTags(tx, note, tagIDs)
	})
}

func (r *NoteRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Note{}, id).Error
}

func replaceTags(tx *gorm.DB, note *models.Note, tagIDs []uint64) error {
	var tags []models.Tag
	if len(tagIDs) > 0 {
		if err := tx.Find(&tags, tagIDs).Error; err != nil {
			return err
		}
	}
	return tx.Model(note).Association("Tags").Replace(tags)
}
