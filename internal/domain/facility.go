package domain

import (
	"fmt"
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
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	FacilityNumber    string             `bson:"facility_number" json:"facility_number"`
	Name              string             `bson:"name" json:"name"`
	Type              string             `bson:"type" json:"type"`
	Address           Address            `bson:"address" json:"address"`
	ContactInfo       ContactInfo        `bson:"contact_info" json:"contact_info"`
	ParkingCapacity   int                `bson:"parking_capacity" json:"parking_capacity"`
	ServicesAvailable []FacilityService  `bson:"services_available" json:"services_available"`
	CreatedAt         primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt         primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type ContactInfo struct {
	Phone string `bson:"phone" json:"phone"`
	Email string `bson:"email" json:"email"`
}

func NewFacility(
	userID primitive.ObjectID,
	facilityNumber string,
	name string,
	facilityType string,
	address Address,
	contactInfo ContactInfo,
	parkingCapacity int,
	servicesAvailable []FacilityService) (*Facility, error) {
	now := time.Now()

	for _, service := range servicesAvailable {
		if !service.IsValid() {
			return nil, fmt.Errorf("invalid facility service: %s", service)
		}
	}

	return &Facility{
		UserID:            userID,
		FacilityNumber:    facilityNumber,
		Name:              name,
		Type:              facilityType,
		Address:           address,
		ContactInfo:       contactInfo,
		ParkingCapacity:   parkingCapacity,
		ServicesAvailable: servicesAvailable,
		CreatedAt:         primitive.NewDateTimeFromTime(now),
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}, nil
}

type FacilityFilter struct {
	UserID          primitive.ObjectID
	StateCode       string
	Type            string
	ServicesInclude []FacilityService
	MinCapacity     *int
	MaxCapacity     *int
	Limit           int64
	Offset          int64
}

// we're going to set a limit of 10 and an offset of 0 by default, but the actual values will be set by the query params
// in the handler. If we want to add any other defaults, we can do so here. Stuff like only showing active facilities, etc.
func NewFacilityFilter() FacilityFilter {
	return FacilityFilter{
		Limit:  10,
		Offset: 0,
		UserID: primitive.NilObjectID,
	}
}

func (f *Facility) ChangeServicesAvailable(servicesAvailable []FacilityService) error {
	for _, service := range servicesAvailable {
		if !service.IsValid() {
			return fmt.Errorf("invalid facility service: %s", service)
		}
	}

	f.ServicesAvailable = servicesAvailable
	f.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	return nil
}
