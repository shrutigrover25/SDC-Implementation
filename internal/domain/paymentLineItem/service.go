package payment

import (
	"time"
	"github.com/google/uuid"
)

type Service interface {
	Create(p PaymentLineItem) (PaymentLineItem, error)
	GetByUID(uid string) (PaymentLineItem, error)
	UpdateStatus(uid string, newStatus string) (PaymentLineItem, error)
	UpdateAmount(uid string, newAmount float64) (PaymentLineItem, error)
	GetByContractor(id string) ([]PaymentLineItem, error)
	GetPendingPayments() ([]PaymentLineItem, error)
	GetVersionHistory(id string) ([]PaymentLineItem, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(p PaymentLineItem) (PaymentLineItem, error) {
	p.UID = uuid.New()
	p.Version = 1
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.IssuedAt.IsZero() {
		p.IssuedAt = time.Now()
	}
	if p.Status == "" {
		p.Status = "not-paid"
	}
	return s.repo.Create(p)
}

func (s *service) GetByUID(uid string) (PaymentLineItem, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) UpdateStatus(uid string, newStatus string) (PaymentLineItem, error) {
	return s.repo.UpdateStatus(uid, newStatus)
}

func (s *service) UpdateAmount(uid string, newAmount float64) (PaymentLineItem, error) {
	return s.repo.Update(uid, func(p PaymentLineItem) PaymentLineItem {
		p.Amount = newAmount
		return p
	})
}

func (s *service) GetByContractor(id string) ([]PaymentLineItem, error) {
	return s.repo.FindLatestByContractor(uuid.MustParse(id))
}

func (s *service) GetPendingPayments() ([]PaymentLineItem, error) {
	return s.repo.FindPendingPayments()
}

func (s *service) GetVersionHistory(id string) ([]PaymentLineItem, error) {
	return s.repo.GetVersionHistory(id)
}
