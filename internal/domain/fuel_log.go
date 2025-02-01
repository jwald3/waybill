package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelLog struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	TruckID          *primitive.ObjectID `bson:"truck_id" json:"truck_id"`
	DriverID         *primitive.ObjectID `bson:"driver_id" json:"driver_id"`
	Date             string              `bson:"date" json:"date"`
	GallonsPurchased float64             `bson:"gallons_purchased" json:"gallons_purchased"`
	PricePerGallon   float64             `bson:"price_per_gallon" json:"price_per_gallon"`
	TotalCost        float64             `bson:"total_cost" json:"total_cost"`
	Location         string              `bson:"location" json:"location"`
	OdometerReading  int                 `bson:"odometer_reading" json:"odometer_reading"`
	CreatedAt        primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt        primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

func NewFuelLog(
	truckId,
	driverId *primitive.ObjectID,
	date,
	location string,
	gallonsPurchased,
	pricePerGallon,
	totalCost float64,
	odometerReading int) (*FuelLog, error) {
	now := time.Now()

	return &FuelLog{
		TruckID:          truckId,
		DriverID:         driverId,
		Date:             date,
		GallonsPurchased: gallonsPurchased,
		PricePerGallon:   pricePerGallon,
		TotalCost:        totalCost,
		Location:         location,
		OdometerReading:  odometerReading,
		CreatedAt:        primitive.NewDateTimeFromTime(now),
		UpdatedAt:        primitive.NewDateTimeFromTime(now),
	}, nil
}
