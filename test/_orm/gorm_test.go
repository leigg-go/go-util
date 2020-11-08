package _orm

import (
	"github.com/leigg-go/go-util/_orm/_gorm"
	"github.com/leigg-go/go-util/_util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestMustInitDef(t *testing.T) {
	dsn := "root:123@tcp(127.0.0.1:3306)/main?charset=utf8mb4&parseTime=True&loc=Local"
	conf := &gorm.Config{
		SkipDefaultTransaction: true,
		//Logger:                 logger.Default,
		//ConnPool:               connPool,
	}
	_gorm.MustInitDef(mysql.Open(dsn), conf)
	db, err := _gorm.DefGormDB.DB()
	_util.PanicIfErr(err, nil)

	// Setting pool params
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)
}
