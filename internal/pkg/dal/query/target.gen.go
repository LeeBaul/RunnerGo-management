// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"kp-management/internal/pkg/dal/model"
)

func newTarget(db *gorm.DB) target {
	_target := target{}

	_target.targetDo.UseDB(db)
	_target.targetDo.UseModel(&model.Target{})

	tableName := _target.targetDo.TableName()
	_target.ALL = field.NewField(tableName, "*")
	_target.ID = field.NewInt64(tableName, "id")
	_target.TeamID = field.NewInt64(tableName, "team_id")
	_target.TargetType = field.NewString(tableName, "target_type")
	_target.Name = field.NewString(tableName, "name")
	_target.ParentID = field.NewInt64(tableName, "parent_id")
	_target.Method = field.NewString(tableName, "method")
	_target.Sort = field.NewInt32(tableName, "sort")
	_target.TypeSort = field.NewInt32(tableName, "type_sort")
	_target.Status = field.NewInt32(tableName, "status")
	_target.Version = field.NewInt32(tableName, "version")
	_target.CreatedUserID = field.NewInt64(tableName, "created_user_id")
	_target.RecentUserID = field.NewInt64(tableName, "recent_user_id")
	_target.Source = field.NewInt32(tableName, "source")
	_target.CreatedAt = field.NewTime(tableName, "created_at")
	_target.UpdatedAt = field.NewTime(tableName, "updated_at")
	_target.DeletedAt = field.NewField(tableName, "deleted_at")

	_target.fillFieldMap()

	return _target
}

type target struct {
	targetDo targetDo

	ALL           field.Field
	ID            field.Int64
	TeamID        field.Int64
	TargetType    field.String
	Name          field.String
	ParentID      field.Int64
	Method        field.String
	Sort          field.Int32
	TypeSort      field.Int32
	Status        field.Int32
	Version       field.Int32
	CreatedUserID field.Int64
	RecentUserID  field.Int64
	Source        field.Int32
	CreatedAt     field.Time
	UpdatedAt     field.Time
	DeletedAt     field.Field

	fieldMap map[string]field.Expr
}

func (t target) Table(newTableName string) *target {
	t.targetDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t target) As(alias string) *target {
	t.targetDo.DO = *(t.targetDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *target) updateTableName(table string) *target {
	t.ALL = field.NewField(table, "*")
	t.ID = field.NewInt64(table, "id")
	t.TeamID = field.NewInt64(table, "team_id")
	t.TargetType = field.NewString(table, "target_type")
	t.Name = field.NewString(table, "name")
	t.ParentID = field.NewInt64(table, "parent_id")
	t.Method = field.NewString(table, "method")
	t.Sort = field.NewInt32(table, "sort")
	t.TypeSort = field.NewInt32(table, "type_sort")
	t.Status = field.NewInt32(table, "status")
	t.Version = field.NewInt32(table, "version")
	t.CreatedUserID = field.NewInt64(table, "created_user_id")
	t.RecentUserID = field.NewInt64(table, "recent_user_id")
	t.Source = field.NewInt32(table, "source")
	t.CreatedAt = field.NewTime(table, "created_at")
	t.UpdatedAt = field.NewTime(table, "updated_at")
	t.DeletedAt = field.NewField(table, "deleted_at")

	t.fillFieldMap()

	return t
}

func (t *target) WithContext(ctx context.Context) *targetDo { return t.targetDo.WithContext(ctx) }

func (t target) TableName() string { return t.targetDo.TableName() }

func (t target) Alias() string { return t.targetDo.Alias() }

func (t *target) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *target) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 16)
	t.fieldMap["id"] = t.ID
	t.fieldMap["team_id"] = t.TeamID
	t.fieldMap["target_type"] = t.TargetType
	t.fieldMap["name"] = t.Name
	t.fieldMap["parent_id"] = t.ParentID
	t.fieldMap["method"] = t.Method
	t.fieldMap["sort"] = t.Sort
	t.fieldMap["type_sort"] = t.TypeSort
	t.fieldMap["status"] = t.Status
	t.fieldMap["version"] = t.Version
	t.fieldMap["created_user_id"] = t.CreatedUserID
	t.fieldMap["recent_user_id"] = t.RecentUserID
	t.fieldMap["source"] = t.Source
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
	t.fieldMap["deleted_at"] = t.DeletedAt
}

func (t target) clone(db *gorm.DB) target {
	t.targetDo.ReplaceDB(db)
	return t
}

type targetDo struct{ gen.DO }

func (t targetDo) Debug() *targetDo {
	return t.withDO(t.DO.Debug())
}

func (t targetDo) WithContext(ctx context.Context) *targetDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t targetDo) ReadDB() *targetDo {
	return t.Clauses(dbresolver.Read)
}

func (t targetDo) WriteDB() *targetDo {
	return t.Clauses(dbresolver.Write)
}

func (t targetDo) Clauses(conds ...clause.Expression) *targetDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t targetDo) Returning(value interface{}, columns ...string) *targetDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t targetDo) Not(conds ...gen.Condition) *targetDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t targetDo) Or(conds ...gen.Condition) *targetDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t targetDo) Select(conds ...field.Expr) *targetDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t targetDo) Where(conds ...gen.Condition) *targetDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t targetDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *targetDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t targetDo) Order(conds ...field.Expr) *targetDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t targetDo) Distinct(cols ...field.Expr) *targetDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t targetDo) Omit(cols ...field.Expr) *targetDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t targetDo) Join(table schema.Tabler, on ...field.Expr) *targetDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t targetDo) LeftJoin(table schema.Tabler, on ...field.Expr) *targetDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t targetDo) RightJoin(table schema.Tabler, on ...field.Expr) *targetDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t targetDo) Group(cols ...field.Expr) *targetDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t targetDo) Having(conds ...gen.Condition) *targetDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t targetDo) Limit(limit int) *targetDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t targetDo) Offset(offset int) *targetDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t targetDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *targetDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t targetDo) Unscoped() *targetDo {
	return t.withDO(t.DO.Unscoped())
}

func (t targetDo) Create(values ...*model.Target) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t targetDo) CreateInBatches(values []*model.Target, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t targetDo) Save(values ...*model.Target) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t targetDo) First() (*model.Target, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Target), nil
	}
}

func (t targetDo) Take() (*model.Target, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Target), nil
	}
}

func (t targetDo) Last() (*model.Target, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Target), nil
	}
}

func (t targetDo) Find() ([]*model.Target, error) {
	result, err := t.DO.Find()
	return result.([]*model.Target), err
}

func (t targetDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Target, err error) {
	buf := make([]*model.Target, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t targetDo) FindInBatches(result *[]*model.Target, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t targetDo) Attrs(attrs ...field.AssignExpr) *targetDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t targetDo) Assign(attrs ...field.AssignExpr) *targetDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t targetDo) Joins(fields ...field.RelationField) *targetDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t targetDo) Preload(fields ...field.RelationField) *targetDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t targetDo) FirstOrInit() (*model.Target, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Target), nil
	}
}

func (t targetDo) FirstOrCreate() (*model.Target, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Target), nil
	}
}

func (t targetDo) FindByPage(offset int, limit int) (result []*model.Target, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t targetDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t targetDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t targetDo) Delete(models ...*model.Target) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *targetDo) withDO(do gen.Dao) *targetDo {
	t.DO = *do.(*gen.DO)
	return t
}
