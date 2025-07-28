package payment

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Create(p PaymentLineItem) (PaymentLineItem, error)
	FindByUID(uid string) (PaymentLineItem, error)
	Update(uid string, updateFn func(PaymentLineItem) PaymentLineItem) (PaymentLineItem, error)
	UpdateStatus(uid string, newStatus string) (PaymentLineItem, error)
	FindLatestByContractor(contractorID uuid.UUID) ([]PaymentLineItem, error)
	FindLatestByContractorInPeriod(contractorID uuid.UUID, start, end time.Time) ([]PaymentLineItem, error)
	FindPendingPayments() ([]PaymentLineItem, error)
	FindByJobUID(jobUID uuid.UUID) ([]PaymentLineItem, error)
	GetVersionHistory(id string) ([]PaymentLineItem, error)
}

type repo struct {
	scd scd.SCDRepository[PaymentLineItem]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[PaymentLineItem](db)}
}

func (r *repo) Create(p PaymentLineItem) (PaymentLineItem, error) {
	return r.scd.Create(p)
}

func (r *repo) FindByUID(uid string) (PaymentLineItem, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, updateFn func(PaymentLineItem) PaymentLineItem) (PaymentLineItem, error) {
	return r.scd.Update(uid, updateFn)
}

func (r *repo) UpdateStatus(uid string, newStatus string) (PaymentLineItem, error) {
	return r.scd.Update(uid, func(p PaymentLineItem) PaymentLineItem {
		p.Status = newStatus
		return p
	})
}

func (r *repo) FindLatestByContractor(contractorID uuid.UUID) ([]PaymentLineItem, error) {
	return r.scd.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		Order("issued_at DESC").
		Find()
}

func (r *repo) FindLatestByContractorInPeriod(contractorID uuid.UUID, start, end time.Time) ([]PaymentLineItem, error) {
	return r.scd.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		BetweenDates(start, end).
		Order("issued_at DESC").
		Find()
}

func (r *repo) FindPendingPayments() ([]PaymentLineItem, error) {
	return r.scd.Query().
		Latest().
		Where("status = ?", "not-paid").
		Order("issued_at ASC").
		Find()
}

func (r *repo) FindByJobUID(jobUID uuid.UUID) ([]PaymentLineItem, error) {
	return r.scd.Query().
		Latest().
		Where("job_uid = ?", jobUID).
		Order("issued_at DESC").
		Find()
}

func (r *repo) GetVersionHistory(id string) ([]PaymentLineItem, error) {
	return r.scd.GetVersionHistory(id)
}
