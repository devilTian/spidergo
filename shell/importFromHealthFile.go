/*
 * @Author: tianye@shimiotech.cn
 * @Date: 2023-12-18 22:15:53
 */
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var str_frag []string
	file, err := os.Open("./file/HKQuantityTypeIdentifierBodyMass.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dsn := "root:123456@tcp(127.0.0.1:3306)/spidertian?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("db connect err~")
	}
	exists_map := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	var dup_count int = 0

	db.Transaction(func(tx *gorm.DB) error {
		for scanner.Scan() {
			str_frag = strings.Split(scanner.Text(), ";")

			// todo 数据去掉重复 按startdate作为key
			if len(str_frag) != 8 || str_frag[0] == "type" {
				fmt.Println("过滤掉无用数据行[" + strings.Join(str_frag, "|") + "]")
				continue
			}
			if exists_map[str_frag[5]] == true {
				dup_count++
				fmt.Println("这条数据已记录过. 时间:[" + str_frag[5] + "] 体重:" + str_frag[7] + "Kg")
				continue
			}
			exists_map[str_frag[5]] = true

			weight, err := strconv.ParseFloat(str_frag[7], 32)
			if err != nil {
				log.Fatal(err)
				return err
			}
			record_time, err := time.Parse("2006-01-02 15:04:05 +0800", str_frag[4])
			if err != nil {
				log.Fatal(err)
				return err
			}
			if err := tx.Table("bw_record").Create(&Bw_record{
				Weight: float32(weight),
				Record_time : record_time.Format("2006-01-02 15:04:05"),
			}).Error; err != nil {
				return err
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
	fmt.Println(len(exists_map))
	fmt.Printf("重复的数据条数:[%d]", dup_count)
	
	
}
