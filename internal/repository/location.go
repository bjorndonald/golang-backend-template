package repository

import (
	"errors"

	"github.com/bjorndonald/golang-backend-template/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type LocationRepositoryInterface interface {
	Create(loc *models.GeoLocation) error
	Find(id uuid.UUID) (*models.GeoLocation, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.GeoLocation, bool, error)
	FindByUserID(value uuid.UUID) ([]*models.GeoLocation, bool, error)
	Save(loc *models.GeoLocation) (*models.GeoLocation, error)
}

type LocationRepository struct {
	database *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepositoryInterface {
	return &LocationRepository{
		database: db,
	}
}

func (a *LocationRepository) Find(id uuid.UUID) (*models.GeoLocation, error) {
	var loc *models.GeoLocation
	err := a.database.First(&loc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func (a *LocationRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.GeoLocation
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *LocationRepository) FindByCondition(condition, value string) (*models.GeoLocation, bool, error) {
	var loc *models.GeoLocation
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.UserID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *LocationRepository) FindByUserID(user_id uuid.UUID) ([]*models.GeoLocation, bool, error) {
	var loc []*models.GeoLocation
	err := a.database.Raw(`SELECT * FROM user_locations WHERE user_id = ?`, user_id).Scan(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc != nil {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *LocationRepository) Create(loc *models.GeoLocation) error {
	return a.database.Model(&loc).Create(loc).Error
}

func (a *LocationRepository) Save(loc *models.GeoLocation) (*models.GeoLocation, error) {

	txn := a.database.Model(loc).Where("id = ?", loc.ID).Updates(&loc).First(loc)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return loc, nil
}
