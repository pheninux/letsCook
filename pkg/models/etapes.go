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


CREATE TABLE `etapes` (
  `id` int(11) NOT NULL,
  `detail` varchar(500) DEFAULT NULL,
  `id_recipes` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `etapes_recipes_FK` (`id_recipes`),
  CONSTRAINT `etapes_recipes_FK` FOREIGN KEY (`id_recipes`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 81,    "detail": "kmnLoutbthCkajdLZUiLZCPNU",    "id_recipes": 78}



*/

// Etapes struct is a row record of the etapes table in the recipe database
type Etapes struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] detail                                         varchar(500)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 500     default: []
	Detail string `gorm:"column:detail;type:varchar;size:500;" json:"detail"`
	//[ 2] id_recipes                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	RecipesID int      `gorm:"column:recipes_id;type:int;" json:"recipes_id"`
	Images    []Images `json:"images" gorm:"foreignKey:EtapesID"`
}

// TableName sets the insert table name for this struct type
func (e *Etapes) TableName() string {
	return "etapes"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Etapes) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Etapes) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Etapes) Validate(action Action) error {
	return nil
}
