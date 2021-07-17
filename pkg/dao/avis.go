package dao

import (
	m "adilhaddad.net/agefice-docs/pkg/models"
	"context"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type AvisDao struct {
	Db *gorm.DB
}

// GetAllAvis is a function to get a slice of record(s) from avis table in the recipe database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (ad *AvisDao) GetAllAvis(ctx context.Context, page, pagesize int64, order string) (results []*m.Avis, totalRows int, err error) {

	resultOrm := DB.Model(&m.Avis{})
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

// GetAvis is a function to get a single record from the avis table in the recipe database
// error - ErrNotFound, db Find error
func (ad *AvisDao) GetAvis(ctx context.Context, argID int32) (record *m.Avis, err error) {
	record = &m.Avis{}
	if err = DB.First(record, argID).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// AddAvis is a function to add a single record to avis table in the recipe database
// error - ErrInsertFailed, db save call failed
func (ad *AvisDao) AddAvis(ctx context.Context, record *m.Avis) (result *m.Avis, RowsAffected int64, err error) {
	db := ad.Db.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateAvis is a function to update a single record from avis table in the recipe database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (ad *AvisDao) UpdateAvis(ctx context.Context, argID int32, updated *m.Avis) (result *m.Avis, RowsAffected int64, err error) {

	result = &m.Avis{}
	db := ad.Db.First(result, argID)
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

// DeleteAvis is a function to delete a single record from avis table in the recipe database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (ad *AvisDao) DeleteAvis(ctx context.Context, argID int32) (rowsAffected int64, err error) {

	record := &m.Avis{}
	db := ad.Db.First(record, argID)
	if db.Error != nil {
		return -1, ErrNoRecordFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
