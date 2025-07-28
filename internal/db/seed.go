package db

import (
	"log"
	"time"

	jobs "mercor/internal/domain/jobs"
	paymentLineItem "mercor/internal/domain/paymentLineItem"
	timelog "mercor/internal/domain/timelog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// -------- SEED JOBS ----------
	jobsToSeed := []jobs.Job{
		{
			ID:           uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), // job_ckbk6oo4hn7pacdgcz9f
			Version:      1,
			UID:          uuid.MustParse("00000000-0000-0000-0000-000000000005"), // job_uid_tm15dj18wal295r3xiea
			Status:       "extended",
			Rate:         20.0,
			Title:        "Software Engineer",
			CompanyID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), // comp_cab5i8o0rvh5arskod
			ContractorID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"), // cont_e0nhseq682vkoc4d
		},
		{
			ID:           uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Version:      2,
			UID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"), // job_uid_ae51ppj9jpt56he2ua3
			Status:       "active",
			Rate:         20.0,
			Title:        "Software Engineer",
			CompanyID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			ContractorID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		},
		{
			ID:           uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Version:      3,
			UID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"), // job_uid_ywij5sh1tvfp5nkq7azav
			Status:       "active",
			Rate:         15.5,
			Title:        "Software Engineer",
			CompanyID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			ContractorID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		},
		{
			ID:           uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"), // job_eysl9r8bhyis7y3lgso00
			Version:      1,
			UID:          uuid.MustParse("00000000-0000-0000-0000-000000000004"), // job_uid_c7pnhvtsgcqm15z8pvh
			Status:       "extended",
			Rate:         30.0,
			Title:        "ML Engineer",
			CompanyID:    uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			ContractorID: uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"), // cont_aezrtdqy9kpdvnhuml
		},
	}
	for _, j := range jobsToSeed {
		if err := db.Create(&j).Error; err != nil {
			log.Fatalf("❌ Failed to seed job %v: %v", j.UID, err)
		}
	}

	contractorID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	// ---------- Seed Timelogs (SCD format) ----------
	timelogID := uuid.MustParse("2d30a4b8-983f-4282-8b54-2f82fb70102a")
	timelogs := []timelog.Timelog{
		{
			ID:           timelogID,
			UID:          uuid.MustParse("1c2e2ca7-a69d-421b-b278-f7f83a49e7e5"),
			Version:      1,
			ContractorID: contractorID,
			StartTime:    time.Date(2025, 7, 26, 20, 26, 0, 0, time.UTC),
			EndTime:      time.Date(2025, 7, 26, 21, 26, 0, 0, time.UTC),
		},
		{
			ID:           timelogID,
			UID:          uuid.MustParse("f31a0700-1c48-4813-ae39-c48110143ee3"),
			Version:      2,
			ContractorID: contractorID,
			StartTime:    time.Date(2025, 7, 26, 20, 26, 0, 0, time.UTC),
			EndTime:      time.Date(2025, 7, 26, 21, 56, 0, 0, time.UTC),
		},
	}
	for _, tl := range timelogs {
		if err := db.Create(&tl).Error; err != nil {
			log.Fatalf("❌ Failed to seed timelog: %v", err)
		}
	}

	// ---------- Seed PaymentLineItems (SCD format) ----------
	paymentID := uuid.MustParse("39358d52-2489-4944-a271-ced8d642980d")
	issuedAt := time.Date(2025, 7, 26, 0, 0, 0, 0, time.UTC)

	payments := []paymentLineItem.PaymentLineItem{
		{
			ID:           paymentID,
			UID:          uuid.MustParse("de1dbf39-3e6c-4d3b-af19-4447e2c26571"),
			Version:      1,
			ContractorID: contractorID,
			Amount:       35.0,
			IssuedAt:     issuedAt,
		},
		{
			ID:           paymentID,
			UID:          uuid.MustParse("9cd2d600-49ae-4b68-8b95-e48c3a68f3ea"),
			Version:      2,
			ContractorID: contractorID,
			Amount:       35.0,
			IssuedAt:     issuedAt,
		},
	}
	for _, p := range payments {
		if err := db.Create(&p).Error; err != nil {
			log.Fatalf("❌ Failed to seed payment: %v", err)
		}
	}

	log.Println("✅ SCD Seed completed for Timelogs and PaymentLineItems.")
}