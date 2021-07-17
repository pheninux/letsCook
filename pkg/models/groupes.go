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


CREATE TABLE `groupes` (
  `id` int(11) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `id_users` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `groupes_users_FK` (`id_users`),
  CONSTRAINT `groupes_users_FK` FOREIGN KEY (`id_users`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 36,    "name": "YNKXWqEcAOUwhlydamRQgyYnT",    "created": "2081-03-06T08:52:53.070189888+01:00",    "id_users": 12}



*/

// Groupes struct is a row record of the groupes table in the recipe database
type Groupes struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] name                                           varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name string `gorm:"column:name;type:varchar;size:100;" json:"name"`
	//[ 2] created                                        datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Created time.Time `gorm:"column:created;type:datetime;" json:"created"`
	//[ 3] id_users                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UsersID int       `gorm:"column:users_id;type:int;" json:"users_id"`
	Recipes []Recipes `json:"-" gorm:"many2many:recipes_groupes"`
}

// TableName sets the insert table name for this struct type
func (g *Groupes) TableName() string {
	return "groupes"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (g *Groupes) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (g *Groupes) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (g *Groupes) Validate(action Action) error {
	return nil
}
