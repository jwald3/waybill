package domain

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// commonly used error strings as vars. if you want to reformat the error or add new ones,
// saving as variables will make it much simpler to ensure standardized errors
var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password does not meet strength requirements")
)

// the regex for emails. You can flesh out the logic for what constitutes a valid email if you'd like
var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// the domain model used in the application. This struct defines the information included in the model
// and contains instructions on how to express the model when parsing it into JSON
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  Password
	Status    Status
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// constructor for new users armed with domain-specific logic and checks. You can create constructors for your own resources
// that include similar checks
func NewUser(email, plaintextPassword string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// hash that password using bcrypt
	pass, err := NewPassword(plaintextPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		Email:     email,
		Password:  pass,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// password behavior/logic
type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, ErrWeakPassword
	}

	// use bcrypt for hashing passwords (you're going to want to store the hashed password in your db)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	// perform logic to actually hash the password
	return Password{hash: string(hashedBytes)}, nil
}

func (p Password) Verify(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	return err == nil
}

func (p Password) Hash() string {
	return p.hash
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Status behavior/logic
type Status string

// the values that Status can be. This type of pattern can help make enum-type control flows
// feel free to add additional statuses here (and ensure that the switch statement below reflects it)
const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBanned   Status = "banned"
)

// register valid statuses as part of the `true` returning case
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusBanned:
		return true
	}
	return false
}

// user status functions. You may want to include additional functions if you add extra statuses.
func (u *User) Activate() error {
	if u.Status == StatusBanned {
		return errors.New("cannot activate banned user")
	}
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Ban() error {
	if u.Status == StatusBanned {
		return errors.New("user already banned")
	}
	u.Status = StatusBanned
	u.UpdatedAt = time.Now()
	return nil
}

// state management used when
type UserEvent struct {
	UserID    int64
	EventType string
	Timestamp time.Time
}

// notification functions related to each status. You can use this as part of a wider pub/sub type system if you like
func (u *User) ToActivatedEvent() UserEvent {
	return UserEvent{
		UserID:    u.ID,
		EventType: "user_activated",
		Timestamp: time.Now(),
	}
}

func (u *User) ToBannedEvent() UserEvent {
	return UserEvent{
		UserID:    u.ID,
		EventType: "user_banned",
		Timestamp: time.Now(),
	}
}

func (u *User) ToInactiveEvent() UserEvent {
	return UserEvent{
		UserID:    u.ID,
		EventType: "user_inactive",
		Timestamp: time.Now(),
	}
}
