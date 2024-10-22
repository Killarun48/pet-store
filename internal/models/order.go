package models

type Order struct {
	ID       int    `json:"id" db:"id" example:"1"`
	PetID    int    `json:"petId" db:"pet_id" example:"1"`
	Quantity int    `json:"quantity" db:"quantity" example:"10"`
	ShipDate string `json:"shipDate" db:"ship_date" example:"2022-01-01T06:29:51.438Z"`
	Status   string `json:"status" db:"status" example:"placed"`
	Complete bool   `json:"complete" db:"complete" example:"true"`
}
