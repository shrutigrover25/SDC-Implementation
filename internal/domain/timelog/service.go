package timelog

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(t Timelog) (Timelog, error)
	GetByUID(uid string) (Timelog, error)
	UpdateDuration(uid string, duration int64) (Timelog, error)
	UpdateType(uid string, newType string) (Timelog, error)
	GetByContractor(id string) ([]Timelog, error)
	GetVersionHistory(id string) ([]Timelog, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(t Timelog) (Timelog, error) {
	t.UID = uuid.New()
	t.Version = 1
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return s.repo.Create(t)
}

func (s *service) GetByUID(uid string) (Timelog, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) UpdateDuration(uid string, duration int64) (Timelog, error) {
	return s.repo.Update(uid, func(t Timelog) Timelog {
		t.Duration = duration
		t.Type = "adjusted"
		return t
	})
}

func (s *service) UpdateType(uid string, newType string) (Timelog, error) {
	return s.repo.Update(uid, func(t Timelog) Timelog {
		t.Type = newType
		return t
	})
}

func (s *service) GetByContractor(id string) ([]Timelog, error) {
	return s.repo.FindLatestByContractor(uuid.MustParse(id))
}

func (s *service) GetVersionHistory(id string) ([]Timelog, error) {
	return s.repo.GetVersionHistory(id)
}
