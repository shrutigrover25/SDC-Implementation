package timelog

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Create(t Timelog) (Timelog, error)
	FindByUID(uid string) (Timelog, error)
	Update(uid string, updateFn func(Timelog) Timelog) (Timelog, error)
	FindLatestByContractor(contractorID uuid.UUID) ([]Timelog, error)
	FindLatestByContractorInPeriod(contractorID uuid.UUID, start, end time.Time) ([]Timelog, error)
	FindByJobUID(jobUID uuid.UUID) ([]Timelog, error)
	GetVersionHistory(id string) ([]Timelog, error)
}

type repo struct {
	scd scd.SCDRepository[Timelog]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[Timelog](db)}
}

func (r *repo) Create(t Timelog) (Timelog, error) {
	return r.scd.Create(t)
}

func (r *repo) FindByUID(uid string) (Timelog, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, updateFn func(Timelog) Timelog) (Timelog, error) {
	return r.scd.Update(uid, updateFn)
}

func (r *repo) FindLatestByContractor(contractorID uuid.UUID) ([]Timelog, error) {
	return r.scd.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		Order("start_time DESC").
		Find()
}

func (r *repo) FindLatestByContractorInPeriod(contractorID uuid.UUID, start, end time.Time) ([]Timelog, error) {
	return r.scd.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		BetweenDates(start, end).
		Order("start_time DESC").
		Find()
}

func (r *repo) FindByJobUID(jobUID uuid.UUID) ([]Timelog, error) {
	return r.scd.Query().
		Latest().
		Where("job_uid = ?", jobUID).
		Order("start_time DESC").
		Find()
}

func (r *repo) GetVersionHistory(id string) ([]Timelog, error) {
	return r.scd.GetVersionHistory(id)
}
