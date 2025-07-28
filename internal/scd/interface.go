package scd

import (
	"time"
	"gorm.io/gorm"
)

// SCDModel defines the contract for SCD entities
type SCDModel[T any] interface {
	TableName() string
	GetID() string
	GetUID() string
	GetVersion() int
	CopyForNewVersion() T
	SetCreatedAt(time.Time) T
	SetUpdatedAt(time.Time) T
}

// SCDQueryBuilder provides a fluent interface for building SCD-aware queries
type SCDQueryBuilder[T SCDModel[T]] interface {
	// Latest version queries
	Latest() SCDQueryBuilder[T]
	Where(query interface{}, args ...interface{}) SCDQueryBuilder[T]
	WhereIn(column string, values interface{}) SCDQueryBuilder[T]
	Order(value interface{}) SCDQueryBuilder[T]
	Limit(limit int) SCDQueryBuilder[T]
	Offset(offset int) SCDQueryBuilder[T]

	// Relationship joins for SCD tables
	JoinLatest(table string, condition string) SCDQueryBuilder[T]

	// Time-based queries
	AsOfDate(date time.Time) SCDQueryBuilder[T]
	BetweenDates(start, end time.Time) SCDQueryBuilder[T]

	// Execution methods
	Find() ([]T, error)
	First() (T, error)
	Count() (int64, error)

	// Advanced queries
	WithHistory() SCDQueryBuilder[T]
	GroupByID() SCDQueryBuilder[T]

	// Raw GORM access when needed
	Raw() *gorm.DB
}

// SCDRepository provides high-level SCD operations
type SCDRepository[T SCDModel[T]] interface {
	// Basic CRUD
	Create(entity T) (T, error)
	FindByUID(uid string) (T, error)
	Update(uid string, updateFn func(T) T) (T, error)
	Delete(uid string) error

	// Query building
	Query() SCDQueryBuilder[T]

	// Batch operations
	CreateBatch(entities []T) error
	UpdateBatch(updates map[string]func(T) T) error

	// Utility methods
	GetLatestVersion(id string) (T, error)
	GetVersionHistory(id string) ([]T, error)
	GetVersionAt(id string, date time.Time) (T, error)
}

// SCDRelationship handles foreign key relationships to specific SCD versions
type SCDRelationship interface {
	// For referencing specific versions of related entities
	ResolveUID(uid string) (interface{}, error)
	ValidateUID(uid string) error
}
