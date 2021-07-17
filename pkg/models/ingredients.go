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


CREATE TABLE `ingredients` (
  `id` int(11) NOT NULL,
  `title` varchar(100) DEFAULT NULL,
  `qt` int(11) DEFAULT NULL,
  `mesure` varchar(2) DEFAULT NULL,
  `descri` varchar(250) DEFAULT NULL,
  `id_recipes` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `ingredients_recipes_FK` (`id_recipes`),
  CONSTRAINT `ingredients_recipes_FK` FOREIGN KEY (`id_recipes`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 56,    "title": "mspwsKxZIeJtaxaojvLsCUtbC",    "qt": 19,    "mesure": "cVoNPZvEePRxufrjcuwDZXDpC",    "descri": "tcZJIpvdMKLwXlxWphBZKuKRc",    "id_recipes": 47}



*/

// Ingredients struct is a row record of the ingredients table in the recipe database
type Ingredients struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] title                                          varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Title string `gorm:"column:title;type:varchar;size:100;" json:"title"`
	//[ 2] qt                                             int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Qt int `gorm:"column:qt;type:int;" json:"qt"`
	//[ 3] mesure                                         varchar(2)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2       default: []
	Mesure string `gorm:"column:mesure;type:varchar;size:2;" json:"mesure"`
	//[ 4] descri                                         varchar(250)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 250     default: []
	Descri string `gorm:"column:descri;type:varchar;size:250;" json:"descri"`
	//[ 5] id_recipes                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	RecipesID int `gorm:"column:recipes_id;type:int;" json:"recipes_id"`
}

// TableName sets the insert table name for this struct type
func (i *Ingredients) TableName() string {
	return "ingredients"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (i *Ingredients) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (i *Ingredients) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (i *Ingredients) Validate(action Action) error {
	return nil
}
