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


CREATE TABLE `avis` (
  `id` int(11) NOT NULL,
  `rate` float DEFAULT NULL,
  `comment` varchar(200) DEFAULT NULL,
  `id_recipes` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `avis_recipes_FK` (`id_recipes`),
  CONSTRAINT `avis_recipes_FK` FOREIGN KEY (`id_recipes`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 32,    "rate": 0.883654,    "comment": "EDHlnVxfMQaymenFcgdntYLMV",    "id_recipes": 87}



*/

// Avis struct is a row record of the avis table in the recipe database
type Avis struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] rate                                           float                null: true   primary: false  isArray: false  auto: false  col: float           len: -1      default: []
	Rate int `gorm:"column:rate;type:float;" json:"rate"`
	//[ 2] comment                                        varchar(200)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	Comment string `gorm:"column:comment;type:varchar;size:200;" json:"comment"`
	//[ 3] id_recipes                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	RecipesID int `gorm:"column:recipes_id;type:int;" json:"recipes_id"`

	UsersID int `json:"users_id"`

	Date time.Time `json:"date"`
}

// TableName sets the insert table name for this struct type
func (a *Avis) TableName() string {
	return "avis"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *Avis) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *Avis) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *Avis) Validate(action Action) error {
	return nil
}
