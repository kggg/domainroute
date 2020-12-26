package main

import (
	"domainroute/models"
	"domainroute/resolv"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	// 1, 解析出域名列表, 并保存
	domainlist, err := models.ReadDomain()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup

	for _, dname := range domainlist {
		dname = strings.SplitN(dname, " ", 2)[0]
		//dname = strings.TrimSuffix(dname, "\n")
		wg.Add(1)
		go func(dname string) {
			defer wg.Done()
			addr, err := resolv.Resolv(dname)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 解析出来的IP地址，加上时间点以域名为文件名存入文件中
			err = models.SaveToFile(dname, addr)
			if err != nil {
				fmt.Println(err)
				return
			}

		}(dname)
	}
	wg.Wait()

	// 2, 设置路由
	// parser route from file route.ini and generate rule
	for _, line := range domainlist {
		dname := strings.SplitN(line, " ", 2)[0]
		rule := strings.SplitN(line, " ", 2)[1]
		rule = strings.TrimSuffix(rule, "\n")
		wg.Add(1)
		go func(dname string) {
			defer wg.Done()
			iplist, err := models.ReadIPFormFile(dname)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, v := range iplist {
				//fmt.Printf("ip route add %s %s\n", v, rule)
				cmd1 := exec.Command("/sbin/ip", "del", v)
				if _, err = cmd1.Output(); err != nil {
					fmt.Println(err)
					continue
				}
				cmd2 := exec.Command("/sbin/ip", "add", v, rule)
				out, err := cmd2.Output()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(string(out))
			}

		}(dname)

	}
	wg.Wait()

}
