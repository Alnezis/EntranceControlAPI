package app

import (
	"EntranceControlAPI/api"
	//"github.com/mailgun/mailgun-go/v4"
	//"github.com/mailgun/mailgun-go/v4"
	"EntranceControlAPI/app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

//var MG *mailgun.MailgunImpl
var CFG *config.Config
var DB *sqlx.DB

func init() {
	CFG = config.InitCfg()
	//mg := mailgun.NewMailgun(CFG.Mailgun.Domain, CFG.Mailgun.PrivateAPIKey)
	//	mg.SetAPIBase(mailgun.APIBaseEU)
	//	MG = mg
	conn := `
           host=` + CFG.Db.Host + `
         dbname=` + CFG.Db.DbName + `
		   user=` + CFG.Db.UserName + `
        sslmode=disable
		   port=` + CFG.Db.Port + `
		password=` + CFG.Db.Password + `
`
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	DB = db
	//	initDb()

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id       varchar primary key,
    hash varchar,
    photo_url    varchar,
    fio    varchar
);`)
	api.CheckErrInfo(err, "test users")
}

func initDb() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id       serial not null primary key,
    phone_number    varchar,
 	role_id  INTEGER REFERENCES roles (id),
    balance numeric(18,2) default 0,
    balance_bonus numeric(18,2) default 0,
    confirmed boolean default false
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS roles
(
    id       serial not null primary key,
    name    varchar
);`)
	//role
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS balance_replenishment
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    sum    numeric(18,2),
	created timestamp
);`)

	api.CheckErrInfo(err, "init db 1")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS codes
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    code    varchar,
	created timestamp
);`)
	api.CheckErrInfo(err, "init db 2")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tokens
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    token    varchar,
	created timestamp
);`)
	api.CheckErrInfo(err, "init db 3")

}
