package models

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//ReadDomain 从文件domainpath中读取需要解析的域名
func ReadDomain() ([]string, error) {
	content, err := readFromFile(domainpath)
	return content, err
}

// SaveToFile
func SaveToFile(dname string, newiplist []string) error {
	filenamepath := iplistpath + "/" + dname

	ok := checkFileExists(filenamepath)
	if !ok {
		f, err := os.Create(filenamepath)
		if err != nil {
			return fmt.Errorf("Create file %s error: %w\n", dname, err)
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
			return fmt.Errorf("%w\n", err)
		}
		// 对比新解析到的地址和旧文件中保存的地址
		iplist, err := Compare(newiplist, oldiplist)
		if err != nil {
			return fmt.Errorf("%w\n", err)
		}

		var alliplist string
		for _, v := range iplist {
			alliplist += v + "\n"
		}
		//fmt.Println(filenamepath)
		err = ioutil.WriteFile(filenamepath, []byte(alliplist), 644)
		if err != nil {
			return fmt.Errorf("write to file %s error: %w\n", filenamepath, err)
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
		return nil, fmt.Errorf("Read file %s error: %w\n", filename, err)
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
		content = append(content, str)
	}
	return content, nil
}

/*
type RouteAddress struct {
	Addr    string
	Created string
}

func (c *RouteAddress) Create(addr, created string) *RouteAddress {
	c.Addr = add
	c.Created = created
	return c
}

func (c *RouteAddress) Parser(dname string) ([]RouteAddress, error) {
	filenamepath := iplistpath + "/" + dname
	content, err := readFromFile(filenamepath)
	if err != nil {
		return nil, err
	}
	var raddress []RouteAddress
	for _, v := range content {
		res := strings.SplitN(v, " ", 2)
		res[1] = strings.TrimSuffix(res[1], "\n")
		raddress = append(raddress, RouteAddress{Addr: res[0], Created: res[1]})
	}
	return raddress, nil
}

*/

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
		fmt.Println("newip ", v1)
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
			continue
		}
		ipmap[strings.SplitN(v2, " ", 2)[0]] = ptime
	}

	var iplist []string
	for k, v := range ipmap {
		fmt.Printf("k is %s, time is %s\n", k, v)
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

func timeConversion(t string) (int64, error) {
	times, err := time.Parse(timeLayout, t)
	if err != nil {
		return 0, fmt.Errorf("Convert tiem error:%w\n", err)
	}
	timeUnix := times.Unix()
	return timeUnix, nil
}
