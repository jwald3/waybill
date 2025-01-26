package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncidentReport struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	Trip           *primitive.ObjectID `bson:"trip_id"`
	TruckID        *primitive.ObjectID `bson:"truck_id"`
	DriverID       *primitive.ObjectID `bson:"driver_id"`
	Type           string              `bson:"type"`
	Description    string              `bson:"description"`
	Date           string              `bson:"date"`
	Location       string              `bson:"location"`
	DamageEstimate float64             `bson:"damage_estimate"`
	CreatedAt      primitive.DateTime  `bson:"created_at"`
	UpdatedAt      primitive.DateTime  `bson:"updated_at"`
}

func NewIncidentReport(trip, truckId, driverId *primitive.ObjectID, incidentType, description, date, location string, damageEstimate float64) (*IncidentReport, error) {
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
