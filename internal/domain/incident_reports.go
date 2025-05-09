package domain

import (
	"fmt"
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
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id"`
	TripID         *primitive.ObjectID `bson:"trip_id,omitempty" json:"trip_id,omitempty"`
	Trip           *Trip               `bson:"trip,omitempty" json:"trip,omitempty"`
	TruckID        *primitive.ObjectID `bson:"truck_id,omitempty" json:"truck_id,omitempty"`
	Truck          *Truck              `bson:"truck,omitempty" json:"truck,omitempty"`
	DriverID       *primitive.ObjectID `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	Driver         *Driver             `bson:"driver,omitempty" json:"driver,omitempty"`
	Type           IncidentType        `bson:"type" json:"type"`
	Description    string              `bson:"description" json:"description"`
	Date           string              `bson:"date" json:"date"`
	Location       string              `bson:"location" json:"location"`
	DamageEstimate float64             `bson:"damage_estimate" json:"damage_estimate"`
	CreatedAt      primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt      primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

func NewIncidentReport(
	userID primitive.ObjectID,
	tripId,
	truckId,
	driverId *primitive.ObjectID,
	incidentType IncidentType,
	description,
	date,
	location string,
	damageEstimate float64) (*IncidentReport, error) {

	if !incidentType.IsValid() {
		return nil, fmt.Errorf("invalid incident report type: %s", incidentType)
	}

	now := time.Now()

	return &IncidentReport{
		UserID:         userID,
		TripID:         tripId,
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

type IncidentReportFilter struct {
	UserID   primitive.ObjectID
	TripID   *primitive.ObjectID
	TruckID  *primitive.ObjectID
	DriverID *primitive.ObjectID
	Type     IncidentType
	Limit    int64
	Offset   int64
}

func NewIncidentReportFilter() IncidentReportFilter {
	return IncidentReportFilter{
		Limit:  10,
		Offset: 0,
		UserID: primitive.NilObjectID,
	}
}
