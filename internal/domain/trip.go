package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	TripNumber    string              `bson:"trip_number"`
	DriverID      *primitive.ObjectID `bson:"driver_id"`
	TruckID       *primitive.ObjectID `bson:"truck_id"`
	StartFacility *primitive.ObjectID `bson:"start_facility_id"`
	EndFacility   *primitive.ObjectID `bson:"end_facility_id"`
	Route         Route               `bson:"route"`
	StartTime     primitive.DateTime  `bson:"start_time"`
	EndTime       primitive.DateTime  `bson:"end_time"`
	Status        string              `bson:"status"`
	Cargo         Cargo               `bson:"cargo"`
	FuelUsage     float64             `bson:"fuel_usage_gallons"`
	DistanceMiles int                 `bson:"distance_miles"`
	Incidents     []Incident          `bson:"incidents"`
	CreatedAt     primitive.DateTime  `bson:"created_at"`
	UpdatedAt     primitive.DateTime  `bson:"updated_at"`
}

type Route struct {
	Origin      string   `bson:"origin"`
	Destination string   `bson:"destination"`
	Waypoints   []string `bson:"waypoints"`
}

type Cargo struct {
	Description string  `bson:"description"`
	Weight      float64 `bson:"weight"`
	Hazmat      bool    `bson:"hazmat"`
}

type Incident struct {
	Type        string `bson:"type"`
	Description string `bson:"description"`
	Location    string `bson:"location"`
	Date        string `bson:"date"`
}

func NewTrip(
	tripNumber,
	status string,
	driverId,
	truckId,
	startFacility,
	endFacility *primitive.ObjectID,
	route Route,
	startTime,
	endTime primitive.DateTime,
	cargo Cargo,
	fuelUsage float64,
	distanceMiles int) (*Trip, error) {
	now := time.Now()

	return &Trip{
		TripNumber:    tripNumber,
		DriverID:      driverId,
		TruckID:       truckId,
		StartFacility: startFacility,
		EndFacility:   endFacility,
		Route:         route,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        status,
		Cargo:         cargo,
		FuelUsage:     fuelUsage,
		DistanceMiles: distanceMiles,
		Incidents:     []Incident{},
		CreatedAt:     primitive.NewDateTimeFromTime(now),
		UpdatedAt:     primitive.NewDateTimeFromTime(now),
	}, nil
}
