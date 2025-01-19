package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Driver struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	FirstName         string             `bson:"first_name"`
	LastName          string             `bson:"lastname"`
	DOB               string             `bson:"dob"`
	LicenseNumber     string             `bson:"license_number"`
	LicenseState      string             `bson:"license_state"`
	LicenseExpiration string             `bson:"license_expiration"`
	Phone             string             `bson:"phone"`
	Email             string             `bson:"email"`
	Address           Address            `bson:"address"`
	EmploymentStatus  string             `bson:"employmen_status"`
	AssignedTruckID   primitive.ObjectID `bson:"assigned_truck_id,omitempty"`
	PastTrips         []TripHistory      `bson:"past_trips"`
	CreatedAt         primitive.DateTime `bson:"created_at"`
	UpdatedAt         primitive.DateTime `bson:"updated_at"`
}

type Address struct {
	Street string `bson:"street"`
	City   string `bson:"city"`
	State  string `bson:"state"`
	Zip    string `bson:"zip"`
}

type TripHistory struct {
	TripID      primitive.ObjectID `bson:"trip_id"`
	StartDate   string             `bson:"start_date"`
	EndDate     string             `bson:"end_date"`
	MilesDriven int                `bson:"miles_driven"`
}
