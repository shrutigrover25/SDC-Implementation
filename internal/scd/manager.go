package scd

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

// SCDManager implements SCDRepository with optimized SCD operations
type SCDManager[T SCDModel[T]] struct {
	db *gorm.DB
}

// queryBuilder implements SCDQueryBuilder for fluent query construction
type queryBuilder[T SCDModel[T]] struct {
	db         *gorm.DB
	manager    *SCDManager[T]
	conditions []condition
	orders     []string
	limitVal   *int
	offsetVal  *int
	latestOnly bool
	withHist   bool
	asOfDate   *time.Time
	startDate  *time.Time
	endDate    *time.Time
}

type condition struct {
	query string
	args  []interface{}
}

func NewManager[T SCDModel[T]](db *gorm.DB) SCDRepository[T] {
	return &SCDManager[T]{db: db}
}

func NewQueryBuilder[T SCDModel[T]](db *gorm.DB, manager *SCDManager[T]) SCDQueryBuilder[T] {
	return &queryBuilder[T]{
		db:      db,
		manager: manager,
	}
}

// SCDRepository implementation

func (m *SCDManager[T]) Create(entity T) (T, error) {
	now := time.Now()
	entity = entity.SetCreatedAt(now)
	entity = entity.SetUpdatedAt(now)
	err := m.db.Create(&entity).Error
	return entity, err
}

func (m *SCDManager[T]) FindByUID(uid string) (T, error) {
	var entity T
	err := m.db.Where("uid = ?", uid).First(&entity).Error
	return entity, err
}

func (m *SCDManager[T]) Update(uid string, updateFn func(T) T) (T, error) {
	current, err := m.FindByUID(uid)
	if err != nil {
		return current, err
	}

	newVersion := current.CopyForNewVersion()
	newVersion = updateFn(newVersion)

	now := time.Now()
	newVersion = newVersion.SetCreatedAt(now)
	newVersion = newVersion.SetUpdatedAt(now)

	err = m.db.Create(&newVersion).Error
	return newVersion, err
}

func (m *SCDManager[T]) Delete(uid string) error {
	return fmt.Errorf("delete operation should be implemented based on business logic")
}

func (m *SCDManager[T]) Query() SCDQueryBuilder[T] {
	return NewQueryBuilder(m.db, m)
}

func (m *SCDManager[T]) CreateBatch(entities []T) error {
	now := time.Now()
	for i := range entities {
		entities[i] = entities[i].SetCreatedAt(now)
		entities[i] = entities[i].SetUpdatedAt(now)
	}
	return m.db.CreateInBatches(entities, 100).Error
}

func (m *SCDManager[T]) UpdateBatch(updates map[string]func(T) T) error {
	for uid, updateFn := range updates {
		_, err := m.Update(uid, updateFn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SCDManager[T]) GetLatestVersion(id string) (T, error) {
	var entity T
	tableName := entity.TableName()

	subQuery := m.db.Table(tableName).
		Select("*, ROW_NUMBER() OVER (PARTITION BY id ORDER BY version DESC) as rn").
		Where("id = ?", id)

	err := m.db.Table("(?) as ranked", subQuery).
		Where("rn = 1").
		Select("*").
		First(&entity).Error

	return entity, err
}

func (m *SCDManager[T]) GetVersionHistory(id string) ([]T, error) {
	var entities []T
	err := m.db.Where("id = ?", id).Order("version ASC").Find(&entities).Error
	return entities, err
}

func (m *SCDManager[T]) GetVersionAt(id string, date time.Time) (T, error) {
	var entity T
	tableName := entity.TableName()

	subQuery := m.db.Table(tableName).
		Select("*, ROW_NUMBER() OVER (PARTITION BY id ORDER BY created_at DESC) as rn").
		Where("id = ? AND created_at <= ?", id, date)

	err := m.db.Table("(?) as ranked", subQuery).
		Where("rn = 1").
		Select("*").
		First(&entity).Error

	return entity, err
}

// SCDQueryBuilder implementation

func (qb *queryBuilder[T]) Latest() SCDQueryBuilder[T] {
	qb.latestOnly = true
	return qb
}

func (qb *queryBuilder[T]) Where(query interface{}, args ...interface{}) SCDQueryBuilder[T] {
	qb.conditions = append(qb.conditions, condition{
		query: fmt.Sprintf("%v", query),
		args:  args,
	})
	return qb
}

func (qb *queryBuilder[T]) WhereIn(column string, values interface{}) SCDQueryBuilder[T] {
	qb.conditions = append(qb.conditions, condition{
		query: fmt.Sprintf("%s IN ?", column),
		args:  []interface{}{values},
	})
	return qb
}

func (qb *queryBuilder[T]) Order(value interface{}) SCDQueryBuilder[T] {
	qb.orders = append(qb.orders, fmt.Sprintf("%v", value))
	return qb
}

func (qb *queryBuilder[T]) Limit(limit int) SCDQueryBuilder[T] {
	qb.limitVal = &limit
	return qb
}

func (qb *queryBuilder[T]) Offset(offset int) SCDQueryBuilder[T] {
	qb.offsetVal = &offset
	return qb
}

func (qb *queryBuilder[T]) JoinLatest(table string, condition string) SCDQueryBuilder[T] {
	return qb
}

func (qb *queryBuilder[T]) AsOfDate(date time.Time) SCDQueryBuilder[T] {
	qb.asOfDate = &date
	return qb
}

func (qb *queryBuilder[T]) BetweenDates(start, end time.Time) SCDQueryBuilder[T] {
	qb.startDate = &start
	qb.endDate = &end
	return qb
}

func (qb *queryBuilder[T]) WithHistory() SCDQueryBuilder[T] {
	qb.withHist = true
	return qb
}

func (qb *queryBuilder[T]) GroupByID() SCDQueryBuilder[T] {
	return qb
}

func (qb *queryBuilder[T]) Raw() *gorm.DB {
	return qb.buildQuery()
}

func (qb *queryBuilder[T]) Find() ([]T, error) {
	var results []T
	err := qb.buildQuery().Find(&results).Error
	return results, err
}

func (qb *queryBuilder[T]) First() (T, error) {
	var result T
	err := qb.buildQuery().First(&result).Error
	return result, err
}

func (qb *queryBuilder[T]) Count() (int64, error) {
	var count int64
	err := qb.buildQuery().Count(&count).Error
	return count, err
}

func (qb *queryBuilder[T]) buildQuery() *gorm.DB {
	var dummy T
	tableName := dummy.TableName()

	var query *gorm.DB

	if qb.latestOnly && !qb.withHist {
		subQuery := qb.db.Table(tableName).
			Select("*, ROW_NUMBER() OVER (PARTITION BY id ORDER BY version DESC) as rn")

		query = qb.db.Table("(?) as latest", subQuery).
			Where("rn = 1").
			Select("id, version, uid, status, rate, title, company_id, contractor_id, created_at, updated_at")
	} else if qb.asOfDate != nil {
		subQuery := qb.db.Table(tableName).
			Select("*, ROW_NUMBER() OVER (PARTITION BY id ORDER BY created_at DESC) as rn").
			Where("created_at <= ?", *qb.asOfDate)

		query = qb.db.Table("(?) as pit", subQuery).
			Where("rn = 1").
			Select("id, version, uid, status, rate, title, company_id, contractor_id, created_at, updated_at")
	} else {
		query = qb.db.Table(tableName)
	}

	for _, cond := range qb.conditions {
		query = query.Where(cond.query, cond.args...)
	}

	if qb.startDate != nil && qb.endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *qb.startDate, *qb.endDate)
	}

	for _, order := range qb.orders {
		query = query.Order(order)
	}

	if qb.limitVal != nil {
		query = query.Limit(*qb.limitVal)
	}
	if qb.offsetVal != nil {
		query = query.Offset(*qb.offsetVal)
	}

	return query
}

// Legacy compatibility

func (m *SCDManager[T]) GetLatest() *gorm.DB {
	return m.Query().Latest().Raw()
}

func (m *SCDManager[T]) FindByUID_Legacy(uid string) (T, error) {
	return m.FindByUID(uid)
}

func (m *SCDManager[T]) Insert(newItem T) error {
	_, err := m.Create(newItem)
	return err
}

func (m *SCDManager[T]) CreateNewVersion(old T) (T, error) {
	newItem := old.CopyForNewVersion()
	now := time.Now()
	newItem = newItem.SetCreatedAt(now)
	newItem = newItem.SetUpdatedAt(now)
	err := m.db.Create(&newItem).Error
	return newItem, err
}
