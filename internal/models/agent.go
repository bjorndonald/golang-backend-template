package models

import "github.com/gofrs/uuid"

type UserAgent struct {
	ID     uuid.UUID `json:"id,omitempty"`
	UserID uuid.UUID `json:"user_id,omitempty"`
	Agent  string    `json:"agent"`
}
