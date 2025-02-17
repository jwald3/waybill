package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaintenanceServiceType string

const (
	ServiceTypeRoutine         MaintenanceServiceType = "ROUTINE_MAINTENANCE"
	ServiceTypeEmergencyRepair MaintenanceServiceType = "EMERGENCY_REPAIR"
)

func (m MaintenanceServiceType) IsValid() bool {
	switch m {
	case ServiceTypeRoutine,
		ServiceTypeEmergencyRepair:
		return true
	}
	return false
}

type MaintenanceLog struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	TruckID     *primitive.ObjectID    `bson:"truck_id,omitempty" json:"truck_id,omitempty"`
	Truck       *Truck                 `bson:"truck,omitempty" json:"truck,omitempty"`
	Date        string                 `bson:"date" json:"date"`
	ServiceType MaintenanceServiceType `bson:"service_type" json:"service_type"`
	Cost        float64                `bson:"cost" json:"cost"`
	Notes       string                 `bson:"notes" json:"notes"`
	Mechanic    string                 `bson:"mechanic" json:"mechanic"`
	Location    string                 `bson:"location" json:"location"`
	CreatedAt   primitive.DateTime     `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime     `bson:"updated_at" json:"updated_at"`
}

func NewMaintenanceLog(
	truckId *primitive.ObjectID,
	date string,
	serviceType MaintenanceServiceType,
	notes,
	mechanic,
	location string,
	cost float64) (*MaintenanceLog, error) {

	if !serviceType.IsValid() {
		return nil, fmt.Errorf("invalid service type provided: %s", serviceType)
	}

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

type MaintenanceLogFilter struct {
	TruckID     *primitive.ObjectID
	ServiceType MaintenanceServiceType
	Limit       int64
	Offset      int64
}

func NewMaintenanceLogFilter() MaintenanceLogFilter {
	return MaintenanceLogFilter{
		Limit:  10,
		Offset: 0,
	}
}
