package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelType string

const (
	FuelTypeDiesel                 FuelType = "DIESEL"
	FuelTypeGasoline               FuelType = "GASOLINE"
	FuelTypeCNG                    FuelType = "CNG"
	FuelTypeLNG                    FuelType = "LNG"
	FuelTypeHydrogen               FuelType = "HYDROGEN"
	FuelTypeElectric               FuelType = "ELECTRIC"
	FuelTypeHybridDieselElectric   FuelType = "HYBRID_DIESEL_ELECTRIC"
	FuelTypeHybridGasolineElectric FuelType = "HYBRID_GASOLINE_ELECTRIC"
	FuelTypeBiodiesel              FuelType = "BIODIESEL"
	FuelTypeRenewableDiesel        FuelType = "RENEWABLE_DIESEL"
)

func (f FuelType) IsValid() bool {
	switch f {
	case FuelTypeDiesel,
		FuelTypeGasoline,
		FuelTypeCNG,
		FuelTypeLNG,
		FuelTypeHydrogen,
		FuelTypeElectric,
		FuelTypeHybridDieselElectric,
		FuelTypeHybridGasolineElectric,
		FuelTypeBiodiesel,
		FuelTypeRenewableDiesel:

		return true
	}
	return false
}

type TrailerType string

const (
	TrailerTypeDryVan        TrailerType = "DRY_VAN"
	TrailerTypeRefrigerated  TrailerType = "REFRIGERATED"
	TrailerTypeFlatBed       TrailerType = "FLAT_BED"
	TrailerTypeTanker        TrailerType = "TANKER"
	TrailerTypeAutoCarrier   TrailerType = "AUTO_CARRIER"
	TrailerTypeLiveStock     TrailerType = "LIVE_STOCK"
	TrailerTypeIntermodal    TrailerType = "INTERMODAL"
	TrailerTypeLogging       TrailerType = "LOGGING"
	TrailerTypePneumaticTank TrailerType = "PNEUMATIC_TANK"
)

func (t TrailerType) IsValid() bool {
	switch t {
	case TrailerTypeDryVan,
		TrailerTypeRefrigerated,
		TrailerTypeFlatBed,
		TrailerTypeTanker,
		TrailerTypeAutoCarrier,
		TrailerTypeLiveStock,
		TrailerTypeIntermodal,
		TrailerTypeLogging,
		TrailerTypePneumaticTank:
		return true
	}
	return false
}

type TruckStatus string

const (
	TruckStatusAvailable        TruckStatus = "AVAILABLE"
	TruckStatusInTransit        TruckStatus = "IN_TRANSIT"
	TruckStatusUnderMaintenance TruckStatus = "UNDER_MAINTENANCE"
	TruckStatusRetired          TruckStatus = "RETIRED"
)

func (s TruckStatus) IsValid() bool {
	switch s {
	case TruckStatusAvailable,
		TruckStatusInTransit,
		TruckStatusUnderMaintenance,
		TruckStatusRetired:
		return true
	}
	return false
}

type MaintenanceServiceType string

const (
	ServiceTypeRoutine         MaintenanceServiceType = "ROUTINE_MAINTENANCE"
	ServiceTypeEmergencyRepair MaintenanceServiceType = "EMERGENCY_REPAIR"
)

func (m MaintenanceServiceType) IsValid() bool {
	switch m {
	case ServiceTypeRoutine,
		ServiceTypeEmergencyRepair:
		return true
	}
	return false
}

type Truck struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty"`
	TruckNumber      string              `bson:"truck_number"`
	VIN              string              `bson:"vin"`
	Make             string              `bson:"make"`
	Model            string              `bson:"model"`
	Year             int                 `bson:"year"`
	LicensePlate     LicensePlate        `bson:"license_plate"`
	Mileage          int                 `bson:"mileage"`
	Status           TruckStatus         `bson:"status"`
	AssignedDriverID *primitive.ObjectID `bson:"assigned_driver_id,omitempty"`
	TrailerType      TrailerType         `bson:"trailer_type"`
	CapacityTons     float64             `bson:"capacity_tons"`
	FuelType         FuelType            `bson:"fuel_type"`
	LastMaintenance  string              `bson:"last_maintenance"`
	CreatedAt        primitive.DateTime  `bson:"created_at"`
	UpdatedAt        primitive.DateTime  `bson:"updated_at"`
}

type LicensePlate struct {
	Number string `bson:"number"`
	State  string `bson:"state"`
}

type MaintenanceRecord struct {
	Date        string                 `bson:"date"`
	ServiceType MaintenanceServiceType `bson:"service_type"`
	Notes       string                 `bson:"notes"`
	Cost        float64                `bson:"cost"`
}

func NewTruck(
	truckNumber,
	vin,
	vehicleMake,
	model string,
	status TruckStatus,
	trailerType TrailerType,
	fuelType FuelType,
	LastMaintenance string,
	year,
	mileage int,
	capacityTons float64,
	licensePlate LicensePlate) (*Truck, error) {
	now := time.Now()

	return &Truck{
		TruckNumber:      truckNumber,
		VIN:              vin,
		Make:             vehicleMake,
		Model:            model,
		Year:             year,
		LicensePlate:     licensePlate,
		Mileage:          mileage,
		Status:           status,
		AssignedDriverID: nil,
		TrailerType:      trailerType,
		CapacityTons:     capacityTons,
		FuelType:         fuelType,
		LastMaintenance:  LastMaintenance,
		CreatedAt:        primitive.NewDateTimeFromTime(now),
		UpdatedAt:        primitive.NewDateTimeFromTime(now),
	}, nil
}
