package orm

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type DB struct {
	Driver string
	User   string
	Passwd string
	Host   string
	DBName string
}

type TelInfo struct {
	Mts       string `gorm:"primary_key;not null;type:varchar(7)"`
	Province  string `gorm:"type:varchar(20);not null"`
	CatName   string `gorm:"type:varchar(20);not null"`
	TelString string `gorm:"-"`
}

// 打开数据库
func Open(d DB) error {
	var (
		err error
		dns string
	)

	switch d.Driver {
	case "mysql":
		dns = fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
			d.User, d.Passwd, d.Host, d.DBName)
	default:
		return fmt.Errorf("Unknow database type[%s]", d.Driver)
	}

	db, err = gorm.Open(d.Driver, dns)
	if err != nil {
		return err
	}

	db.SingularTable(true)

	return nil
}

// 创建表
func CreateTable() error {

	if db.HasTable(&TelInfo{}) {
		return fmt.Errorf("要创建的表已经存在")
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB").Set("gorm:table_options",
		"charset=utf8").AutoMigrate(&TelInfo{})

	return nil
}

// 插入数据
func Insert(tel_info TelInfo) error {

	if tel_info.Mts == "" {
		return fmt.Errorf("插入数据失败, Mts字段为空")
	}

	return db.Create(&tel_info).Error

}

// 查询
func Select(table, columns, mts string) (TelInfo, error) {
	var (
		info TelInfo
		err  error
	)

	res, err := db.Table(table).Select(columns).Where("mts = ?", mts).Rows()
	if err != nil {
		return info, err
	}

	defer res.Close()

	if !res.Next() {
		err = fmt.Errorf("查询电话的信息不存在数据库中, Mts[%s]", mts)
	} else {
		res.Scan(&info.Mts, &info.Province, &info.CatName)
	}

	return info, err
}
