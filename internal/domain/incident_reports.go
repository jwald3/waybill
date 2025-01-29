package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncidentType string

const (
	IncidentTypeMechanicalFailure IncidentType = "MECHANICAL_FAILURE"
	IncidentTypeTrafficAccident   IncidentType = "TRAFFIC_ACCIDENT"
	IncidentTypeCargoDamage       IncidentType = "CARGO_DAMAGE"
	IncidentTypeTheft             IncidentType = "THEFT"
	IncidentTypeWeatherDelay      IncidentType = "WEATHER_DELAY"
	IncidentTypeRouteDeviation    IncidentType = "ROUTE_DEVIATION"
	IncidentTypeFuelShortage      IncidentType = "FUEL_SHORTAGE"
	IncidentTypeDriverIllness     IncidentType = "DRIVER_ILLNESS"
)

func (i IncidentType) IsValid() bool {
	switch i {
	case IncidentTypeMechanicalFailure,
		IncidentTypeTrafficAccident,
		IncidentTypeCargoDamage,
		IncidentTypeTheft,
		IncidentTypeWeatherDelay,
		IncidentTypeRouteDeviation,
		IncidentTypeFuelShortage,
		IncidentTypeDriverIllness:
		return true
	}
	return false
}

type IncidentReport struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	Trip           *primitive.ObjectID `bson:"trip_id"`
	TruckID        *primitive.ObjectID `bson:"truck_id"`
	DriverID       *primitive.ObjectID `bson:"driver_id"`
	Type           IncidentType        `bson:"type"`
	Description    string              `bson:"description"`
	Date           string              `bson:"date"`
	Location       string              `bson:"location"`
	DamageEstimate float64             `bson:"damage_estimate"`
	CreatedAt      primitive.DateTime  `bson:"created_at"`
	UpdatedAt      primitive.DateTime  `bson:"updated_at"`
}

func NewIncidentReport(
	trip,
	truckId,
	driverId *primitive.ObjectID,
	incidentType IncidentType,
	description,
	date,
	location string,
	damageEstimate float64) (*IncidentReport, error) {
	now := time.Now()

	return &IncidentReport{
		Trip:           trip,
		TruckID:        truckId,
		DriverID:       driverId,
		Type:           incidentType,
		Description:    description,
		Date:           date,
		Location:       location,
		DamageEstimate: damageEstimate,
		CreatedAt:      primitive.NewDateTimeFromTime(now),
		UpdatedAt:      primitive.NewDateTimeFromTime(now),
	}, nil
}
