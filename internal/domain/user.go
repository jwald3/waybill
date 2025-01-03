package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password does not meet strength requirements")
)

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  Password
	Status    Status
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// password behavior/logic
type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, ErrWeakPassword
	}

	// perform logic to actually hash the password
	return Password{hash: "hashed_password"}, nil
}

func (p Password) Verify(plaintext string) bool {
	// perform actual verification
	return true
}

func (p Password) Hash() string {
	return p.hash
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Status behavior/logic
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBanned   Status = "banned"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusBanned:
		return true
	}
	return false
}

// user logic
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

type UserEvent struct {
	UserID    int64
	EventType string
	Timestamp time.Time
}

func (u *User) ToActivatedEvent() UserEvent {
	return UserEvent{
		UserID:    u.ID,
		EventType: "user_activated",
		Timestamp: time.Now(),
	}
}
