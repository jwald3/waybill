package domain

import (
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
	DriverID        *primitive.ObjectID `bson:"driver_id" json:"driver_id"`
	Driver          *Driver             `bson:"driver" json:"driver"`
	TruckID         *primitive.ObjectID `bson:"truck_id" json:"truck_id"`
	Truck           *Truck              `bson:"truck" json:"truck"`
	StartFacilityID *primitive.ObjectID `bson:"start_facility_id" json:"start_facility_id"`
	StartFacility   *Facility           `bson:"start_facility" json:"start_facility"`
	EndFacilityID   *primitive.ObjectID `bson:"end_facility_id" json:"end_facility_id"`
	EndFacility     *Facility           `bson:"end_facility" json:"end_facility"`
	Route           Route               `bson:"route" json:"route"`
	StartTime       primitive.DateTime  `bson:"start_time" json:"start_time"`
	EndTime         primitive.DateTime  `bson:"end_time" json:"end_time"`
	Status          TripStatus          `bson:"status" json:"status"`
	Cargo           Cargo               `bson:"cargo" json:"cargo"`
	FuelUsage       float64             `bson:"fuel_usage_gallons" json:"fuel_usage_gallons"`
	DistanceMiles   int                 `bson:"distance_miles" json:"distance_miles"`
	CreatedAt       primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

type Route struct {
	Origin      string   `bson:"origin" json:"origin"`
	Destination string   `bson:"destination" json:"destination"`
	Waypoints   []string `bson:"waypoints" json:"waypoints"`
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
	route Route,
	startTime,
	endTime primitive.DateTime,
	cargo Cargo,
	fuelUsage float64,
	distanceMiles int) (*Trip, error) {
	now := time.Now()

	return &Trip{
		TripNumber:      tripNumber,
		DriverID:        driverId,
		TruckID:         truckId,
		StartFacilityID: startFacilityID,
		EndFacilityID:   endFacilityID,
		Route:           route,
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          status,
		Cargo:           cargo,
		FuelUsage:       fuelUsage,
		DistanceMiles:   distanceMiles,
		CreatedAt:       primitive.NewDateTimeFromTime(now),
		UpdatedAt:       primitive.NewDateTimeFromTime(now),
	}, nil
}
