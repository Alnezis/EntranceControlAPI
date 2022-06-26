package user

import (
	"EntranceControlAPI/api"
	"EntranceControlAPI/app"
	"crypto/md5"
	"encoding/hex"
	"time"
)

func IsConfirmed(phoneNumber int) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from users where confirmed=true and phone_number=$1)", phoneNumber)
	api.CheckErrInfo(err, "IsConfirmed")
	return i
}

func Exists(id string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select * from users where id=$1)", id)
	api.CheckErrInfo(err, "Exists")
	return i
}

//func IsTokenValid(serverKey, nick, password string) bool {
//	var i bool
//	err := app.DB.Get(&i, "select exists(select id from accounts where serverkey=$1 and password=$2)", serverKey, password)
//	api.CheckErrInfo(err, "IsTokenValid")
//	return i
//}

func Confirmed(id int) {
	_, err := app.DB.Exec("update users set confirmed=true where id=$1", id)
	api.CheckErrInfo(err, "Confirmed")
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

//func (i User) CreateToken() string {
//	var t, err = token.GenerateJWT(i.ID, i.PhoneNumber)
//	api.CheckErrInfo(err, "GenerateJWT")
//	_, err = app.DB.Exec(`INSERT INTO tokens (user_id, token, created) VALUES ($1, $2, $3)`, i.ID, t, time.Now())
//	api.CheckErrInfo(err, "CreateToken")
//	return t
//}

func IsToken(id int, token string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from tokens where id=$1 and token=$2)", id, token)
	api.CheckErrInfo(err, "IsToken")
	return i
}

func IsStaticHash(id, hash string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from users where id=$1 and hash=$2)", id, hash)
	api.CheckErrInfo(err, "IsPhotoStatic")
	return i
}
func UpdateHash(id, url string) {
	_, err := app.DB.Exec("update users set hash=$1 where id=$2", url, id)
	api.CheckErrInfo(err, "UpdateHash")
}

func IsStaticFIO(id, fio string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from users where id=$1 and fio=$2)", id, fio)
	api.CheckErrInfo(err, "IsPhotoStatic")
	return i
}
func UpdateFIO(id, fio string) {
	_, err := app.DB.Exec("update users set fio=$1 where id=$2", fio, id)
	api.CheckErrInfo(err, "UpdateHash")
}

func UpdatePhotoURL(id, url string) {
	_, err := app.DB.Exec("update users set photo_url=$1 where id=$2", url, id)
	api.CheckErrInfo(err, "UpdatePhotoURL")
}

type User struct {
	ID   string `json:"id,omitempty"`
	Hash string `json:"hash,omitempty"`

	PhotoURL string `json:"photo_url" db:"photo_url"`
	FIO      string `json:"fio" db:"fio"`
}

func GetUser(id string) User {
	var i User
	err := app.DB.Get(&i, `select * from users where id=$1`, id)
	api.CheckErrInfo(err, "GetUser")
	return i
}

func NewUser(id, photoURL, fio, hash string) string {
	err := app.DB.Get(&id, `INSERT INTO users (id, hash, photo_url, fio) VALUES ($1,$2,$3,$4) returning id`, id, hash, photoURL, fio)
	api.CheckErrInfo(err, "NewUser")
	return id
}

func NewCode(id, code int) bool {
	_, err := app.DB.Exec(`INSERT INTO codes (user_id, code, created) VALUES ($1, $2, $3)`,
		id, code, time.Now())
	api.CheckErrInfo(err, "NewCode")
	return true
}

func CheckCode(id int, code string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from codes where user_id=$1 and code=$2)", id, code)
	api.CheckErrInfo(err, "CheckCode")
	return i
}
