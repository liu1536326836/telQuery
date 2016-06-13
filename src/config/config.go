package config

import (
	"fmt"
	"os"

	"clog"
	"orm"
	"pongo"

	"github.com/BurntSushi/toml"
)

var Conf conf

type conf struct {
	Web   WebInfo
	DB    orm.DB
	Log   clog.LogInfo
	Pongo pongo.PongorOption
	Lib   LibInfo
}

type WebInfo struct {
	Listen   string
	Fasthttp bool
}

type LibInfo struct {
	CallAPI bool
}

func ParseConfig(path string) {
	_, err := toml.DecodeFile(path, &Conf)
	if err != nil {
		fmt.Printf("解析配置文件失败, %s\n", err)
		os.Exit(1)
	}
}
