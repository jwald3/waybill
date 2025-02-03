package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmploymentStatus string

const (
	EmploymentStatusActive     EmploymentStatus = "ACTIVE"
	EmploymentStatusSuspended  EmploymentStatus = "SUSPENDED"
	EmploymentStatusTerminated EmploymentStatus = "TERMINATED"
)

func (e EmploymentStatus) IsValid() bool {
	switch e {
	case EmploymentStatusActive,
		EmploymentStatusSuspended,
		EmploymentStatusTerminated:
		return true
	}
	return false
}

type Driver struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName         string              `bson:"first_name" json:"first_name"`
	LastName          string              `bson:"last_name" json:"last_name"`
	DOB               string              `bson:"dob" json:"dob"`
	LicenseNumber     string              `bson:"license_number" json:"license_number"`
	LicenseState      string              `bson:"license_state" json:"license_state"`
	LicenseExpiration string              `bson:"license_expiration" json:"license_expiration"`
	Phone             string              `bson:"phone" json:"phone"`
	Email             string              `bson:"email" json:"email"`
	Address           Address             `bson:"address" json:"address"`
	EmploymentStatus  EmploymentStatus    `bson:"employment_status" json:"employment_status"`
	AssignedTruckID   *primitive.ObjectID `bson:"assigned_truck_id,omitempty" json:"assigned_truck_id,omitempty"`
	AssignedTruck     *Truck              `bson:"assigned_truck,omitempty" json:"assigned_truck,omitempty"`
	CreatedAt         primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt         primitive.DateTime  `bson:"updated_at" json:"updated_at"`
}

type Address struct {
	Street string `bson:"street" json:"street"`
	City   string `bson:"city" json:"city"`
	State  string `bson:"state" json:"state"`
	Zip    string `bson:"zip" json:"zip"`
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
		EmploymentStatus:  EmploymentStatusActive,
		AssignedTruckID:   nil,
		CreatedAt:         primitive.NewDateTimeFromTime(now),
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}, nil
}

func (d *Driver) ChangeEmploymentStatus(newStatus EmploymentStatus) error {
	if d.EmploymentStatus == EmploymentStatusTerminated {
		return fmt.Errorf("cannot change status from TERMINATED")
	}

	if !newStatus.IsValid() {
		return fmt.Errorf("invalid employment status: %s", newStatus)
	}

	d.EmploymentStatus = newStatus
	d.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	return nil
}
