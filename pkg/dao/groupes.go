package dao

import (
	model "adilhaddad.net/agefice-docs/pkg/models"
	"context"
	"github.com/guregu/null"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type GroupesDao struct {
	Db *gorm.DB
}

// GetAllGroupes is a function to get a slice of record(s) from groupes table in the recipe database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (gd *GroupesDao) GetAllGroupes(ctx context.Context, page, pagesize int64, order string) (results []*model.Groupes, totalRows int, err error) {

	resultOrm := gd.Db.Model(&model.Groupes{})
	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		err = ErrNoRecordFound
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetGroupes is a function to get a single record from the groupes table in the recipe database
// error - ErrNotFound, db Find error
func (gd *GroupesDao) GetGroupes(ctx context.Context, argID int32) (record *model.Groupes, err error) {
	record = &model.Groupes{}
	if err = gd.Db.First(record, argID).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// AddGroupes is a function to add a single record to groupes table in the recipe database
// error - ErrInsertFailed, db save call failed
func (gd *GroupesDao) AddGroupes(ctx context.Context, record *model.Groupes) (result *model.Groupes, RowsAffected int64, err error) {
	db := gd.Db.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateGroupes is a function to update a single record from groupes table in the recipe database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (gd *GroupesDao) UpdateGroupes(ctx context.Context, argID int32, updated *model.Groupes) (result *model.Groupes, RowsAffected int64, err error) {

	result = &model.Groupes{}
	db := gd.Db.First(result, argID)
	if err = db.Error; err != nil {
		return nil, -1, ErrNoRecordFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// DeleteGroupes is a function to delete a single record from groupes table in the recipe database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (gd *GroupesDao) DeleteGroupes(ctx context.Context, argID int32) (rowsAffected int64, err error) {

	record := &model.Groupes{}
	db := gd.Db.First(record, argID)
	if db.Error != nil {
		return -1, ErrNoRecordFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
