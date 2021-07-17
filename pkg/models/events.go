package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


CREATE TABLE `events` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `descri` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 72,    "name": "FQJFvAjKPxSesvtesIlyMPdmn",    "descri": "TrQOIHmVGiWHncicgZpnUBLwp"}



*/

// Events struct is a row record of the events table in the recipe database
type Events struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	ID int `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	//[ 1] name                                           varchar(50)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Name string `gorm:"column:name;type:varchar;size:50;" json:"name"`
	//[ 2] descri                                         varchar(250)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 250     default: []
	Descri  string    `gorm:"column:descri;type:varchar;size:250;" json:"descri"`
	Recipes []Recipes `gorm:"many2many:recipes_events" json:"-"`
}

// TableName sets the insert table name for this struct type
func (e *Events) TableName() string {
	return "events"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Events) BeforeSave(db *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Events) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Events) Validate(action Action) error {
	return nil
}
