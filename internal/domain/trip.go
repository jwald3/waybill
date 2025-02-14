package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/jwald3/lollipop/pkg/statemachine"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripStatus string

const (
	TripStatusScheduled      TripStatus = "SCHEDULED"
	TripStatusInTransit      TripStatus = "IN_TRANSIT"
	TripStatusCompleted      TripStatus = "COMPLETED"
	TripStatusFailedDelivery TripStatus = "FAILED_DELIVERY"
	TripStatusCanceled       TripStatus = "CANCELED"
	MaxNoteLength                       = 1000
)

func (s TripStatus) IsValid() bool {
	switch s {
	case TripStatusScheduled,
		TripStatusInTransit,
		TripStatusCompleted,
		TripStatusFailedDelivery,
		TripStatusCanceled:
		return true
	}
	return false
}

type Trip struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	TripNumber      string              `bson:"trip_number" json:"trip_number"`
	DriverID        *primitive.ObjectID `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	Driver          *Driver             `bson:"driver,omitempty" json:"driver,omitempty"`
	TruckID         *primitive.ObjectID `bson:"truck_id,omitempty" json:"truck_id,omitempty"`
	Truck           *Truck              `bson:"truck,omitempty" json:"truck,omitempty"`
	StartFacilityID *primitive.ObjectID `bson:"start_facility_id,omitempty" json:"start_facility_id,omitempty"`
	StartFacility   *Facility           `bson:"start_facility,omitempty" json:"start_facility,omitempty"`
	EndFacilityID   *primitive.ObjectID `bson:"end_facility_id,omitempty" json:"end_facility_id,omitempty"`
	EndFacility     *Facility           `bson:"end_facility,omitempty" json:"end_facility,omitempty"`
	DepartureTime   TimeWindow          `bson:"departure_time" json:"departure_time"`
	ArrivalTime     TimeWindow          `bson:"arrival_time" json:"arrival_time"`
	Status          TripStatus          `bson:"status" json:"status"`
	Cargo           Cargo               `bson:"cargo" json:"cargo"`
	FuelUsage       float64             `bson:"fuel_usage_gallons" json:"fuel_usage_gallons"`
	DistanceMiles   int                 `bson:"distance_miles" json:"distance_miles"`
	Notes           []TripNote          `bson:"notes" json:"notes"`
	CreatedAt       primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime  `bson:"updated_at" json:"updated_at"`
	StateMachine    *statemachine.StateMachine
}

type TimeWindow struct {
	Scheduled primitive.DateTime  `bson:"scheduled" json:"scheduled"`
	Actual    *primitive.DateTime `bson:"actual,omitempty" json:"actual,omitempty"`
}

type TripNote struct {
	NoteTimestamp time.Time `bson:"note_timestamp" json:"note_timestamp"`
	Content       string    `json:"content" bson:"content"`
}

type Cargo struct {
	Description string  `bson:"description" json:"description"`
	Weight      float64 `bson:"weight" json:"weight"`
	Hazmat      bool    `bson:"hazmat" json:"hazmat"`
}

func NewTrip(
	tripNumber string,
	driverId,
	truckId,
	startFacilityID,
	endFacilityID *primitive.ObjectID,
	departureTime,
	arrivalTime TimeWindow,
	cargo Cargo,
	fuelUsage float64,
	distanceMiles int) (*Trip, error) {

	sm := statemachine.NewStateMachine(TripStatusScheduled)

	// Add all valid transitions
	sm.AddTransition(TripStatusScheduled, TripStatusInTransit)
	sm.AddTransition(TripStatusScheduled, TripStatusCanceled)
	sm.AddTransition(TripStatusInTransit, TripStatusCompleted)
	sm.AddTransition(TripStatusInTransit, TripStatusFailedDelivery)

	now := time.Now()

	trip := &Trip{
		TripNumber:      tripNumber,
		DriverID:        driverId,
		TruckID:         truckId,
		StartFacilityID: startFacilityID,
		EndFacilityID:   endFacilityID,
		DepartureTime:   departureTime,
		ArrivalTime:     arrivalTime,
		Status:          TripStatusScheduled,
		Cargo:           cargo,
		FuelUsage:       fuelUsage,
		DistanceMiles:   distanceMiles,
		Notes:           make([]TripNote, 0),
		CreatedAt:       primitive.NewDateTimeFromTime(now),
		UpdatedAt:       primitive.NewDateTimeFromTime(now),
		StateMachine:    sm,
	}

	sm.SetEntryAction(TripStatusInTransit, func() error {
		trip.Status = TripStatusInTransit
		return nil
	})

	sm.SetEntryAction(TripStatusCompleted, func() error {
		trip.Status = TripStatusCompleted
		return nil
	})

	sm.SetEntryAction(TripStatusFailedDelivery, func() error {
		trip.Status = TripStatusFailedDelivery
		return nil
	})

	sm.SetEntryAction(TripStatusCanceled, func() error {
		trip.Status = TripStatusCanceled
		return nil
	})

	return trip, nil
}

func (t *Trip) BeginTrip(departureTime time.Time) error {
	if err := t.StateMachine.Transition(TripStatusInTransit); err != nil {
		return fmt.Errorf("failed to start trip from status %s: %w", t.Status, err)
	}

	departure := primitive.NewDateTimeFromTime(departureTime)

	t.DepartureTime = TimeWindow{
		Scheduled: t.DepartureTime.Scheduled,
		Actual:    &departure,
	}

	t.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	return nil
}

func (t *Trip) CancelTrip() error {
	if err := t.StateMachine.Transition(TripStatusCanceled); err != nil {
		return fmt.Errorf("failed to cancel trip from status %s: %w", t.Status, err)
	}

	t.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	return nil
}

func (t *Trip) CompleteTripSuccessfully(arrivalTime time.Time) error {
	if err := t.StateMachine.Transition(TripStatusCompleted); err != nil {
		return fmt.Errorf("failed to complete trip from status %s: %w", t.Status, err)
	}

	now := time.Now()

	arrival := primitive.NewDateTimeFromTime(arrivalTime)
	t.ArrivalTime = TimeWindow{
		Scheduled: t.ArrivalTime.Scheduled,
		Actual:    &arrival,
	}
	t.UpdatedAt = primitive.NewDateTimeFromTime(now)
	return nil
}

func (t *Trip) CompleteTripUnsuccessfully(arrivalTime time.Time) error {
	if err := t.StateMachine.Transition(TripStatusFailedDelivery); err != nil {
		return fmt.Errorf("failed to mark trip as failed delivery from status %s: %w", t.Status, err)
	}

	now := time.Now()

	arrival := primitive.NewDateTimeFromTime(arrivalTime)
	t.ArrivalTime = TimeWindow{
		Scheduled: t.ArrivalTime.Scheduled,
		Actual:    &arrival,
	}
	t.UpdatedAt = primitive.NewDateTimeFromTime(now)
	return nil
}

func (t *Trip) AddNote(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("note content cannot be empty")
	}
	if len(content) > MaxNoteLength {
		return fmt.Errorf("note content exceeds maximum length of %d characters", MaxNoteLength)
	}

	note := TripNote{
		NoteTimestamp: time.Now(),
		Content:       content,
	}
	t.Notes = append(t.Notes, note)
	return nil
}
