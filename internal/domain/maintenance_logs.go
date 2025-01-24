package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaintenanceLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TruckID     primitive.ObjectID `bson:"truck_id"`
	Date        string             `bson:"date"`
	ServiceType string             `bson:"service_type"`
	Cost        float64            `bson:"cost"`
	Notes       string             `bson:"notes"`
	Mechanic    string             `bson:"mechanic"`
	Location    string             `bson:"location"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at"`
}

func NewMaintenanceLog(
	truckId primitive.ObjectID,
	date,
	serviceType,
	notes,
	mechanic,
	location string,
	cost float64) (*MaintenanceLog, error) {

	now := time.Now()

	return &MaintenanceLog{
		TruckID:     truckId,
		Date:        date,
		ServiceType: serviceType,
		Cost:        cost,
		Notes:       notes,
		Mechanic:    mechanic,
		Location:    location,
		CreatedAt:   primitive.NewDateTimeFromTime(now),
		UpdatedAt:   primitive.NewDateTimeFromTime(now),
	}, nil
}
