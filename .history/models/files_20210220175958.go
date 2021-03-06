package models

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

//ReadDomain 从文件domainpath中读取需要解析的域名
func ReadDomain() ([]string, error) {
	content, err := readFromFile(domainpath)
	errA := formatParser(content)
	if errA != nil {
		return nil, fmt.Errorf("%w", errA)
	}
	return content, err
}

//formatParser 检测route.ini格式是否正确
func formatParser(content []string) error {
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

// SaveToFile 保存域名dname解析到的IP列表iplist到文件中
func SaveToFile(dname string, newiplist []string) error {
	filenamepath := iplistpath + "/" + dname

	ok := checkFileExists(filenamepath)
	if !ok {
		f, err := os.Create(filenamepath)
		if err != nil {
			return fmt.Errorf("Create file %s error: %w", dname, err)
		}
		defer f.Close()

		var ipstring string
		now := time.Now().Format(timeLayout)

		for _, v := range newiplist {
			ipstring += v + " " + now + "\n"
		}
		f.WriteString(ipstring)

	} else {
		//从文件中读取已经保存的地址信息
		oldiplist, err := readFromFile(filenamepath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		// 对比新解析到的地址和旧文件中保存的地址
		iplist, err := Compare(newiplist, oldiplist)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		var alliplist string
		for _, v := range iplist {
			alliplist += v + "\n"
		}
		//fmt.Println(filenamepath)
		err = ioutil.WriteFile(filenamepath, []byte(alliplist), 644)
		if err != nil {
			return fmt.Errorf("write to file %s error: %w", filenamepath, err)
		}
	}
	return nil
}

// ReadFromFile 读文件dname, 获取里面的内容返回[]byte, 错误返回error。
// 主要目的是与新解析出来的域名和地址列表对照。
func readFromFile(filename string) ([]string, error) {
	var content []string

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Read file %s error: %w", filename, err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("readline error: %w", err)
		}
		str = strings.TrimSuffix(str, "\n")
		if len(str) == 0 {
			continue
		}
		content = append(content, str)
	}
	return content, nil
}

// ReadIPFormFile 从存储文件中读取IP列表
func ReadIPFormFile(dname string) ([]string, error) {
	filenamepath := iplistpath + "/" + dname
	content, err := readFromFile(filenamepath)
	if err != nil {
		return nil, err
	}
	for i, v := range content {
		content[i] = strings.Split(v, " ")[0]
	}
	return content, nil
}

// Compare 对比新解析得到的IP列表与已经存在文件的IP列表， 如果IP已存在，则更新此IP的存储时间， 如果没有则追加到文件的末尾
func Compare(newiplist, oldiplist []string) ([]string, error) {
	ipmap := make(map[string]string)
	now := time.Now().Format(timeLayout)

	for _, v1 := range newiplist {
		ipmap[v1] = now
	}

	for _, v2 := range oldiplist {

		//如果之前的保存时间超过半年没有更新， 则去掉这个IP地址
		ptime := strings.SplitN(v2, " ", 2)[1]
		ptime = strings.TrimSuffix(ptime, "\n")
		pretime, err := timeConversion(ptime)
		if err != nil {
			continue
		}
		if time.Now().Unix()-pretime >= 15552000 {
			//同时要删除路由表中相关该地址的路由
			err := delroute(strings.Split(v2, " ")[0])
			if err != nil {
				return nil, err
			}
			continue
		}
		ipmap[strings.SplitN(v2, " ", 2)[0]] = ptime
	}

	var iplist []string
	for k, v := range ipmap {
		iplist = append(iplist, k+" "+v)
	}
	return iplist, nil
}

func checkFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 获取路由表总数
func getRouteTables() ([]string, error) {
	content, err := readFromFile(routeTablePath)
	if err != nil {
		return nil, err
	}
	var tables []string
	for _, v := range content {
		if strings.HasPrefix(v, "#") || strings.HasPrefix(v, "255") || strings.HasPrefix(v, "254") || strings.HasPrefix(v, "253") || strings.HasPrefix(v, "250") {
			continue
		}
		if len(v) == 0 || strings.HasPrefix(v, "0") {
			continue
		}
		tables = append(tables, strings.Split(v, " ")[1])
	}
	return tables, nil
}
