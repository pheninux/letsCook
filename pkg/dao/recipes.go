package dao

import (
	model "adilhaddad.net/agefice-docs/pkg/models"
	"context"
	"github.com/guregu/null"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	//"github.com/jinzhu/gorm"
	"strings"
	"time"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type RecipesDao struct {
	Db *gorm.DB
}

// GetAllRecipes is a function to get a slice of record(s) from recipes table in the recipe database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (rd *RecipesDao) GetAllRecipes(ctx context.Context, page, pagesize int64, order string) (results []*model.Recipes, totalRows int64, err error) {

	resultOrm := rd.Db.Model(&model.Recipes{})
	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(int(offset)).Limit(int(pagesize))
	} else {
		resultOrm = resultOrm.Limit(int(pagesize))
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Preload("Images").
		Preload("Events").
		Preload("Etapes").Preload("Etapes.Images").
		Preload("Avis").Preload("Urls").Preload("Ingredients").
		Find(&results).Error; err != nil {

		if strings.Contains(err.Error(), ErrScanColumns.Error()) {
			err = ErrScanColumns
		}
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetRecipes is a function to get a single record from the recipes table in the recipe database
// error - ErrNotFound, db Find error
func (rd *RecipesDao) GetRecipe(ctx context.Context, argID int32) (record *model.Recipes, err error) {
	record = &model.Recipes{}
	if err = rd.Db.First(record, argID).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// AddRecipes is a function to add a single record to recipes table in the recipe database
// error - ErrInsertFailed, db save call failed
func (rd *RecipesDao) AddRecipe(ctx context.Context, record *model.Recipes) (result *model.Recipes, RowsAffected int64, err error) {
	db := rd.Db.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateRecipes is a function to update a single record from recipes table in the recipe database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (rd *RecipesDao) UpdateRecipe(ctx context.Context, argID int32, updated *model.Recipes) (result *model.Recipes, RowsAffected int64, err error) {

	result = &model.Recipes{}
	db := rd.Db.First(result, argID)
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

// DeleteRecipes is a function to delete a single record from recipes table in the recipe database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (rd *RecipesDao) DeleteRecipe(ctx context.Context, argID int32) (rowsAffected int64, err error) {

	record := &model.Recipes{}
	db := rd.Db.First(record, argID)
	if db.Error != nil {
		return -1, ErrNoRecordFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
