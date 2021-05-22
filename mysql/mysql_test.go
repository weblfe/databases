package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

type User struct {
	ID   int64  `json:"id" gorose:"id"`
	Name string `json:"user_nicename" gorose:"user_nicename"`
}

func (user User) TableName() string {
	return "user"
}

func (user User) New() interface{} {
	return &User{}
}

func (user User) Collection() interface{} {
	var arr []*User
	return arr
}

type UserModel struct {
	Model
}

func NewUserModel() *UserModel {
	var m = &UserModel{}
	return m
}

func getConfig() (*DbManager, error) {
	var (
		file, _ = filepath.Abs("./config/database.yml")
		v       = viper.New()
	)
	// fmt.Println(file)
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	// fmt.Println(v.AllKeys())
	return NewDbMangerByViper(*v)
}

func TestDbManager_Get(t *testing.T) {
	var db, err = getConfig()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v", db.Get("mysql").GetOrm().Table(""))
}

func TestNewUser(t *testing.T) {
	var db, err = getConfig()
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(db.Get(db.Default).GetOrm())
	var user = NewUserModel()
	user.Init(db)
	user.Bind(&User{})
	data := user.GetMapper().New()
	errs := user.Table(data).Where("id", "=", 1).Select()
	if errs != nil {
		t.Error(errs)
	}
	var user2 []User
	err = user.Query(&user2).Where("id","<>", 43302).Limit(2).Select()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user2)

}
