package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	//	"time"
)

const URL = "http://localhost:8080/tel?tel="

//const URL = "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel="

//
func main() {
	ch := make(chan bool)
	go get(ch, 14000000000)
	go get(ch, 15000000000)
	go get(ch, 18000000000)
	<-ch
	<-ch
	<-ch
}

// 1[3|4|5|8][0-9]\d{4,8}
func get(ch chan<- bool, i int) {
	j := i

	for {
		if j > i+999999999 {
			ch <- true
			return
		}

		fmt.Printf("获取电话号码为%d的信息\n", j)

		resp, err := http.Get(URL + strconv.Itoa(j))
		if err != nil {
			fmt.Printf("Get, %s", err)
			return
		}

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//		time.Sleep(200 * time.Millisecond)
		resp.Body.Close()

		j = j + 10000
	}
}
