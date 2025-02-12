package domain

import (
	"errors"
	"fmt"
)

var ErrDriverNotFound = errors.New("driver not found")
var ErrFacilityNotFound = errors.New("facility not found")
var ErrFuelLogNotFound = errors.New("fuel log not found")
var ErrIncidentReportNotFound = errors.New("incident report not found")
var ErrMaintenanceLogNotFound = errors.New("maintenance log not found")
var ErrTripNotFound = errors.New("trip not found")
var ErrTruckNotFound = errors.New("truck not found")

type TripStateError struct {
	CurrentState TripStatus
	DesiredState TripStatus
}

func (e *TripStateError) Error() string {
	return fmt.Sprintf("invalid state transition from %s to %s", e.CurrentState, e.DesiredState)
}

type TruckStateError struct {
	CurrentState TruckStatus
	DesiredState TruckStatus
}

func (e *TruckStateError) Error() string {
	return fmt.Sprintf("invalid state transition from %s to %s", e.CurrentState, e.DesiredState)
}
