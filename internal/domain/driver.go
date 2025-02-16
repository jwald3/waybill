package domain

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	statemachine "github.com/jwald3/lollipop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email string

func NewEmail(email string) (Email, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("invalid email format: %w", err)
	}
	return Email(email), nil
}

type PhoneNumber string

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	if !phoneRegex.MatchString(phone) {
		return "", fmt.Errorf("invalid phone number format")
	}
	return PhoneNumber(phone), nil
}

type EmploymentStatus string

const (
	EmploymentStatusActive     EmploymentStatus = "ACTIVE"
	EmploymentStatusSuspended  EmploymentStatus = "SUSPENDED"
	EmploymentStatusTerminated EmploymentStatus = "TERMINATED"
)

type Driver struct {
	ID                primitive.ObjectID         `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName         string                     `bson:"first_name" json:"first_name"`
	LastName          string                     `bson:"last_name" json:"last_name"`
	DOB               string                     `bson:"dob" json:"dob"`
	LicenseNumber     string                     `bson:"license_number" json:"license_number"`
	LicenseState      string                     `bson:"license_state" json:"license_state"`
	LicenseExpiration string                     `bson:"license_expiration" json:"license_expiration"`
	Phone             PhoneNumber                `bson:"phone" json:"phone"`
	Email             Email                      `bson:"email" json:"email"`
	Address           Address                    `bson:"address" json:"address"`
	EmploymentStatus  EmploymentStatus           `bson:"employment_status" json:"employment_status"`
	CreatedAt         primitive.DateTime         `bson:"created_at" json:"created_at"`
	UpdatedAt         primitive.DateTime         `bson:"updated_at" json:"updated_at"`
	StateMachine      *statemachine.StateMachine `bson:"-" json:"-"`
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

	validEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	validPhone, err := NewPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	driver := &Driver{
		FirstName:         firstName,
		LastName:          lastName,
		DOB:               dateOfBirth,
		LicenseNumber:     licenseNumber,
		LicenseState:      licenseState,
		LicenseExpiration: licenseExpiration,
		Phone:             validPhone,
		Email:             validEmail,
		Address:           address,
		EmploymentStatus:  EmploymentStatusActive,
		CreatedAt:         primitive.NewDateTimeFromTime(now),
		UpdatedAt:         primitive.NewDateTimeFromTime(now),
	}

	if err := driver.InitializeStateMachine(); err != nil {
		return nil, fmt.Errorf("failed to initialize state machine: %w", err)
	}

	return driver, nil
}

func (d *Driver) InitializeStateMachine() error {
	sm := statemachine.NewStateMachine(d.EmploymentStatus)

	sm.AddSimpleTransition(EmploymentStatusActive, EmploymentStatusSuspended)
	sm.AddSimpleTransition(EmploymentStatusActive, EmploymentStatusTerminated)

	sm.AddSimpleTransition(EmploymentStatusSuspended, EmploymentStatusActive)
	sm.AddSimpleTransition(EmploymentStatusSuspended, EmploymentStatusTerminated)

	sm.SetEntryAction(EmploymentStatusActive, func() error {
		d.EmploymentStatus = EmploymentStatusActive
		return nil
	})

	sm.SetEntryAction(EmploymentStatusSuspended, func() error {
		d.EmploymentStatus = EmploymentStatusSuspended
		return nil
	})

	sm.SetEntryAction(EmploymentStatusTerminated, func() error {
		d.EmploymentStatus = EmploymentStatusTerminated
		return nil
	})

	d.StateMachine = sm

	return nil
}

func (d *Driver) SuspendDriver() error {
	if err := d.StateMachine.Transition(EmploymentStatusSuspended); err != nil {
		return fmt.Errorf("failed to suspend driver from status %s: %w", d.EmploymentStatus, err)
	}

	now := time.Now()

	d.UpdatedAt = primitive.NewDateTimeFromTime(now)

	return nil
}

func (d *Driver) TerminateDriver() error {
	if err := d.StateMachine.Transition(EmploymentStatusTerminated); err != nil {
		return fmt.Errorf("failed to terminate driver from status %s: %w", d.EmploymentStatus, err)
	}

	now := time.Now()

	d.UpdatedAt = primitive.NewDateTimeFromTime(now)

	return nil
}

func (d *Driver) ActivateDriver() error {
	if err := d.StateMachine.Transition(EmploymentStatusActive); err != nil {
		return fmt.Errorf("failed to activate driver from status %s: %w", d.EmploymentStatus, err)
	}

	now := time.Now()

	d.UpdatedAt = primitive.NewDateTimeFromTime(now)

	return nil
}
