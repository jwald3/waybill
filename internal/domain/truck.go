package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Truck struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty"`
	TruckNumber        string              `bson:"truck_number"`
	VIN                string              `bson:"vin"`
	Make               string              `bson:"make"`
	Model              string              `bson:"model"`
	Year               int                 `bson:"year"`
	LicensePlate       LicensePlate        `bson:"license_plate"`
	Mileage            int                 `bson:"mileage"`
	Status             string              `bson:"status"`
	AssignedDriverID   primitive.ObjectID  `bson:"assigned_driver_id,omitempty"`
	TrailerType        string              `bson:"trailer_type"`
	CapacityTons       float64             `bson:"capacity_tons"`
	FuelType           string              `bson:"fuel_type"`
	LastMaintenance    string              `bson:"last_maintenance"`
	MaintenanceRecords []MaintenanceRecord `bson:"maintenance_records"`
	CreatedAt          primitive.DateTime  `bson:"created_at"`
	UpdatedAt          primitive.DateTime  `bson:"updated_at"`
}

type LicensePlate struct {
	Number string `bson:"number"`
	State  string `bson:"state"`
}

type MaintenanceRecord struct {
	Date        string  `bson:"date"`
	ServiceType string  `bson:"service_type"`
	Notes       string  `bson:"notes"`
	Cost        float64 `bson:"cost"`
}

func NewTruck(
	truckNumber,
	vin,
	vehicleMake,
	model,
	status,
	trailerType,
	fuelType,
	LastMaintenance string,
	year,
	mileage int,
	capacityTons float64,
	assignedDriverID primitive.ObjectID,
	licensePlate LicensePlate) (*Truck, error) {
	now := time.Now()

	return &Truck{
		TruckNumber:        truckNumber,
		VIN:                vin,
		Make:               vehicleMake,
		Model:              model,
		Year:               year,
		LicensePlate:       licensePlate,
		Mileage:            mileage,
		Status:             status,
		AssignedDriverID:   assignedDriverID,
		TrailerType:        trailerType,
		CapacityTons:       capacityTons,
		FuelType:           fuelType,
		LastMaintenance:    LastMaintenance,
		MaintenanceRecords: make([]MaintenanceRecord, 0),
		CreatedAt:          primitive.NewDateTimeFromTime(now),
		UpdatedAt:          primitive.NewDateTimeFromTime(now),
	}, nil
}
