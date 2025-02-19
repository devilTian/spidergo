/*
 * @Author: tianye@shimiotech.cn
 * @Date: 2023-12-04 11:07:56
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Bw_record struct {
	Weight      float32
	Record_time string
	Weekday     string
}

var chnNumChar = [10]string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}

var str_date_format = "2006-01-02"

func LinkMysql(db string) (res *gorm.DB) {
	//dsn := fmt.Sprintf("spidertian:Devil@86615@tcp(43.143.240.94:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", db)
	dsn := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", db)
	res, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print("db connect err~")
	}
	return
}

func main() {
	r := gin.Default()

	r.GET("/baidu", func(c *gin.Context) {
		/*
		if ret, err := http.NewRequest("GET", "https://mhaoma.baidu.com/pages/search-result/search-result?search=13693314520", nil); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(ret.Response)
		*/
		if resp, err := http.Get("https://mhaoma.baidu.com/pages/search-result/search-result?search=13693314520"); err != nil {
			fmt.Println(err.Error())
		} else {
			ret, _ := io.ReadAll(resp.Body)
			fmt.Printf("ret is %s", ret)
		}
	})

	r.GET("/bodyweight", func(c *gin.Context) {
		var db = LinkMysql("spidertian")
		var bw_record_list []Bw_record
		var cur_time = time.Now()
		var start_time = cur_time.AddDate(0, 0, -7).Format(str_date_format) + " 00:00:00"
		var end_time = cur_time.Format(str_date_format) + " 23:59:59"
		if err := db.Table("bw_record").Where("record_time BETWEEN ? AND ?", start_time, end_time).Order("record_time DESC").Find(&bw_record_list).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "数据查询失败",
			})
			return
		}
		for k, _ := range bw_record_list {
			record_time, _ := time.Parse("2006-01-02T15:04:05+08:00", bw_record_list[k].Record_time)
			bw_record_list[k].Weekday = fmt.Sprintf("星期%s", chnNumChar[int(record_time.Weekday())])
		}
		fmt.Println(time.Now())
		c.JSON(http.StatusOK, bw_record_list)
	})
	r.POST("/bodyweight", func(c *gin.Context) {
		var bw_record Bw_record
		// 绑定参数
		if err := c.BindJSON(&bw_record); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "上报数据格式有误. msg:",
				"err": err.Error(),
			})
			return
		}
		// 参数娇艳
		if _, err := time.Parse("2006-01-02 15:04:05", bw_record.Record_time); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "上报时间数据格式有误" + err.Error(),
			})
			return
		}
		var db = LinkMysql("spidertian")
		// 写入数据
		ret := db.Table("bw_record").Select("weight", "record_time").Create(&bw_record)
		if ret.Error != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "数据插入失败",
			})
			return
		}
		// 成功返回
		c.JSON(http.StatusOK, gin.H{
			"msg": "上报成功!",
		})
	})
	// 下面是活动相关的两个功能接口
	r.POST("/activity_info", func(c *gin.Context) {
		var req struct {
			Lat  float32
			Lon  float32
			Name string
			Tel  string
		}
		// 绑定参数
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "上报数据格式有误.",
				"status": 1,
			})
			return
		}
		if req.Lat == float32(0) || req.Lon == float32(0) {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "您的地址信息有误.",
				"status": 1,
			})
			return
		} else if req.Tel == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "手机号不能为空.",
				"status": 1,
			})
			return
		} else if match, _ := regexp.MatchString(`^1[3456789]\d{9}$`, req.Tel); !match {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "手机号格式不正确.",
				"status": 1,
			})
			return
		}
		res, err := http.Get(fmt.Sprintf("https://apis.map.qq.com/ws/geocoder/v1/?location=%f,%f&key=YPQBZ-7MTES-IP3O3-6FORR-PKBZF-IJBIX&get_poi=1", req.Lat, req.Lon))
		if err != nil || res.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "解析地址出错.",
				"status": 1,
			})
			return
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "解析地址出错2.",
				"status": 1,
			})
			return
		}
		var qq_api_res struct {
			Result struct {
				Address_component struct {
					City string
				}
				Formatted_addresses struct {
					Recommend string
				}
			}
		}
		if err := json.Unmarshal(body, &qq_api_res); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":    "第三方服务出错.",
				"status": 1,
			})
			return
		}
		// 写入数据
		var db = LinkMysql("micro_activity")
		type insert_data struct {
			Cn_name    *string
			Tel        string
			Latitude   float32
			Longtitude float32
			Address    string
		}
		ret := db.Table("user_info").Select("cn_name", "tel", "latitude", "longtitude", "address").Create(&insert_data{
			Cn_name:    &req.Name,
			Tel:        req.Tel,
			Latitude:   req.Lat,
			Longtitude: req.Lon,
			Address:    qq_api_res.Result.Address_component.City + qq_api_res.Result.Formatted_addresses.Recommend,
		})
		if ret.Error != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "数据插入失败",
			})
			return
		}
		// 成功返回
		c.JSON(http.StatusOK, gin.H{
			"msg":    "上报成功1!",
			"status": 0,
		})
	})
	r.GET("/activity_info", func(c *gin.Context) {
		var db = LinkMysql("micro_activity")
		var data []struct {
			Address    string
			Cn_name    string
			Tel        string
			Latitude   float32
			Longtitude float32
		}
		db.Table("user_info").Find(&data)
		// 成功返回
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})
	r.GET("/activity_detail", func(c *gin.Context) {
		req := c.Request.URL.Query()
		if len(req["product_id"]) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "产品id信息有误.",
			})
			return
		}
		if info, err := os.ReadFile(fmt.Sprintf("./activity_detail_file/detail_info_%s.json", req["product_id"][0])); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "获取不到产品信息.",
				"err": err.Error(),
			})
			return
		} else {
			var json_data struct {
				List []struct {
					Desc      []string `json:"desc"`
					Sub_title string   `json:"sub_title"`
				}
			}
			err := json.Unmarshal(info, &json_data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "获取产品内容有误.",
					"err": err.Error(),
				})
				return
			}
			// 成功返回
			c.JSON(http.StatusOK, gin.H{
				"msg":  "成功",
				"list": json_data.List,
			})
		}
	})
	r.Group("/admin").GET("/a", func(c *gin.Context) {
		// 用zerolog包写一条日志
		fileRes, err := os.OpenFile("/Users/tianye/go/src/github.com/spidergo/file/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("日志文件打不开啊", err)
			return
		}
		defer fileRes.Close()
		// zerolog初始化
		logger := zerolog.New(fileRes).With().Timestamp().Logger()
		// 记录一条没用的日志
		logger.Info().Dict("data", zerolog.Dict().Fields(map[string]interface{}{
			"msg":  "success",
			"list": map[string]int{"a": 1},
		})).Msg("first log msg is Hello World~~~")
		// ### end zerolog ###
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"list": map[string]int{"a": 1},
		})
	})
	//strings.ReplaceAll()
	endless.ListenAndServe("127.0.0.1:8080", r)
}
