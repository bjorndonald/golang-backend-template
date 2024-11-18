package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PlantRepositoryInterface interface {
	Create(page *models.CurrentPlantRate) error
	Find() ([]*models.CurrentPlantRate, error)
	Exists(id uuid.UUID) (bool, error)
	Delete(id uuid.UUID) (*models.CurrentPlantRate, error)
	FindByCondition(condition, value string) (*models.CurrentPlantRate, bool, error)
	Save(page *models.CurrentPlantRate) (*models.CurrentPlantRate, error)
}

type PlantRepository struct {
	database *gorm.DB
}

func NewPlantRepository(db *gorm.DB) PlantRepositoryInterface {
	return &PlantRepository{
		database: db,
	}
}

func (a *PlantRepository) Delete(id uuid.UUID) (*models.CurrentPlantRate, error) {
	var plant *models.CurrentPlantRate
	err := a.database.Raw(`DELETE FROM plant_rates WHERE id = ?`, id).Scan(&plant).Error
	if err != nil {
		return nil, err
	}
	return plant, nil
}

func (a *PlantRepository) Find() ([]*models.CurrentPlantRate, error) {
	var rates []*models.CurrentPlantRate
	err := a.database.Raw(`SELECT * FROM plant_rates `).Scan(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (a *PlantRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.CurrentPlantRate
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *PlantRepository) FindByCondition(condition, value string) (*models.CurrentPlantRate, bool, error) {
	var loc *models.CurrentPlantRate
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.ID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *PlantRepository) Create(loc *models.CurrentPlantRate) error {
	return a.database.Model(&loc).Create(loc).Error
}

func (a *PlantRepository) Save(loc *models.CurrentPlantRate) (*models.CurrentPlantRate, error) {

	txn := a.database.Model(loc).Where("id = ?", loc.ID).Updates(&loc).First(loc)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return loc, nil
}
