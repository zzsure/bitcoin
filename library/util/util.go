package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
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
