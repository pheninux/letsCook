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


CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `hashed_password` char(60) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `active` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 51,    "name": "QFnGHXurSWBBSsZiykyPqxGsF",    "email": "CqEpPaAIpTkCGchtguJBrFvKZ",    "hashed_password": "munqPNBfanMyDquwBZUJPpSUZ",    "created": "2306-01-18T03:05:34.178553836+01:00",    "active": 20}



*/

// Users_ struct is a row record of the users table in the recipe database
type Users struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] name                                           varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Name string `gorm:"column:name;type:varchar;size:255;" json:"name"`
	//[ 2] email                                          varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Email string `gorm:"column:email;type:varchar;size:255;" json:"email"`
	//[ 3] hashed_password                                char(60)             null: true   primary: false  isArray: false  auto: false  col: char            len: 60      default: []
	HashedPassword []byte `sql:"-" gorm:"column:hashed_password;type:char;size:60;" json:"hashed_password"`
	//[ 4] created                                        datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Created time.Time `gorm:"column:created;type:datetime;" json:"created"`
	//[ 5] active                                         tinyint              null: true   primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: []
	Active bool `gorm:"column:active;type:bool;" json:"active"`

	AvatarUrl string `json:"avatar_url"`

	Avis []Avis `json:"avis" gorm:"foreignKey:UsersID"`
}

func (u *Users) NewUser(name, email, pass, avatarUrl string) {
	u.Email = email
	u.Name = name
	u.HashedPassword = []byte(pass)
	u.AvatarUrl = avatarUrl

}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "users"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users) BeforeSave(tx *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users) Validate(action Action) error {
	return nil
}
