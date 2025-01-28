package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FacilityService string

const (
	FacilityServiceRepairs          FacilityService = "REPAIRS"
	FacilityServiceLoadingUnloading FacilityService = "LOADING_UNLOADING"
	FacilityServiceLodging          FacilityService = "LODGING"
	FacilityServiceFueling          FacilityService = "FUELING"
)

func (f FacilityService) IsValid() bool {
	switch f {
	case FacilityServiceFueling,
		FacilityServiceLoadingUnloading,
		FacilityServiceLodging,
		FacilityServiceRepairs:
		return true
	}
	return false
}

type Facility struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty"`
	FacilityNumber    string               `bson:"facility_number"`
	Name              string               `bson:"name"`
	Type              string               `bson:"type"`
	Address           Address              `bson:"address"`
	ContactInfo       ContactInfo          `bson:"contact_info"`
	ParkingCapacity   int                  `bson:"parking_capacity"`
	ServicesAvailable []FacilityService    `bson:"services_available"`
	AssignedTrucks    []primitive.ObjectID `bson:"assigned_trucks"`
	CreatedAt         primitive.DateTime   `bson:"created_at"`
	UpdatedAt         primitive.DateTime   `bson:"updated_at"`
}

type ContactInfo struct {
	Phone string `bson:"phone"`
	Email string `bson:"email"`
}

func NewFacility(
	facilityNumber string,
	name string,
	facilityType string,
	address Address,
	contactInfo ContactInfo,
	parkingCapacity int,
	servicesAvailable []FacilityService) (*Facility, error) {
	now := time.Now()

	return &Facility{
		FacilityNumber:    facilityNumber,
		Name:              name,
		Type:              facilityType,
		Address:           address,
		ContactInfo:       contactInfo,
		ParkingCapacity:   parkingCapacity,
		ServicesAvailable: servicesAvailable,
		AssignedTrucks:    []primitive.ObjectID{},
		CreatedAt:         primitive.NewDateTimeFromTime(now),
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}, nil
}
