package main

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateOnly represents a date without time
type DateOnly struct {
	time.Time
}

// Scan implements the Scanner interface for database reads
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	
	switch v := value.(type) {
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
	case time.Time:
		d.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
	default:
		return fmt.Errorf("cannot scan %T into DateOnly", value)
	}
	return nil
}

// Value implements the driver Valuer interface for database writes
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}

// MarshalJSON implements json.Marshaler
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02"))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time = time.Time{}
		return nil
	}
	
	str := string(data[1 : len(data)-1]) // Remove quotes
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// Status enum for job application status
type Status string

const (
	SUBMITTED        Status = "SUBMITTED"
	REJECTED         Status = "REJECTED"
	PHONE_SCREEN     Status = "PHONE_SCREEN"
	REMOTE_INTERVIEW Status = "REMOTE_INTERVIEW"
	ON_SITE_INTERVIEW Status = "ON_SITE_INTERVIEW"
)

// JobApplication represents a job application
type JobApplication struct {
	AppId         uint      `gorm:"primaryKey;autoIncrement" json:"appId"`
	Company       string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_company_position_date" json:"company"`
	Position      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_company_position_date" json:"position"`
	Location      string    `gorm:"type:varchar(255)" json:"location"`
	SalaryRange   string    `gorm:"type:varchar(100)" json:"salaryRange"`
	WorkplaceType string    `gorm:"type:varchar(50)" json:"workplaceType"`
	Status        Status    `gorm:"type:varchar(50);default:SUBMITTED" json:"status"`
	Notes         string    `gorm:"type:text" json:"notes"`
	Website       string    `gorm:"type:varchar(500)" json:"website"`
	DateApplied   DateOnly  `gorm:"type:varchar(10);uniqueIndex:idx_company_position_date" json:"dateApplied"`
}

// TableName specifies the table name for GORM
func (JobApplication) TableName() string {
	return "apps"
}
