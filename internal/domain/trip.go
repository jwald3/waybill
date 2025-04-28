package domain

import (
	"fmt"
	"strings"
	"time"

	statemachine "github.com/jwald3/lollipop"
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

type Trip struct {
	ID              primitive.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          primitive.ObjectID         `bson:"user_id" json:"user_id"`
	TripNumber      string                     `bson:"trip_number" json:"trip_number"`
	DriverID        *primitive.ObjectID        `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	Driver          *Driver                    `bson:"driver,omitempty" json:"driver,omitempty"`
	TruckID         *primitive.ObjectID        `bson:"truck_id,omitempty" json:"truck_id,omitempty"`
	Truck           *Truck                     `bson:"truck,omitempty" json:"truck,omitempty"`
	StartFacilityID *primitive.ObjectID        `bson:"start_facility_id,omitempty" json:"start_facility_id,omitempty"`
	StartFacility   *Facility                  `bson:"start_facility,omitempty" json:"start_facility,omitempty"`
	EndFacilityID   *primitive.ObjectID        `bson:"end_facility_id,omitempty" json:"end_facility_id,omitempty"`
	EndFacility     *Facility                  `bson:"end_facility,omitempty" json:"end_facility,omitempty"`
	DepartureTime   TimeWindow                 `bson:"departure_time" json:"departure_time"`
	ArrivalTime     TimeWindow                 `bson:"arrival_time" json:"arrival_time"`
	Status          TripStatus                 `bson:"status" json:"status"`
	Cargo           Cargo                      `bson:"cargo" json:"cargo"`
	FuelUsage       float64                    `bson:"fuel_usage_gallons" json:"fuel_usage_gallons"`
	DistanceMiles   int                        `bson:"distance_miles" json:"distance_miles"`
	Notes           []TripNote                 `bson:"notes" json:"notes"`
	CreatedAt       primitive.DateTime         `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime         `bson:"updated_at" json:"updated_at"`
	StateMachine    *statemachine.StateMachine `bson:"-" json:"-"`
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
	userID primitive.ObjectID,
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

	now := time.Now()

	trip := &Trip{
		UserID:          userID,
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
	}

	if err := trip.InitializeStateMachine(); err != nil {
		return nil, fmt.Errorf("failed to initialize state machine: %w", err)
	}

	return trip, nil
}

type TripFilter struct {
	UserID          primitive.ObjectID
	DriverID        *primitive.ObjectID
	TruckID         *primitive.ObjectID
	StartFacilityID *primitive.ObjectID
	EndFacilityID   *primitive.ObjectID
	Limit           int64
	Offset          int64
}

func NewTripFilter() TripFilter {
	return TripFilter{
		Limit:  10,
		Offset: 0,
	}
}

func (t *Trip) InitializeStateMachine() error {
	sm := statemachine.NewStateMachine(t.Status)

	// Add all valid transitions
	sm.AddSimpleTransition(TripStatusScheduled, TripStatusInTransit)
	sm.AddSimpleTransition(TripStatusScheduled, TripStatusCanceled)
	sm.AddSimpleTransition(TripStatusInTransit, TripStatusCompleted)
	sm.AddSimpleTransition(TripStatusInTransit, TripStatusFailedDelivery)

	// Configure entry actions
	sm.SetEntryAction(TripStatusInTransit, func() error {
		t.Status = TripStatusInTransit
		return nil
	})

	sm.SetEntryAction(TripStatusCompleted, func() error {
		t.Status = TripStatusCompleted
		return nil
	})

	sm.SetEntryAction(TripStatusFailedDelivery, func() error {
		t.Status = TripStatusFailedDelivery
		return nil
	})

	sm.SetEntryAction(TripStatusCanceled, func() error {
		t.Status = TripStatusCanceled
		return nil
	})

	t.StateMachine = sm
	return nil
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
