package jobs

import (
	"time"
	"github.com/google/uuid"
)

type Job struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	UID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Version      int
	Status       string
	Rate         float64
	Title        string
	CompanyID    uuid.UUID
	ContractorID uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Job) TableName() string { return "jobs" }

func (j Job) GetID() string   { return j.ID.String() }
func (j Job) GetUID() string  { return j.UID.String() }
func (j Job) GetVersion() int { return j.Version }

func (j Job) SetCreatedAt(t time.Time) Job {
	j.CreatedAt = t
	return j
}

func (j Job) SetUpdatedAt(t time.Time) Job {
	j.UpdatedAt = t
	return j
}

func (j Job) CopyForNewVersion() Job {
	return Job{
		ID:           j.ID,
		Status:       j.Status,
		Rate:         j.Rate,
		Title:        j.Title,
		CompanyID:    j.CompanyID,
		ContractorID: j.ContractorID,
		Version:      j.Version + 1,
		UID:          uuid.New(),
		// CreatedAt and UpdatedAt will be set by the SCD manager
	}
}
