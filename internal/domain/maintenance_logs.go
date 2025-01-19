package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

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
