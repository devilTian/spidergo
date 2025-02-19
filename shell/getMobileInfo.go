package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type phone_data struct {
	Number string
	Area_code string
	Province string
	Carrier string
	City string
}
type resp_data struct {
	Code int
	Msg string
	Phone phone_data
	Result []int
}

var URL string = "https://haoma.sginput.qq.com/xcx/search?query=%s"
var mobile_chan chan phone_data

func main() {
	fmt.Println("Start!!!")
	var wg sync.WaitGroup
	mobile_chan = make(chan phone_data, 1000)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go handle(i, &wg)
	}
	insertToDb()
	wg.Wait()
	fmt.Println("Finish!!!")
}

func insertToDb() {
	v := <- mobile_chan
	fmt.Println(v)
}

func handle(i int, wg *sync.WaitGroup) {
	var res_data resp_data
	defer wg.Done()
	mobile := fmt.Sprintf("1358199959%d", i)
	// 调用接口获取归属地数据
	resp, err := http.Get(fmt.Sprintf(URL, mobile))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 读取响应数据
	if ret, err := io.ReadAll(resp.Body); err != nil {
		fmt.Printf("fail to read response body, err: %s", err)
		return
	} else if err := json.Unmarshal(ret, &res_data); err != nil {
		// 将响应数据解析成json
		fmt.Printf("fail to parse json, err: %s", err)
		return
	}
	mobile_chan <- res_data.Phone
}