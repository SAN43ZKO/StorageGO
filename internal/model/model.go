package model

import (
	"errors"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Rubles      int       `json:"Rubles"`
	Pennies     int       `json:"Pennies"`
	Quantity    int       `json:"Quantity"`
	CreatedAt   time.Time `json:"Created_at"`
	UpdatedAt   time.Time `json:"Updated_at"`
}

func (req *Product) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Rubles <= 0 {
		return errors.New("rubles cant be negative or zero")
	}
	if req.Quantity < 0 {
		return errors.New("quantity cant be negative")
	}
	return nil
}
