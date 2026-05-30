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
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	ShippingAddress string    `json:"shipping_address"`
	Notes           string    `json:"notes"`
	PrintId         string    `json:"print_id"`
	Total           float64   `json:"total"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type AdminUser struct {
	Id           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type DashboardStats struct {
	TotalOrders   int     `json:"total_orders"`
	Pending       int     `json:"pending"`
	InProgress    int     `json:"in_progress"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type CreateFilamentPayload struct {
	Name        string  `json:"name"`
	AmountGrams int     `json:"amount_grams"`
	TotalPrice  float64 `json:"total_price"`
}

type AdminRepository interface {
	GetAdminByEmail(email string) (*AdminUser, error)
	GetOrders(status string) ([]*Order, error)
	GetOrderByID(id string) (*Order, error)
	UpdateOrderStatus(id, status string) error
	GetFilaments() ([]*Filament, error)
	CreateFilament(name string, amountGrams int, costPerGram float64) (*Filament, error)
	UpdateFilament(id, name string, amountGrams int, costPerGram float64) (*Filament, error)
	DeleteFilament(id string) error
	GetDashboardStats() (*DashboardStats, error)
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
	ColorHex      string  `json:"color_hex"`
	IniFilePath   string  `json:"ini_file_path"`

	// Slicer config — populated from the ini file, not stored in the DB.
	FilamentType              string `json:"filament_type"`
	Temperature               string `json:"temperature"`
	BedTemperature            string `json:"bed_temperature"`
	FirstLayerTemperature     string `json:"first_layer_temperature"`
	FirstLayerBedTemperature  string `json:"first_layer_bed_temperature"`
	FilamentDiameter          string `json:"filament_diameter"`
	ExtrusionMultiplier       string `json:"extrusion_multiplier"`
	FilamentDensity           string `json:"filament_density"`
	FilamentCost              string `json:"filament_cost"`
}

type OrderRepository interface {
	GetOrders() ([]*Order, error)
	CreateAndReturnOrder(name, email, shippingAddress, notes, printID string, total float64) (*Order, error)
}

type FilamentRepository interface {
	GetFilaments() ([]*Filament, error)
	GetFilamentByID(id string) (*Filament, error)
	CreateFilament(id, name string, currentAmount int, costPerGram float64, colorHex, iniFilePath string) (*Filament, error)
	UpdateFilament(id string, name *string, currentAmount *int, costPerGram *float64, colorHex *string, iniFilePath *string) error
	DeleteFilament(id string) error
}

type PrintRepository interface {
	GetPrints() ([]*Print, error)
}
