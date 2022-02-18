// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------

package comm

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/huichen/sego"
	"github.com/mozillazg/go-pinyin"
	"github.com/wengoldx/wing/logger"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

// Variables of Sego
var Segmenter sego.Segmenter

// Try try-catch-finaly method
func Try(do func(), catcher func(error), finaly ...func()) {
	defer func() {
		if i := recover(); i != nil {
			execption := errors.New(fmt.Sprint(i))
			logger.E("catched exception:", i)
			catcher(execption)
			if len(finaly) > 0 {
				finaly[0]()
			}
		}
	}()
	do()
}

// Condition return the trueData when pass the condition, or return falseData
//
// `USAGE` :
//
//	// use as follow to return diffrent type value, but the input
//	// true and false params MUST BE no-nil datas.
//	a := Condition(condition, trueString, falseString)	// return interface{}
//	b := Condition(condition, trueInt, falseInt).(int)
//	c := Condition(condition, trueInt64, falseInt64).(int64)
//	d := Condition(condition, trueFloat, falseFloat).(float64)
//	e := Condition(condition, trueDur, falseDur).(time.Duration)
//	f := Condition(condition, trueString, falseString).(string)
func Condition(condition bool, trueData interface{}, falseData interface{}) interface{} {
	if condition {
		return trueData
	}
	return falseData
}

// Contain check the given string list if contains item
func Contain(list *[]string, item string) bool {
	for _, v := range *list {
		if v == item {
			return true
		}
	}
	return false
}

// To2Digits fill zero if input digit not enough 2
func To2Digits(input interface{}) string {
	return fmt.Sprintf("%02d", input)
}

// To2Digits fill zero if input digit not enough 3
func To3Digits(input interface{}) string {
	return fmt.Sprintf("%03d", input)
}

// ToNDigits fill zero if input digit not enough N
func ToNDigits(input interface{}, n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", input)
}

// ToMap transform given struct data to map data, the transform struct
// feilds must using json tag to mark the map key.
//
// ---
//
//	type struct Sample {
//		Name string `json:"name"`
//	}
//	d := Sample{ Name : "name_value" }
//	md, _ := comm.ToMap(d)
//	// md data format is {
//	//     "name" : "name_value"
//	// }
func ToMap(input interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	buf, err := json.Marshal(input)
	if err != nil {
		logger.E("Marshal input struct err:", err)
		return nil, err
	}

	// json buffer decode to map
	d := json.NewDecoder(bytes.NewReader(buf))
	d.UseNumber()
	if err = d.Decode(&out); err != nil {
		logger.E("Decode json data to map err:", err)
		return nil, err
	}

	return out, nil
}

// ToXMLString transform given struct data to xml string
func ToXMLString(input interface{}) (string, error) {
	buf, err := xml.Marshal(input)
	if err != nil {
		logger.E("Marshal input to XML err:", err)
		return "", err
	}
	return string(buf), nil
}

// ToXMLReplace transform given struct data to xml string, ant then
// replace indicated fileds or values, to form param must not empty,
// but the to param allow set empty when use to remove all form keyworlds.
func ToXMLReplace(input interface{}, from, to string) (string, error) {
	xmlout, err := ToXMLString(input)
	if err != nil {
		return "", err
	}

	trimsrc := strings.TrimSpace(from)
	if trimsrc != "" {
		logger.D("Replace xml string from:", trimsrc, "to:", to)
		xmlout = strings.Replace(xmlout, trimsrc, to, -1)
	}
	return xmlout, nil
}

// JoinLines combine strings into multiple lines
func JoinLines(inputs ...string) string {
	packet := ""
	for _, line := range inputs {
		packet += line + "\n"
	}
	return packet
}

// GetSortKey get first letter of Chinese Pinyin
func GetSortKey(str string) string {
	if str == "" { // check the input param
		return "*"
	}

	// get the first char and verify if it is a~Z char
	firstChar, sortKey := []rune(str)[0], ""
	isAZchar, err := regexp.Match("[a-zA-Z]", []byte(str))
	if err != nil {
		logger.E("Regexp match err:", err)
		return "*"
	}

	if isAZchar {
		sortKey = string(unicode.ToUpper(firstChar))
	} else {
		if unicode.Is(unicode.Han, firstChar) { // chinese
			str1 := pinyin.LazyConvert(string(firstChar), nil)
			s := []rune(str1[0])
			sortKey = string(unicode.ToUpper(s[0]))
		} else if unicode.IsNumber(firstChar) { // number
			sortKey = string(firstChar)
		} else { // other language
			sortKey = "#"
		}
	}
	return sortKey
}

// GetKeyWords get more search keywords
func GetKeyWords(str string) []string {
	segments := Segmenter.Segment([]byte(str))

	// use search mode 'true' to get more search keywords
	words := sego.SegmentsToSlice(segments, true)
	var keywords []string
	for _, v := range words {
		if v == " " {
			continue
		}
		keywords = append(keywords, v)
	}
	return keywords
}

// RemoveDuplicate remove duplicate data from array
func RemoveDuplicate(oldArr []string) []string {
	newArr := make([]string, 0)
	for i := 0; i < len(oldArr); i++ {
		repeat := false
		for j := i + 1; j < len(oldArr); j++ {
			if oldArr[i] == oldArr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, oldArr[i])
		}
	}
	return newArr
}

// IgnoreSysSignalPIPE ignore system PIPE signal
func IgnoreSysSignalPIPE() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGPIPE)
	go func() {
		for {
			select {
			case sig := <-sc:
				if sig == syscall.SIGPIPE {
					logger.E("!! IGNORE BROKEN PIPE SIGNAL !!")
				}
			}
		}
	}()
}

// AccessAllowOriginBy allow cross domain access for the given origins
func AccessAllowOriginBy(category int, origins string, allowCredentials bool) {
	beego.InsertFilter("*", category, cors.Allow(&cors.Options{
		AllowAllOrigins:  !allowCredentials,
		AllowCredentials: allowCredentials,
		AllowOrigins:     []string{origins}, // use to set allow Origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))
}

// AccessAllowOriginByLocal allow cross domain access for localhost,
// the port number must config in /conf/app.conf file like :
//
// ---
//
//	; Server port of HTTP
//	httpport=3200
func AccessAllowOriginByLocal(category int, allowCredentials bool) {
	if beego.BConfig.Listen.HTTPPort > 0 {
		localhosturl := fmt.Sprintf("http://127.0.0.1:%v/", beego.BConfig.Listen.HTTPPort)
		AccessAllowOriginBy(category, localhosturl, allowCredentials)
	}
}

// ExecuteServer start and excute backend server
func ExecuteServer(allowCredentials ...bool) {
	IgnoreSysSignalPIPE()
	if len(allowCredentials) > 0 {
		AccessAllowOriginBy(beego.BeforeRouter, "*", allowCredentials[0])
		AccessAllowOriginBy(beego.BeforeStatic, "*", allowCredentials[0])
	} else {
		AccessAllowOriginBy(beego.BeforeRouter, "*", false)
		AccessAllowOriginBy(beego.BeforeStatic, "*", false)
	}

	// just output log to file on prod mode
	if beego.BConfig.RunMode != "dev" &&
		logger.GetLevel() != logger.LevelDebug {
		beego.BeeLogger.DelLogger(logs.AdapterConsole)
	}
	beego.Run()
}
