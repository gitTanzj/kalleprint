package types

import (
	"time"
)

type PostOrderPayload struct {
}

type Customer struct {
	Id         string
	Name       string
	Address    string
	Created_at time.Time
}

type Order struct {
	Id          string    `json:"id"`
	CustomerId string    `json:"customer_id"`
	PrintId string `json:"print_id"`
	Total       float64   `json:"total"`
	CreatedAt   time.Time `json:"created_at"`
}

type Print struct {
	Id                 string  `json:"id"`
	FilamentName       string  `json:"filament_name"`
	PathToFile         string  `json:"path_to_file"`
	FilamentAmountUsed float64 `json:"filament_amount_used"`
	EstimatedPrintTime string  `json:"estimated_print_time"`
}

type Filament struct {
	Id            string  `json:"id"`
	FilamentName  string  `json:"filament_name"`
	CurrentAmount int     `json:"current_amount"`
	CostPerGram   float64 `json:"cost_per_gram"`
}

type OrderRepository interface {
	GetOrders() ([]*Order, error)
}

type FilamentRepository interface {
	GetFilaments() ([]*Filament, error)
	GetFilamentByID(id string) (*Filament, error)
}

type PrintRepository interface {
	GetPrints() ([]*Print, error)
}
