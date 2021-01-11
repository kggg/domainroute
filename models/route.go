package models

import (
	"bytes"
	"domainroute/errno"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

//HandleRoute 根据给定的地址addr及路由规则rule处理路由， 错误返回error
func HandleRoute(addr, rule string) error {
	err := addroute(addr, rule)
	return err
}

func checkroute(addr, table string) (bool, error) {
	cmd := &exec.Cmd{}
	cmd = exec.Command("/sbin/ip", "route", "list", "table", table)
	out, err := cmd.Output()
	if err != nil {
		now := time.Now().Format(timeLayout)
		return false, fmt.Errorf("%s check route of %s error: %w", now, addr, err)
	}

	content := bytes.Split(out, []byte("\n"))
	for _, v := range content {
		if bytes.Contains(v, []byte(addr)) {
			return true, nil
		}
	}
	return false, nil

}

//添加addr地址到路由表中， 如果rule中有指定路由表，则添加到指定的table中
func addroute(addr, rule string) error {
	cmd := &exec.Cmd{}
	tables, err := getRouteTables()
	if err != nil {
		return err
	}
	cmdstrs := strings.Split(rule, " ")
	now := time.Now().Format(timeLayout)
	if len(tables) >= 1 {
		for _, table := range tables {
			ok, err := checkroute(addr, table)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			if ok {
				return errno.ExistRoute
			}
			cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1], "table", table)
			_, err = cmd.Output()
			if err != nil {
				return fmt.Errorf("%s add %s route error: %w", now, addr, err)
			}
			fmt.Printf("%s add route %s for table %s successed\n", now, addr, table)
		}

	}
	ok, err := checkroute(addr, "main")
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	if ok {
		return errno.ExistRoute
	}
	cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1])
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s add %s route error: %w", now, addr, err)
	}
	fmt.Printf("%s add route %s for table main successed\n", now, addr)

	return nil
}

// 删除路由, 用于删除过期的路由
func delroute(addr string) error {
	cmd := &exec.Cmd{}
	tables, err := getRouteTables()
	if err != nil {
		return err
	}
	now := time.Now().Format(timeLayout)
	if len(tables) >= 1 {
		for _, table := range tables {
			ok, err := checkroute(addr, table)
			if err != nil {
				return fmt.Errorf("Error: %w", err)
			}
			if ok {
				cmd = exec.Command("/sbin/ip", "route", "del", addr, "table", table)
				_, err = cmd.Output()
				if err != nil {
					return fmt.Errorf("%s Delete %s route error: %w", now, addr, err)
				}
				fmt.Printf("%s Remove route for %s successed\n", now, addr)
			}
		}
	}
	ok, err := checkroute(addr, "main")
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	if ok {
		cmd = exec.Command("/sbin/ip", "route", "del", addr)
		_, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("%s Delete %s route error: %w", now, addr, err)
		}
		fmt.Printf("%s Remove route for %s successed\n", now, addr)
	}

	return nil
}
