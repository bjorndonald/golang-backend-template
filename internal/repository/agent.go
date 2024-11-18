package repository

import (
	"errors"

	"github.com/bjorndonald/lasgcce/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AgentRepositoryInterface interface {
	Create(loc *models.UserAgent) error
	Find(id uuid.UUID) (*models.UserAgent, error)
	Exists(id uuid.UUID) (bool, error)
	FindByCondition(condition, value string) (*models.UserAgent, bool, error)
	FindByUserID(value uuid.UUID) ([]*models.UserAgent, bool, error)
	Save(loc *models.UserAgent) (*models.UserAgent, error)
}

type AgentRepository struct {
	database *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepositoryInterface {
	return &AgentRepository{
		database: db,
	}
}

func (a *AgentRepository) Find(id uuid.UUID) (*models.UserAgent, error) {
	var agent *models.UserAgent
	err := a.database.First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (a *AgentRepository) Exists(id uuid.UUID) (bool, error) {
	var agent *models.UserAgent
	err := a.database.First(&agent, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if agent == nil {
		return false, nil
	}
	return true, nil
}

func (a *AgentRepository) FindByCondition(condition, value string) (*models.UserAgent, bool, error) {
	var loc *models.UserAgent
	err := a.database.Where(condition, value).Find(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc.UserID.String() != "" {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *AgentRepository) FindByUserID(user_id uuid.UUID) ([]*models.UserAgent, bool, error) {
	var loc []*models.UserAgent
	err := a.database.Raw(`SELECT * FROM user_Agens WHERE user_id = ?`, user_id).Scan(&loc).Error
	if err != nil {
		return nil, false, err
	}
	if loc != nil {
		return loc, true, nil
	}
	return nil, false, nil
}

func (a *AgentRepository) Create(agent *models.UserAgent) error {
	return a.database.Model(&agent).Create(agent).Error
}

func (a *AgentRepository) Save(agent *models.UserAgent) (*models.UserAgent, error) {

	txn := a.database.Model(agent).Where("id = ?", agent.ID).Updates(&agent).First(agent)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return agent, nil
}
