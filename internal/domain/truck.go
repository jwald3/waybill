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

type Truck struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TruckNumber      string              `bson:"truck_number" json:"truck_number"`
	VIN              string              `bson:"vin" json:"vin"`
	Make             string              `bson:"make" json:"make"`
	Model            string              `bson:"model" json:"model"`
	Year             int                 `bson:"year" json:"year"`
	LicensePlate     LicensePlate        `bson:"license_plate" json:"license_plate"`
	Mileage          int                 `bson:"mileage" json:"mileage"`
	Status           TruckStatus         `bson:"status" json:"status"`
	AssignedDriverID *primitive.ObjectID `bson:"assigned_driver_id,omitempty" json:"assigned_driver_id,omitempty"`
	AssignedDriver   *Driver             `bson:"assigned_driver,omitempty" json:"assigned_driver,omitempty"`
	TrailerType      TrailerType         `bson:"trailer_type" json:"trailer_type"`
	CapacityTons     float64             `bson:"capacity_tons" json:"capacity_tons"`
	FuelType         FuelType            `bson:"fuel_type" json:"fuel_type"`
	LastMaintenance  string              `bson:"last_maintenance" json:"last_maintenance"`
	CreatedAt        primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt        primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

type LicensePlate struct {
	Number string `bson:"number" json:"number"`
	State  string `bson:"state" json:"state"`
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
