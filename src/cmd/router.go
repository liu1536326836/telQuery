package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"clog"
	"config"
	"libs"
	"orm"

	"github.com/labstack/echo"
)

var ErrNotVaildNumber = errors.New("Not a vaild telephone number")

const (
	URL       = "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel="
	SHOW_PATH = "show.html"
)

func router(e *echo.Echo) {
	e.GET("/tel", handleTelInfo)
}

func handleTelInfo(c echo.Context) error {
	var (
		info orm.TelInfo
		err  error
	)

	// 获取查询的电话号码
	tel := c.QueryParam("tel")
	if libs.IsValidTelNumber(tel) {
		clog.Info("查询的电话号码为: %s", tel)
	} else {
		clog.Error("查询的电话号码格式错误: %s", tel)
		return ErrNotVaildNumber
	}

	// 先在数据库中查找，找不到在淘宝中查找
	info, err = orm.Select("tel_info", "*", tel[:7])
	if err != nil {
		clog.Error("在数据库中查找数据失败, %s", err)

		if config.Conf.Lib.CallAPI {
			ctx, err := libs.GetHTMLContent(URL + tel)
			if err != nil {
				clog.Error("获取电话信息失败, TEL[%s], %s", tel, err)
				return err
			}

			err = json.Unmarshal([]byte(ctx), &info)
			if err != nil {
				clog.Error("解析成JSON格式失败, %s", err)
				return err
			}

			if err := orm.Insert(info); err != nil {
				clog.Error("往数据库中插入数据失败, %s", err)
			}
		}
	}

	info.TelString = tel
	clog.Info("电话号码的详细信息: %v", info)

	err = c.Render(http.StatusOK, SHOW_PATH, map[string]interface{}{"info": info})
	if err != nil {
		clog.Error("解析模板失败, %s", err)
		return err
	}

	return nil
}
