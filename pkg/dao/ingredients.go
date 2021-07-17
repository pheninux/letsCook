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

type IngredientsDao struct {
	Db *gorm.DB
}

// GetAllIngredients is a function to get a slice of record(s) from ingredients table in the recipe database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (ingD IngredientsDao) GetAllIngredients(ctx context.Context, page, pagesize int64, order string) (results []*model.Ingredients, totalRows int, err error) {

	resultOrm := ingD.Db.Model(&model.Ingredients{})
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

// GetIngredients is a function to get a single record from the ingredients table in the recipe database
// error - ErrNotFound, db Find error
func (ingD IngredientsDao) GetIngredients(ctx context.Context, argID int32) (record *model.Ingredients, err error) {
	record = &model.Ingredients{}
	if err = ingD.Db.First(record, argID).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// AddIngredients is a function to add a single record to ingredients table in the recipe database
// error - ErrInsertFailed, db save call failed
func (ingD IngredientsDao) AddIngredients(ctx context.Context, record *model.Ingredients) (result *model.Ingredients, RowsAffected int64, err error) {
	db := ingD.Db.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateIngredients is a function to update a single record from ingredients table in the recipe database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (ingD IngredientsDao) UpdateIngredients(ctx context.Context, argID int32, updated *model.Ingredients) (result *model.Ingredients, RowsAffected int64, err error) {

	result = &model.Ingredients{}
	db := ingD.Db.First(result, argID)
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

// DeleteIngredients is a function to delete a single record from ingredients table in the recipe database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (ingD IngredientsDao) DeleteIngredients(ctx context.Context, argID int32) (rowsAffected int64, err error) {

	record := &model.Ingredients{}
	db := ingD.Db.First(record, argID)
	if db.Error != nil {
		return -1, ErrNoRecordFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
