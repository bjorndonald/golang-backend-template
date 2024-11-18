package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type SpecsRepositoryInterface interface {
	Create(page *models.Specs) error
	Find(id uuid.UUID) (*models.Specs, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.Specs, bool, error)
	FindByUserID(value uuid.UUID) ([]*models.Specs, bool, error)
	Save(page *models.Specs) (*models.Specs, error)
}

type SpecsRepository struct {
	database *gorm.DB
}

func NewSpecsRepository(db *gorm.DB) SpecsRepositoryInterface {
	return &SpecsRepository{
		database: db,
	}
}

func (a *SpecsRepository) Find(id uuid.UUID) (*models.Specs, error) {
	var loc *models.Specs
	err := a.database.First(&loc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func (a *SpecsRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.Specs
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *SpecsRepository) FindByCondition(condition, value string) (*models.Specs, bool, error) {
	var specs *models.Specs
	err := a.database.Where(condition, value).Find(&specs).Error
	if err != nil {
		return nil, false, err
	}
	if specs.ID.String() != "" {
		return specs, true, nil
	}
	return nil, false, nil
}

func (a *SpecsRepository) FindByUserID(user_id uuid.UUID) ([]*models.Specs, bool, error) {
	var specs []*models.Specs
	err := a.database.Raw(`SELECT * FROM user_locations WHERE user_id = ?`, user_id).Scan(&specs).Error
	if err != nil {
		return nil, false, err
	}
	if specs != nil {
		return specs, true, nil
	}
	return nil, false, nil
}

func (a *SpecsRepository) Create(spec *models.Specs) error {
	return a.database.Model(&spec).Create(spec).Error
}

func (a *SpecsRepository) Save(spec *models.Specs) (*models.Specs, error) {

	txn := a.database.Model(spec).Where("id = ?", spec.ID).Updates(&spec).First(spec)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return spec, nil
}
