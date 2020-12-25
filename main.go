package main

import (
	"domainroute/models"
	"domainroute/resolv"
	"fmt"
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
		dname = strings.TrimSuffix(dname, "\n")
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

	// add the route into system
}
