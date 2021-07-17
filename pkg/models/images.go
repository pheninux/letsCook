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


CREATE TABLE `images` (
  `id` int(11) NOT NULL,
  `src` varchar(700) DEFAULT NULL,
  `typ` varchar(10) DEFAULT NULL,
  `size` int(11) DEFAULT NULL,
  `id_etapes` int(11) NOT NULL,
  `id_recipes` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `images_etapes_FK` (`id_etapes`),
  KEY `images_recipes0_FK` (`id_recipes`),
  CONSTRAINT `images_etapes_FK` FOREIGN KEY (`id_etapes`) REFERENCES `etapes` (`id`),
  CONSTRAINT `images_recipes0_FK` FOREIGN KEY (`id_recipes`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 33,    "src": "rcCKfveSLpCThEVeOokdXCZXJ",    "typ": "CQlltusxecgUJHvVulEGNjaUZ",    "size": 26,    "id_etapes": 76,    "id_recipes": 4}



*/

// Images struct is a row record of the images table in the recipe database
type Images struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] src                                            varchar(700)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 700     default: []
	Src string `gorm:"column:src;type:varchar;size:700;" json:"src"`
	//[ 2] typ                                            varchar(10)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 10      default: []
	Typ string `gorm:"column:typ;type:varchar;size:10;" json:"typ"`
	//[ 3] size                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Size int `gorm:"column:size;type:int;" json:"size"`
	//[ 4] id_etapes                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	EtapesID int `gorm:"column:etapes_id;type:int;" json:"etapes_id"`
	//[ 5] id_recipes                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	RecipesID int `gorm:"column:recipes_id;type:int;" json:"recipes_id"`
}

// TableName sets the insert table name for this struct type
func (i *Images) TableName() string {
	return "images"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (i *Images) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (i *Images) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (i *Images) Validate(action Action) error {
	return nil
}
