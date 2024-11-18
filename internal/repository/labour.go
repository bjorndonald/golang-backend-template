package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type LabourRepositoryInterface interface {
	Create(page *models.LabourRate) error
	Find() ([]*models.LabourRate, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.LabourRate, bool, error)
	Save(page *models.LabourRate) (*models.LabourRate, error)
	Delete(id string) (*models.LabourRate, error)
}

type LabourRepository struct {
	database *gorm.DB
}

func NewLabourRepository(db *gorm.DB) LabourRepositoryInterface {
	return &LabourRepository{
		database: db,
	}
}

func (a *LabourRepository) Delete(id string) (*models.LabourRate, error) {
	var labour *models.LabourRate
	err := a.database.Raw(`DELETE FROM labour_rates WHERE id = ?`, id).Scan(&labour).Error
	if err != nil {
		return nil, err
	}
	return labour, nil
}

func (a *LabourRepository) Find() ([]*models.LabourRate, error) {
	var rates []*models.LabourRate
	err := a.database.Raw(`SELECT * FROM labour_rates `).Scan(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (a *LabourRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.LabourRate
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *LabourRepository) FindByCondition(condition, value string) (*models.LabourRate, bool, error) {
	var loc *models.LabourRate
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.ID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *LabourRepository) Create(loc *models.LabourRate) error {
	return a.database.Model(&loc).Create(loc).Error
}

func (a *LabourRepository) Save(loc *models.LabourRate) (*models.LabourRate, error) {

	txn := a.database.Model(loc).Where("id = ?", loc.ID).Updates(&loc).First(loc)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return loc, nil
}
