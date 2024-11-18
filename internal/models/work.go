package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Page struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	PageTitle string    `json:"pagetitle"`
	SubTitle  string    `json:"subtitle"`
	Section   string    `json:"section"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LabourRate struct {
	ID             uuid.UUID `json:"id,omitempty"`
	WorkIdentifier string    `json:"work_identifier"`
	Title          string    `json:"title"`
	SubTitle       string    `json:"subtitle"`
	Labour         string    `json:"labour"`
	Unit           string    `json:"unit"`
	PriceInNaira   float32   `json:"price_in_naira"`
}

type CurrentPlantRate struct {
	ID             uuid.UUID `json:"id,omitempty"`
	WorkIdentifier string    `json:"work_identifier"`
	Title          string    `json:"title"`
	SubTitle       string    `json:"subtitle"`
	Equipment      string    `json:"equipment"`
	Rental         float32   `json:"rental"`
	Diesel         float32   `json:"diesel"`
	TotalPrice     float32   `json:"total_price"`
}

type MaterialPrice struct {
	ID             uuid.UUID `json:"id,omitempty"`
	WorkIdentifier string    `json:"work_identifier"`
	Title          string    `json:"title"`
	SubTitle       string    `json:"subtitle"`
	Material       string    `json:"material"`
	Unit           string    `json:"unit"`
	Price          float32   `json:"price"`
}

type Specs struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Section    uuid.UUID `json:"section"`
	SubSection string    `json:"subsection"`
	Title      string    `json:"title"`
	Label      string    `json:"label"`
	DataType   string    `json:"data_type"`
	Value      string    `json:"value"`
}

type Section struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Title string    `json:"title"`
}

type Rates struct {
	ID                       uuid.UUID `json:"id,omitempty"`
	Section                  uuid.UUID `json:"section"`
	Title                    string    `json:"title"`
	SubTitle                 string    `json:"subtitle"`
	Description              string    `json:"description"`
	Signature                string    `json:"signature"`
	Label                    string    `json:"label"`
	Unit                     string    `json:"unit"`
	NettRate                 float32   `json:"nett_rate"`
	LabourRates              float32   `json:"labour_rates"`
	IndeginousContractorRate float32   `json:"indeginous_contractor_rate"`
	MediumContractorRate     float32   `json:"medium_contractor_rate"`
	HighContractorRate       float32   `json:"high_contractor_rate"`
}
