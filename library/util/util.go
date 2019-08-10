package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"gitlab.azbit.cn/web/bitcoin/library/util/net"
)

func GetTodayDay() string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	day := time.Now().In(cstSh).Format("2006-01-02")
	return day
}

func GetTodayDayByUnix(s int64) string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	str := time.Unix(s, 0).In(cstSh).Format("2006-01-02T15:04:05Z07:00")
	return str
}

func GenUUID() string {
	return uuid.NewV4().String()
}

func GenServerUUID() string {
	ip, mac := net.NewLAN().NetInfo()
	return fmt.Sprintf("%s-%s", ip, mac)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func UnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func Float64Precision(f float64, prec int, round bool) float64 {
	pow10_n := math.Pow10(prec)
	if round {
		return math.Trunc(f+0.5/pow10_n) * pow10_n / pow10_n
	}
	return math.Trunc((f)*pow10_n) / pow10_n
}

func IntToString(i int) string {
	return fmt.Sprintf("%v", i)
}

func GetBackNum(num, divisor int) (int, int) {
	if divisor <= 0 {
		return 0, 0
	}
	if num < divisor {
		return 0, 0
	}
	d := num / divisor
	r := d % 10
	return d, r
}
