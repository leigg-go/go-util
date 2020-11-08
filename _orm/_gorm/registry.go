package _gorm

import (
	"fmt"
	"github.com/leigg-go/go-util/_util"
	"gorm.io/gorm"
	"sync"
)

/*
MySQL registry with gorm, app can directly call them to initialize MySQL gorm DB client.
*/

var (
	lock      sync.Mutex
	DefGormDB *gorm.DB
)

// init default pool on registry
func MustInitDef(dialer gorm.Dialector, conf *gorm.Config) {
	lock.Lock()
	defer lock.Unlock()
	_util.Must(DefGormDB == nil, fmt.Errorf("_gorm: DefGormDB already exists"))

	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(dialer, conf)
	_util.PanicIfErr(err, nil, fmt.Sprintf("_gorm: gorm.Open err:%v", err))

	DefGormDB = db
}

func MustInit(dialer gorm.Dialector, conf *gorm.Config) *gorm.DB {
	db, err := gorm.Open(dialer, conf)
	_util.PanicIfErr(err, nil, fmt.Sprintf("_gorm: gorm.Open err:%v", err))
	return db
}

func Close() error {
	lock.Lock()
	defer lock.Unlock()
	if DefGormDB != nil {
		db, _ := DefGormDB.DB()
		err := db.Close()
		DefGormDB = nil
		return err
	}
	return nil
}
