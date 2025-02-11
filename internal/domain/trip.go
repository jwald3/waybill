package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripStatus string

const (
	TripStatusScheduled      TripStatus = "SCHEDULED"
	TripStatusInTransit      TripStatus = "IN_TRANSIT"
	TripStatusCompleted      TripStatus = "COMPLETED"
	TripStatusFailedDelivery TripStatus = "FAILED_DELIVERY"
	TripStatusCanceled       TripStatus = "CANCELED"
)

func (s TripStatus) IsValid() bool {
	switch s {
	case TripStatusScheduled,
		TripStatusInTransit,
		TripStatusCompleted,
		TripStatusFailedDelivery,
		TripStatusCanceled:
		return true
	}
	return false
}

type Trip struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	TripNumber      string              `bson:"trip_number" json:"trip_number"`
	DriverID        *primitive.ObjectID `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	Driver          *Driver             `bson:"driver,omitempty" json:"driver,omitempty"`
	TruckID         *primitive.ObjectID `bson:"truck_id,omitempty" json:"truck_id,omitempty"`
	Truck           *Truck              `bson:"truck,omitempty" json:"truck,omitempty"`
	StartFacilityID *primitive.ObjectID `bson:"start_facility_id,omitempty" json:"start_facility_id,omitempty"`
	StartFacility   *Facility           `bson:"start_facility,omitempty" json:"start_facility,omitempty"`
	EndFacilityID   *primitive.ObjectID `bson:"end_facility_id,omitempty" json:"end_facility_id,omitempty"`
	EndFacility     *Facility           `bson:"end_facility,omitempty" json:"end_facility,omitempty"`
	DepartureTime   TimeWindow          `bson:"departure_time" json:"departure_time"`
	ArrivalTime     TimeWindow          `bson:"arrival_time" json:"arrival_time"`
	Status          TripStatus          `bson:"status" json:"status"`
	Cargo           Cargo               `bson:"cargo" json:"cargo"`
	FuelUsage       float64             `bson:"fuel_usage_gallons" json:"fuel_usage_gallons"`
	DistanceMiles   int                 `bson:"distance_miles" json:"distance_miles"`
	CreatedAt       primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

type TimeWindow struct {
	Scheduled primitive.DateTime  `bson:"scheduled" json:"scheduled"`
	Actual    *primitive.DateTime `bson:"actual,omitempty" json:"actual,omitempty"`
}

type Cargo struct {
	Description string  `bson:"description" json:"description"`
	Weight      float64 `bson:"weight" json:"weight"`
	Hazmat      bool    `bson:"hazmat" json:"hazmat"`
}

func NewTrip(
	tripNumber string,
	status TripStatus,
	driverId,
	truckId,
	startFacilityID,
	endFacilityID *primitive.ObjectID,
	departureTime,
	arrivalTime TimeWindow,
	cargo Cargo,
	fuelUsage float64,
	distanceMiles int) (*Trip, error) {

	if !status.IsValid() {
		return nil, fmt.Errorf("invalid status provided: %s", status)
	}

	now := time.Now()

	return &Trip{
		TripNumber:      tripNumber,
		DriverID:        driverId,
		TruckID:         truckId,
		StartFacilityID: startFacilityID,
		EndFacilityID:   endFacilityID,
		DepartureTime:   departureTime,
		ArrivalTime:     arrivalTime,
		Status:          status,
		Cargo:           cargo,
		FuelUsage:       fuelUsage,
		DistanceMiles:   distanceMiles,
		CreatedAt:       primitive.NewDateTimeFromTime(now),
		UpdatedAt:       primitive.NewDateTimeFromTime(now),
	}, nil
}
