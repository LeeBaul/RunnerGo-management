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

	"kp-management/cmd/generate/internal/pkg/dal/model"
)

func newSetting(db *gorm.DB, opts ...gen.DOOption) setting {
	_setting := setting{}

	_setting.settingDo.UseDB(db, opts...)
	_setting.settingDo.UseModel(&model.Setting{})

	tableName := _setting.settingDo.TableName()
	_setting.ALL = field.NewAsterisk(tableName)
	_setting.ID = field.NewInt64(tableName, "id")
	_setting.UserID = field.NewInt64(tableName, "user_id")
	_setting.UserIdentify = field.NewString(tableName, "user_identify")
	_setting.TeamID = field.NewInt64(tableName, "team_id")
	_setting.CreatedAt = field.NewTime(tableName, "created_at")
	_setting.UpdatedAt = field.NewTime(tableName, "updated_at")
	_setting.DeletedAt = field.NewField(tableName, "deleted_at")

	_setting.fillFieldMap()

	return _setting
}

type setting struct {
	settingDo settingDo

	ALL          field.Asterisk
	ID           field.Int64
	UserID       field.Int64 // 用户id
	UserIdentify field.String
	TeamID       field.Int64 // 当前团队id
	CreatedAt    field.Time
	UpdatedAt    field.Time
	DeletedAt    field.Field

	fieldMap map[string]field.Expr
}

func (s setting) Table(newTableName string) *setting {
	s.settingDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s setting) As(alias string) *setting {
	s.settingDo.DO = *(s.settingDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *setting) updateTableName(table string) *setting {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.UserID = field.NewInt64(table, "user_id")
	s.UserIdentify = field.NewString(table, "user_identify")
	s.TeamID = field.NewInt64(table, "team_id")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeletedAt = field.NewField(table, "deleted_at")

	s.fillFieldMap()

	return s
}

func (s *setting) WithContext(ctx context.Context) *settingDo { return s.settingDo.WithContext(ctx) }

func (s setting) TableName() string { return s.settingDo.TableName() }

func (s setting) Alias() string { return s.settingDo.Alias() }

func (s *setting) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *setting) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 7)
	s.fieldMap["id"] = s.ID
	s.fieldMap["user_id"] = s.UserID
	s.fieldMap["user_identify"] = s.UserIdentify
	s.fieldMap["team_id"] = s.TeamID
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["deleted_at"] = s.DeletedAt
}

func (s setting) clone(db *gorm.DB) setting {
	s.settingDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s setting) replaceDB(db *gorm.DB) setting {
	s.settingDo.ReplaceDB(db)
	return s
}

type settingDo struct{ gen.DO }

func (s settingDo) Debug() *settingDo {
	return s.withDO(s.DO.Debug())
}

func (s settingDo) WithContext(ctx context.Context) *settingDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s settingDo) ReadDB() *settingDo {
	return s.Clauses(dbresolver.Read)
}

func (s settingDo) WriteDB() *settingDo {
	return s.Clauses(dbresolver.Write)
}

func (s settingDo) Session(config *gorm.Session) *settingDo {
	return s.withDO(s.DO.Session(config))
}

func (s settingDo) Clauses(conds ...clause.Expression) *settingDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s settingDo) Returning(value interface{}, columns ...string) *settingDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s settingDo) Not(conds ...gen.Condition) *settingDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s settingDo) Or(conds ...gen.Condition) *settingDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s settingDo) Select(conds ...field.Expr) *settingDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s settingDo) Where(conds ...gen.Condition) *settingDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s settingDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *settingDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s settingDo) Order(conds ...field.Expr) *settingDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s settingDo) Distinct(cols ...field.Expr) *settingDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s settingDo) Omit(cols ...field.Expr) *settingDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s settingDo) Join(table schema.Tabler, on ...field.Expr) *settingDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s settingDo) LeftJoin(table schema.Tabler, on ...field.Expr) *settingDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s settingDo) RightJoin(table schema.Tabler, on ...field.Expr) *settingDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s settingDo) Group(cols ...field.Expr) *settingDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s settingDo) Having(conds ...gen.Condition) *settingDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s settingDo) Limit(limit int) *settingDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s settingDo) Offset(offset int) *settingDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s settingDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *settingDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s settingDo) Unscoped() *settingDo {
	return s.withDO(s.DO.Unscoped())
}

func (s settingDo) Create(values ...*model.Setting) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s settingDo) CreateInBatches(values []*model.Setting, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s settingDo) Save(values ...*model.Setting) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s settingDo) First() (*model.Setting, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Setting), nil
	}
}

func (s settingDo) Take() (*model.Setting, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Setting), nil
	}
}

func (s settingDo) Last() (*model.Setting, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Setting), nil
	}
}

func (s settingDo) Find() ([]*model.Setting, error) {
	result, err := s.DO.Find()
	return result.([]*model.Setting), err
}

func (s settingDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Setting, err error) {
	buf := make([]*model.Setting, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s settingDo) FindInBatches(result *[]*model.Setting, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s settingDo) Attrs(attrs ...field.AssignExpr) *settingDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s settingDo) Assign(attrs ...field.AssignExpr) *settingDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s settingDo) Joins(fields ...field.RelationField) *settingDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s settingDo) Preload(fields ...field.RelationField) *settingDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s settingDo) FirstOrInit() (*model.Setting, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Setting), nil
	}
}

func (s settingDo) FirstOrCreate() (*model.Setting, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Setting), nil
	}
}

func (s settingDo) FindByPage(offset int, limit int) (result []*model.Setting, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s settingDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s settingDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s settingDo) Delete(models ...*model.Setting) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *settingDo) withDO(do gen.Dao) *settingDo {
	s.DO = *do.(*gen.DO)
	return s
}
