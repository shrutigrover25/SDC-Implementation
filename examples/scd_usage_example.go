package examples

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"mercor/internal/scd"
	"mercor/internal/domain/jobs"
	"mercor/internal/domain/timelog"
	"mercor/internal/domain/paymentLineItem"
)

// This file demonstrates how to use the new SCD abstraction for various common queries

// Example 1: Get all active Jobs for a company (latest version is active)
func GetActiveJobsForCompany(db *gorm.DB, companyID uuid.UUID) ([]jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Query().
		Latest().
		Where("company_id = ?", companyID).
		Where("status = ?", "active").
		Order("created_at DESC").
		Find()
}

// Example 2: Get all active Jobs for a contractor (latest version is active)
func GetActiveJobsForContractor(db *gorm.DB, contractorID uuid.UUID) ([]jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		Where("status IN ?", []string{"active", "extended"}).
		Order("title ASC").
		Find()
}

// Example 3: Get all PaymentLineItems for a contractor in a particular period (latest versions only)
func GetPaymentLineItemsForContractorInPeriod(db *gorm.DB, contractorID uuid.UUID, start, end time.Time) ([]payment.PaymentLineItem, error) {
	paymentRepo := scd.NewManager[payment.PaymentLineItem](db)
	
	return paymentRepo.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		BetweenDates(start, end).
		Order("issued_at DESC").
		Find()
}

// Example 4: Get all Timelogs for a contractor in a particular period (latest versions only)
func GetTimelogsForContractorInPeriod(db *gorm.DB, contractorID uuid.UUID, start, end time.Time) ([]timelog.Timelog, error) {
	timelogRepo := scd.NewManager[timelog.Timelog](db)
	
	return timelogRepo.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		BetweenDates(start, end).
		Where("type = ?", "captured").
		Order("start_time DESC").
		Find()
}

// Example 5: Update a Job (creates new version automatically)
func UpdateJobRate(db *gorm.DB, jobUID string, newRate float64) (jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Update(jobUID, func(j jobs.Job) jobs.Job {
		j.Rate = newRate
		return j
	})
}

// Example 6: Update Payment Status (creates new version automatically)
func MarkPaymentAsPaid(db *gorm.DB, paymentUID string) (payment.PaymentLineItem, error) {
	paymentRepo := scd.NewManager[payment.PaymentLineItem](db)
	
	return paymentRepo.Update(paymentUID, func(p payment.PaymentLineItem) payment.PaymentLineItem {
		p.Status = "paid"
		return p
	})
}

// Example 7: Get version history for audit trails
func GetJobVersionHistory(db *gorm.DB, jobID string) ([]jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	return jobRepo.GetVersionHistory(jobID)
}

// Example 8: Point-in-time query - Get state as of a specific date
func GetJobStateAsOfDate(db *gorm.DB, jobID string, date time.Time) (jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	return jobRepo.GetVersionAt(jobID, date)
}

// Example 9: Complex query with multiple conditions
func GetHighValueActiveJobs(db *gorm.DB, minRate float64) ([]jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Query().
		Latest().
		Where("status = ?", "active").
		Where("rate >= ?", minRate).
		Order("rate DESC").
		Limit(50).
		Find()
}

// Example 10: Batch operations
func UpdateMultipleJobStatuses(db *gorm.DB, updates map[string]string) error {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	updateFunctions := make(map[string]func(jobs.Job) jobs.Job)
	for uid, newStatus := range updates {
		status := newStatus // Capture in closure
		updateFunctions[uid] = func(j jobs.Job) jobs.Job {
			j.Status = status
			return j
		}
	}
	
	return jobRepo.UpdateBatch(updateFunctions)
}

// Example 11: Count queries
func CountActiveJobsForCompany(db *gorm.DB, companyID uuid.UUID) (int64, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Query().
		Latest().
		Where("company_id = ?", companyID).
		Where("status = ?", "active").
		Count()
}

// Example 12: Using Raw GORM access for complex joins
func GetJobsWithTimelogSummary(db *gorm.DB, contractorID uuid.UUID) (*gorm.DB, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	latestJobsQuery := jobRepo.Query().
		Latest().
		Where("contractor_id = ?", contractorID).
		Raw()
	
	return latestJobsQuery.
		Select("jobs.*, COUNT(timelogs.uid) as timelog_count").
		Joins("LEFT JOIN timelogs ON jobs.uid = timelogs.job_uid").
		Group("jobs.uid"), nil
}

// Example 13: Creating entities with proper SCD setup
func CreateNewJob(db *gorm.DB, title string, rate float64, companyID, contractorID uuid.UUID) (jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	job := jobs.Job{
		ID:           uuid.New(),
		UID:          uuid.New(),
		Version:      1,
		Title:        title,
		Rate:         rate,
		Status:       "active",
		CompanyID:    companyID,
		ContractorID: contractorID,
	}
	
	return jobRepo.Create(job)
}

// Example 14: Handling foreign key relationships to specific versions
func CreatePaymentLineItem(db *gorm.DB, amount float64, contractorID uuid.UUID, jobUID, timelogUID uuid.UUID) (payment.PaymentLineItem, error) {
	paymentRepo := scd.NewManager[payment.PaymentLineItem](db)
	
	payment := payment.PaymentLineItem{
		ID:           uuid.New(),
		UID:          uuid.New(),
		Version:      1,
		Amount:       amount,
		Status:       "not-paid",
		ContractorID: contractorID,
		JobUID:       jobUID,
		TimelogUID:   timelogUID,
		IssuedAt:     time.Now(),
	}
	
	return paymentRepo.Create(payment)
}

// Example 15: Advanced filtering with WhereIn
func GetJobsForMultipleCompanies(db *gorm.DB, companyIDs []uuid.UUID) ([]jobs.Job, error) {
	jobRepo := scd.NewManager[jobs.Job](db)
	
	return jobRepo.Query().
		Latest().
		WhereIn("company_id", companyIDs).
		Where("status = ?", "active").
		Order("company_id, title").
		Find()
}
