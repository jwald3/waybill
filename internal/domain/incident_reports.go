package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type IncidentReport struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Trip           primitive.ObjectID `bson:"trip_id"`
	TruckID        primitive.ObjectID `bson:"truck_id"`
	DriverID       primitive.ObjectID `bson:"driver_id"`
	Type           string             `bson:"type"`
	Description    string             `bson:"description"`
	Date           string             `bson:"date"`
	Location       string             `bson:"location"`
	DamageEstimate float64            `bson:"damage_estimate"`
	CreatedAt      primitive.DateTime `bson:"created_at"`
	UpdatedAt      primitive.DateTime `bson:"updated_at"`
}
