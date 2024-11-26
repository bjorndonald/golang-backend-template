package models

import "github.com/gofrs/uuid"

type UserAgent struct {
	ID          uuid.UUID `json:"id,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	Platform    string    `json:"platform"`
	OS          string    `json:"os"`
	BrowserName string    `json:"browser_name"`
	Mobile      bool      `json:"mobile"`
	Model       string    `json:"model"`
}
