package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type SectionRepositoryInterface interface {
	Create(page *models.Section) error
	Find(id uuid.UUID) (*models.Section, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.Section, bool, error)
	Save(page *models.Section) (*models.Section, error)
}

type SectionRepository struct {
	database *gorm.DB
}

func NewSectionRepository(db *gorm.DB) SectionRepositoryInterface {
	return &SectionRepository{
		database: db,
	}
}

func (a *SectionRepository) Find(id uuid.UUID) (*models.Section, error) {
	var section *models.Section
	err := a.database.First(&section, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return section, nil
}

func (a *SectionRepository) Exists(id uuid.UUID) (bool, error) {
	var section *models.Section
	err := a.database.First(&section, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if section == nil {
		return false, nil
	}
	return true, nil
}

func (a *SectionRepository) FindByCondition(condition, value string) (*models.Section, bool, error) {
	var section *models.Section
	err := a.database.Where(condition, value).Find(&section).Error
	if err != nil {
		return nil, false, err
	}
	if section.ID.String() != "" {
		return section, true, nil
	}
	return nil, false, nil
}

func (a *SectionRepository) Create(section *models.Section) error {
	return a.database.Model(&section).Create(section).Error
}

func (a *SectionRepository) Save(section *models.Section) (*models.Section, error) {

	txn := a.database.Model(section).Where("id = ?", section.ID).Updates(&section).First(section)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return section, nil
}
