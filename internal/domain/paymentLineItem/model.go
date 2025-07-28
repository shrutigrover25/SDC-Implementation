package payment

import (
	"time"
	"github.com/google/uuid"
)

type PaymentLineItem struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	UID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Version      int
	JobUID       uuid.UUID  // Foreign key to specific job version
	TimelogUID   uuid.UUID  // Foreign key to specific timelog version
	Amount       float64
	Status       string     // "not-paid", "paid", "failed", etc.
	ContractorID uuid.UUID
	IssuedAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (PaymentLineItem) TableName() string { return "payment_line_items" }
func (p PaymentLineItem) GetID() string { return p.ID.String() }
func (p PaymentLineItem) GetUID() string { return p.UID.String() }
func (p PaymentLineItem) GetVersion() int { return p.Version }

func (p PaymentLineItem) SetCreatedAt(t time.Time) PaymentLineItem {
	p.CreatedAt = t
	return p
}
func (p PaymentLineItem) SetUpdatedAt(t time.Time) PaymentLineItem {
	p.UpdatedAt = t
	return p
}

func (p PaymentLineItem) CopyForNewVersion() PaymentLineItem {
	return PaymentLineItem{
		ID:           p.ID,
		JobUID:       p.JobUID,
		TimelogUID:   p.TimelogUID,
		Amount:       p.Amount,
		Status:       p.Status,
		ContractorID: p.ContractorID,
		IssuedAt:     p.IssuedAt,
		Version:      p.Version + 1,
		UID:          uuid.New(),
		// CreatedAt and UpdatedAt will be set by the SCD manager
	}
}
