package orm

import (
	"testing"
)

var d = DB{
	Driver: "mysql",
	User:   "root",
	Passwd: "rocky",
	Host:   "127.0.0.1:3306",
	DBName: "tel_info",
}

var info = TelInfo{
	Mts:       "333",
	Province:  "222",
	CatName:   "国家",
	TelString: "444",
}

func TestDBApi(t *testing.T) {
	err := Open(d)
	if err != nil {
		t.Error(err)
	}

	err = CreateTable()
	if err != nil {
		t.Error(err)
	}

	err = Insert(info)
	if err != nil {
		t.Error(err)
	}

	info, err := Select("tel_info", "*", "33")
	if err != nil {
		t.Log("未找到")
	} else {
		t.Log(info)
	}
}
