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


CREATE TABLE `urls` (
  `id` int(11) NOT NULL,
  `url` varchar(200) DEFAULT NULL,
  `id_recipes` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `urls_recipes_FK` (`id_recipes`),
  CONSTRAINT `urls_recipes_FK` FOREIGN KEY (`id_recipes`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 28,    "url": "ojVpvvHllSJaeZhZOVlqTcatr",    "id_recipes": 48}



*/

// Urls struct is a row record of the urls table in the recipe database
type Urls struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] url                                            varchar(200)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	URL string `gorm:"column:url;type:varchar;size:200;" json:"url"`
	//[ 2] id_recipes                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	RecipesID int `gorm:"column:recipes_id;type:int;" json:"recipes_id"`
}

// TableName sets the insert table name for this struct type
func (u *Urls) TableName() string {
	return "urls"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Urls) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Urls) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Urls) Validate(action Action) error {
	return nil
}
