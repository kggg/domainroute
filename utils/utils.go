package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

//TimeConversion conversion t to unix timestamp
func TimeConversion(t string) (int64, error) {
	times, err := time.Parse(TimeLayout, t)
	if err != nil {
		return 0, fmt.Errorf("Convert tiem error:%w", err)
	}
	timeUnix := times.Unix()
	return timeUnix, nil
}

// ReadFromFile 读文件dname, 获取里面的内容返回[]byte, 错误返回error。
// 主要目的是与新解析出来的域名和地址列表对照。
func ReadFromFile(filename string) ([]string, error) {
	var content []string

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Read file %s error: %w", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineText := scanner.Text()
		content = append(content, lineText)

	}
	return content, nil
}

func CheckFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//formatParser 检测route.ini格式是否正确
func FormatParser(content []string) error {
	for _, v := range content {
		sub := strings.Fields(v)
		//检测域名格式
		var m = regexp.MustCompile("^[a-zA-Z0-9\\-\\.]+(\\.[a-z]{2,4})$")
		match := m.MatchString(sub[0])
		if !match {
			return fmt.Errorf("Validate route.ini domain: invalid [%s]", sub[0])
		}
		//检测中间字符是否是via
		if sub[1] != "via" {
			return fmt.Errorf("Validate file route.ini format error: invalid [%s]", sub[1])
		}
		//检测最后一格 的IP地址
		ok := net.ParseIP(sub[2])
		if ok == nil {
			return fmt.Errorf("validate route.ini ipaddress error [%s]", sub[2])
		}

	}
	return nil
}

// Resolv 解析dname得到IP地址返回， 错误返回error
func Resolv(dname string) ([]string, error) {
	dname = strings.TrimSpace(dname)
	iplist, err := net.LookupHost(dname)
	if err != nil {
		return nil, fmt.Errorf("Resolv %s error: %w", dname, err)
	}
	var newslice []string
	for i := 0; i < len(iplist); i++ {
		if strings.Contains(iplist[i], ":") {
			continue
		}
		newslice = append(newslice, iplist[i])
	}
	return newslice, nil
}
