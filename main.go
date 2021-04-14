package main

import (
	"domainroute/app"
	"domainroute/errno"
	"domainroute/utils"
	"fmt"
	"sync"
)

func main() {
	config, err := app.NewConfig()
	if err != nil {
		panic(err)
	}
	// 1, 解析出域名列表, 并保存
	domains, err := config.GetDomain()
	if err != nil {
		panic(err)
	}
	//fmt.Println(domains)
	var wg sync.WaitGroup
	ch1 := make(chan struct{}, 5)
	for _, dname := range domains {
		ch1 <- struct{}{}
		domain := dname.Domainname
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			addr, err := utils.Resolv(domain)
			if err != nil {
				fmt.Println(err)
				return
			}
			switch config.Mode {
			case "file":
				// 解析出来的IP地址，加上时间点以域名为文件名存入文件中
				err = config.SaveToFile(domain, addr)
				if err != nil {
					fmt.Println(err)
					return
				}
			case "mysql":
				err = config.SaveToDB(domain, addr)
				if err != nil {
					fmt.Println(err)
					return
				}
			default:
				fmt.Println("Please setup app mode into conf/app.ini")
				return
			}
			<-ch1
		}(domain)
	}
	wg.Wait()

	// 2, 设置路由
	// parser route from file route.ini and generate rule

	ch2 := make(chan struct{}, 5)
	for _, line := range domains {
		ch2 <- struct{}{}
		wg.Add(1)
		go func(domain, gateway string) {
			defer wg.Done()
			var iplist []string
			switch config.Mode {
			case "file":
				iplist, err = config.ReadIPFormFile(domain)
				if err != nil {
					fmt.Println(err)
					return
				}
			case "mysql":
				err := config.ConnectionDB()
				if err != nil {
					fmt.Println(err)
					return
				}
				iplist, err = config.ReadIPFromDB(domain)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			for _, v := range iplist {
				err := config.HandleRoute(v, gateway)
				if err != nil {
					//重复的路由错误不需要打印
					if err == errno.ExistRoute {
						continue
					}
					fmt.Println(err)
				}
			}
			<-ch2
		}(line.Domainname, line.Gateway)

	}
	wg.Wait()

}
