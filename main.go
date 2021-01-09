package main

import (
	"fmt"
	"strings"
	"sync"

	"domainroute/errno"
	"domainroute/models"
	"domainroute/resolv"
)

func main() {
	// 1, 解析出域名列表, 并保存
	domainlist, err := models.ReadDomain()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	//限制goroutine的数量为5
	ch1 := make(chan struct{}, 5)
	ch2 := make(chan struct{}, 5)
	for _, dname := range domainlist {
		ch1 <- struct{}{}
		dname = strings.SplitN(dname, " ", 2)[0]
		//dname = strings.TrimSuffix(dname, "\n")
		wg.Add(1)
		go func(dname string) {
			defer wg.Done()

			//解析域名得到IP列表
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
			<-ch1
		}(dname)
	}
	wg.Wait()

	// 2, 设置路由
	// parser route from file route.ini and generate rule
	for _, line := range domainlist {
		ch2 <- struct{}{}
		content := strings.SplitN(line, " ", 2)
		wg.Add(1)
		go func(content []string) {

			defer wg.Done()
			iplist, err := models.ReadIPFormFile(content[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, v := range iplist {
				//需要检测路由表是否有重复的规则
				content[1] = strings.TrimSuffix(content[1], "\n")
				err := models.HandleRoute(v, content[1])
				if err != nil {
					//重复的路由错误不需要打印
					if err == errno.ExistRoute {
						continue
					}
					fmt.Println(err)
				}
			}
			<-ch2
		}(content)

	}
	wg.Wait()

}
