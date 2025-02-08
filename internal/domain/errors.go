package domain

import "errors"

var ErrDriverNotFound = errors.New("driver not found")
var ErrFacilityNotFound = errors.New("facility not found")
var ErrFuelLogNotFound = errors.New("fuel log not found")
var ErrIncidentReportNotFound = errors.New("incident report not found")
var ErrMaintenanceLogNotFound = errors.New("maintenance log not found")
var ErrTripNotFound = errors.New("trip not found")
var ErrTruckNotFound = errors.New("truck not found")
