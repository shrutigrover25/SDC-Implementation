package jobs

import (
	"github.com/google/uuid"
)

type Service interface {
	CreateJob(j Job) (Job, error)
	GetByUID(uid string) (Job, error)
	UpdateJob(uid string, title string, rate float64) (Job, error)
	UpdateStatus(uid, status string) (Job, error)
	GetActiveJobsByCompany(companyID string) ([]Job, error)
	GetJobsByContractor(contractorID string) ([]Job, error)
	GetAllActiveJobs() ([]Job, error)
	GetJobHistory(id string) ([]Job, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateJob(j Job) (Job, error) {
	j.Version = 1
	j.UID = uuid.New()
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return s.repo.Create(j)
}

func (s *service) GetByUID(uid string) (Job, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) UpdateJob(uid string, title string, rate float64) (Job, error) {
	return s.repo.Update(uid, func(j Job) Job {
		j.Title = title
		j.Rate = rate
		return j
	})
}

func (s *service) UpdateStatus(uid, status string) (Job, error) {
	return s.repo.UpdateStatus(uid, status)
}

func (s *service) GetActiveJobsByCompany(companyID string) ([]Job, error) {
	id := uuid.MustParse(companyID)
	return s.repo.FindLatestByCompany(id)
}

func (s *service) GetJobsByContractor(contractorID string) ([]Job, error) {
	id := uuid.MustParse(contractorID)
	return s.repo.FindLatestByContractor(id)
}

func (s *service) GetAllActiveJobs() ([]Job, error) {
	return s.repo.FindActiveJobs()
}

func (s *service) GetJobHistory(id string) ([]Job, error) {
	return s.repo.GetVersionHistory(id)
}
