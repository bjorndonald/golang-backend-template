package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PageRepositoryInterface interface {
	Create(page *models.Page) error
	Find() ([]*models.Page, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.Page, bool, error)
	Delete(id string) (*models.Page, error)
	Save(page *models.Page) (*models.Page, error)
}

type PageRepository struct {
	database *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepositoryInterface {
	return &PageRepository{
		database: db,
	}
}

func (a *PageRepository) Find() ([]*models.Page, error) {
	var pages []*models.Page
	err := a.database.Raw(`SELECT * FROM pages `).Scan(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (a *PageRepository) Exists(id uuid.UUID) (bool, error) {
	var location *models.Page
	err := a.database.First(&location, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if location == nil {
		return false, nil
	}
	return true, nil
}

func (a *PageRepository) FindByCondition(condition, value string) (*models.Page, bool, error) {
	var loc *models.Page
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.ID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *PageRepository) Delete(id string) (*models.Page, error) {
	var page *models.Page
	err := a.database.Raw(`DELETE FROM pages WHERE id = ?`, id).Scan(&page).Error
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (a *PageRepository) Create(loc *models.Page) error {
	return a.database.Model(&loc).Create(loc).Error
}

func (a *PageRepository) Save(loc *models.Page) (*models.Page, error) {

	txn := a.database.Model(loc).Where("id = ?", loc.ID).Updates(&loc).First(loc)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return loc, nil
}
