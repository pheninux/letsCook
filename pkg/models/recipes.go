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


CREATE TABLE `recipes` (
  `id` int(11) NOT NULL,
  `title` varchar(100) NOT NULL,
  `descri` varchar(250) DEFAULT NULL,
  `obs` varchar(250) DEFAULT NULL,
  `categorie` varchar(25) DEFAULT NULL,
  `preparation` int(11) DEFAULT NULL,
  `typ` varchar(25) DEFAULT NULL,
  `cuisson` int(11) DEFAULT NULL,
  `repos` int(11) DEFAULT NULL,
  `lvl` varchar(25) DEFAULT NULL,
  `nbr_pers` int(11) DEFAULT NULL,
  `cout` float DEFAULT NULL,
  `share` varchar(25) DEFAULT NULL,
  `valide` int(11) DEFAULT NULL,
  `id_users` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `recipes_users_FK` (`id_users`),
  CONSTRAINT `recipes_users_FK` FOREIGN KEY (`id_users`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "id": 60,    "title": "ObkadoAvSwDMpGCpALoUQKSsJ",    "descri": "hRXIsMPILsSiORgJDKxFyVwMs",    "obs": "LPxuhqcaocInDScqgYWSmkvja",    "categorie": "ERIICNESURycrdnYOHmlukoNW",    "preparation": 19,    "typ": "kIsusbiRMKIUrorYgUIbGqhWt",    "cuisson": 55,    "repos": 44,    "lvl": "MTNwOOmSdRmNtiyNImHVubhvZ",    "nbr_pers": 48,    "cout": 0.87331206,    "share": "BCvZSIYCPHQaifkHlIqVnmmsb",    "valide": 14,    "id_users": 7}



*/

// Recipes struct is a row record of the recipes table in the recipe database
type Recipes struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	ID int `gorm:"primary_key;auto_increment;column:id;type:int;" json:"id"`
	//[ 1] title                                          varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Title string `gorm:"column:title;type:varchar;size:100;" json:"title"`
	//[ 2] descri                                         varchar(250)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 250     default: []
	Descri string `gorm:"column:descri;type:varchar;size:250;" json:"descri"`
	//[ 3] obs                                            varchar(250)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 250     default: []
	Obs string `gorm:"column:obs;type:varchar;size:250;" json:"obs"`
	//[ 4] categorie                                      varchar(25)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 25      default: []
	Categorie string `gorm:"column:categorie;type:varchar;size:25;" json:"categorie"`
	//[ 5] preparation                                    int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Preparation int `gorm:"column:preparation;type:int;" json:"preparation"`
	//[ 6] typ                                            varchar(25)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 25      default: []
	Typ string `gorm:"column:typ;type:varchar;size:25;" json:"typ"`
	//[ 7] cuisson                                        int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Cuisson int `gorm:"column:cuisson;type:int;" json:"cuisson"`
	//[ 8] repos                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Repos int `gorm:"column:repos;type:int;" json:"repos"`
	//[ 9] lvl                                            varchar(25)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 25      default: []
	Lvl string `gorm:"column:lvl;type:varchar;size:25;" json:"lvl"`
	//[10] nbr_pers                                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	NbrPers int `gorm:"column:nbr_pers;type:int;" json:"nbr_pers"`
	//[11] cout                                           float                null: true   primary: false  isArray: false  auto: false  col: float           len: -1      default: []
	Cout float64 `gorm:"column:cout;type:float;" json:"cout"`
	//[12] share                                          varchar(25)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 25      default: []
	Share string `gorm:"column:share;type:varchar;size:25;" json:"share"`
	//[13] valide                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Valide int `gorm:"column:valide;type:int;" json:"valide"`
	//[14] id_users                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UsersID     int           `gorm:"column:users_id;type:int;" json:"users_id"`
	Urls        []Urls        `json:"urls" gorm:"foreignKey:RecipesID"`
	Etapes      []Etapes      `json:"etapes"  gorm:"foreignKey:RecipesID"`
	Ingredients []Ingredients `gorm:"foreignKey:RecipesID" json:"ingredients"`
	Events      []Events      `json:"events" gorm:"many2many:recipes_events"`
	Images      []Images      `json:"images" gorm:"foreignKey:RecipesID"`
	Avis        []Avis        `json:"avis" gorm:"foreignKey:RecipesID"`
	Groupes     []Groupes     `json:"groupes" gorm:"many2many:recipes_groupes"`
}

// TableName sets the insert table name for this struct type
func (r *Recipes) TableName() string {
	return "recipes"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Recipes) BeforeSave(db *gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Recipes) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Recipes) Validate(action Action) error {
	return nil
}
