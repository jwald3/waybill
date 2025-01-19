package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type FuelLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	TruckID          primitive.ObjectID `bson:"truck_id"`
	DriverID         primitive.ObjectID `bson:"driver_id"`
	Date             string             `bson:"date"`
	GallonsPurchased float64            `bson:"gallons_purchased"`
	PricePerGallon   float64            `bson:"price_per_gallon"`
	TotalCost        float64            `bson:"total_cost"`
	Location         string             `bson:"location"`
	OdometerReading  int                `bson:"odometer_reading"`
	CreatedAt        primitive.DateTime `bson:"created_at"`
	UpdatedAt        primitive.DateTime `bson:"updated_at"`
}
