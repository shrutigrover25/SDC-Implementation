package jobs

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Create(job Job) (Job, error)
	FindByUID(uid string) (Job, error)
	Update(uid string, updateFn func(Job) Job) (Job, error)
	UpdateStatus(uid string, newStatus string) (Job, error)
	FindLatestByCompany(companyID uuid.UUID) ([]Job, error)
	FindLatestByContractor(contractorID uuid.UUID) ([]Job, error)
	FindActiveJobs() ([]Job, error)
	GetVersionHistory(id string) ([]Job, error)
}

type repo struct {
	scd scd.SCDRepository[Job]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[Job](db)}
}

func (r *repo) Create(j Job) (Job, error) {
	return r.scd.Create(j)
}

func (r *repo) FindByUID(uid string) (Job, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, updateFn func(Job) Job) (Job, error) {
	return r.scd.Update(uid, updateFn)
}

func (r *repo) UpdateStatus(uid string, newStatus string) (Job, error) {
	return r.scd.Update(uid, func(j Job) Job {
		j.Status = newStatus
		return j
	})
}

func (r *repo) FindLatestByCompany(companyID uuid.UUID) ([]Job, error) {
	return r.scd.Query().
		Latest().
		Where("company_id = ?", companyID).
		Where("status = ?", "active").
		Order("created_at DESC").
		Find()
}

func (r *repo) FindLatestByContractor(contractorID uuid.UUID) ([]Job, error) {
	return r.scd.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		Where("status IN ?", []string{"active", "extended"}).
		Order("created_at DESC").
		Find()
}

func (r *repo) FindActiveJobs() ([]Job, error) {
	return r.scd.Query().
		Latest().
		Where("status = ?", "active").
		Order("title ASC").
		Find()
}

func (r *repo) GetVersionHistory(id string) ([]Job, error) {
	return r.scd.GetVersionHistory(id)
}
