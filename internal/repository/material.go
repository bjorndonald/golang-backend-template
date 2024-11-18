package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type MaterialRepositoryInterface interface {
	Create(price *models.MaterialPrice) error
	Find() ([]*models.MaterialPrice, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.MaterialPrice, bool, error)
	Save(price *models.MaterialPrice) (*models.MaterialPrice, error)
	Delete(id string) (*models.MaterialPrice, error)
}

type MaterialRepository struct {
	database *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) MaterialRepositoryInterface {
	return &MaterialRepository{
		database: db,
	}
}

func (a *MaterialRepository) Find() ([]*models.MaterialPrice, error) {
	var prices []*models.MaterialPrice
	err := a.database.Raw(`SELECT * FROM material_prices `).Scan(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func (a *MaterialRepository) Delete(id string) (*models.MaterialPrice, error) {
	var price *models.MaterialPrice
	err := a.database.Raw(`DELETE FROM material_prices WHERE id = ?`, id).Scan(&price).Error
	if err != nil {
		return nil, err
	}
	return price, nil
}

func (a *MaterialRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.MaterialPrice
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *MaterialRepository) FindByCondition(condition, value string) (*models.MaterialPrice, bool, error) {
	var loc *models.MaterialPrice
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.ID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *MaterialRepository) Create(price *models.MaterialPrice) error {
	return a.database.Model(&price).Create(price).Error
}

func (a *MaterialRepository) Save(price *models.MaterialPrice) (*models.MaterialPrice, error) {

	txn := a.database.Model(price).Where("id = ?", price.ID).Updates(&price).First(price)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return price, nil
}
