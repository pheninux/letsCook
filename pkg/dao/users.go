package dao

import (
	m "adilhaddad.net/agefice-docs/pkg/models"
	"context"
	"gorm.io/gorm"
	//"github.com/jinzhu/gorm"
	"github.com/guregu/null"
	"github.com/satori/go.uuid"
	"time"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type UsersDao struct {
	Db *gorm.DB
}

// GetAllUsers_ is a function to get a slice of record(s) from users table in the recipe database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func (ud *UsersDao) GetAllUsers(ctx context.Context, page, pagesize int64, order string) (results []*m.Users, totalRows int64, err error) {

	resultOrm := ud.Db.Model(&m.Users{})
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

	if err = resultOrm.Find(&results).Error; err != nil {
		err = ErrNoRecordFound
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetUsers_ is a function to get a single record from the users table in the recipe database
// error - ErrNotFound, db Find error
func (ud *UsersDao) GetUser(ctx context.Context, argID int32) (record *m.Users, err error) {
	record = &m.Users{}
	if err = ud.Db.Omit("hashed_password").First(&record, argID).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// AddUsers_ is a function to add a single record to users table in the recipe database
// error - ErrInsertFailed, db save call failed
func (ud *UsersDao) AddUser(ctx context.Context, record *m.Users) (result *m.Users, RowsAffected int64, err error) {
	db := ud.Db.Create(&record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateUsers_ is a function to update a single record from users table in the recipe database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func (ud *UsersDao) UpdateUser(ctx context.Context, argID int32, updated *m.Users) (result *m.Users, RowsAffected int64, err error) {

	result = &m.Users{}
	db := ud.Db.First(&result, argID)
	if err = db.Error; err != nil {
		return nil, -1, ErrNoRecordFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = db.Save(&result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// DeleteUsers_ is a function to delete a single record from users table in the recipe database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func (ud *UsersDao) DeleteUser(ctx context.Context, argID int32) (rowsAffected int64, err error) {

	record := &m.Users{}
	db := ud.Db.First(record, argID)
	if db.Error != nil {
		return -1, ErrNoRecordFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}

// GetUserByEmail is a function to get a single record from the users table in the recipe database
// error - ErrNotFound, db Find error
func (ud *UsersDao) GetUserByEmail(ctx context.Context, argEmail string) (record *m.Users, err error) {
	record = &m.Users{}
	if err = ud.Db.First(&record, "email = ?", argEmail).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}

// GetUserByEmail is a function to get a single record from the users table in the recipe database
// error - ErrNotFound, db Find error
func (ud *UsersDao) GetUserByEmailAndMdp(ctx context.Context, argEmail, mdp string) (record *m.Users, err error) {
	record = &m.Users{}
	if err = ud.Db.First(&record, "email = ?", argEmail).Error; err != nil {
		err = ErrNoRecordFound
		return record, err
	}

	return record, nil
}
