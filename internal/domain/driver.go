package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Driver struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty"`
	FirstName         string              `bson:"first_name"`
	LastName          string              `bson:"last_name"`
	DOB               string              `bson:"dob"`
	LicenseNumber     string              `bson:"license_number"`
	LicenseState      string              `bson:"license_state"`
	LicenseExpiration string              `bson:"license_expiration"`
	Phone             string              `bson:"phone"`
	Email             string              `bson:"email"`
	Address           Address             `bson:"address"`
	EmploymentStatus  string              `bson:"employment_status"`
	AssignedTruckID   *primitive.ObjectID `bson:"assigned_truck_id,omitempty"`
	CreatedAt         primitive.DateTime  `bson:"created_at"`
	UpdatedAt         primitive.DateTime  `bson:"updated_at"`
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

func NewDriver(
	firstName,
	lastName,
	dateOfBirth,
	licenseNumber,
	licenseState,
	licenseExpiration,
	phoneNumber,
	email string,
	address Address) (*Driver, error) {
	now := time.Now()

	return &Driver{
		FirstName:         firstName,
		LastName:          lastName,
		DOB:               dateOfBirth,
		LicenseNumber:     licenseNumber,
		LicenseState:      licenseState,
		LicenseExpiration: licenseExpiration,
		Phone:             phoneNumber,
		Email:             email,
		Address:           address,
		EmploymentStatus:  "active",
		AssignedTruckID:   nil,
		CreatedAt:         primitive.NewDateTimeFromTime(now),
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}, nil
}
