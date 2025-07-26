package codes

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getGormConn() *gorm.DB {
	dsn := "host=localhost user=postgres password=lyan123 dbname=gorm_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}

func IsError(err error) {
	if err != nil {
		panic(err)
	}
}

func getSqlxConn() *sqlx.DB {
	dsn := "host=localhost user=postgres password=lyan123 dbname=gorm_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}
