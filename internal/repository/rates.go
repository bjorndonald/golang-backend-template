package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type RatesRepositoryInterface interface {
	Create(page *models.Rates) error
	Find(id uuid.UUID) (*models.Rates, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.Rates, bool, error)
	FindByUserID(value uuid.UUID) ([]*models.Rates, bool, error)
	Save(page *models.Rates) (*models.Rates, error)
}

type RatesRepository struct {
	database *gorm.DB
}

func NewRatesRepository(db *gorm.DB) RatesRepositoryInterface {
	return &RatesRepository{
		database: db,
	}
}

func (a *RatesRepository) Find(id uuid.UUID) (*models.Rates, error) {
	var loc *models.Rates
	err := a.database.First(&loc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func (a *RatesRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.Rates
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *RatesRepository) FindByCondition(condition, value string) (*models.Rates, bool, error) {
	var rates *models.Rates
	err := a.database.Where(condition, value).Find(&rates).Error
	if err != nil {
		return nil, false, err
	}
	if rates.ID.String() != "" {
		return rates, true, nil
	}
	return nil, false, nil
}

func (a *RatesRepository) FindByUserID(user_id uuid.UUID) ([]*models.Rates, bool, error) {
	var rates []*models.Rates
	err := a.database.Raw(`SELECT * FROM user_locations WHERE user_id = ?`, user_id).Scan(&rates).Error
	if err != nil {
		return nil, false, err
	}
	if rates != nil {
		return rates, true, nil
	}
	return nil, false, nil
}

func (a *RatesRepository) Create(spec *models.Rates) error {
	return a.database.Model(&spec).Create(spec).Error
}

func (a *RatesRepository) Save(spec *models.Rates) (*models.Rates, error) {

	txn := a.database.Model(spec).Where("id = ?", spec.ID).Updates(&spec).First(spec)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return spec, nil
}
