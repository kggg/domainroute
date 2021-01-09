package models

import (
	"bytes"
	"domainroute/errno"
	"fmt"
	"os/exec"
	"strings"
)

func HandleRoute(addr, rule string) error {
	err := addroute(addr, rule)
	return err
}

func checkroute(addr, table string) (bool, error) {
	cmd := &exec.Cmd{}
	cmd = exec.Command("/sbin/ip", "route", "list", "table", table)
	out, err := cmd.Output()
	if err != nil {
		return false, err
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
	tables, err := parserRouteTables()
	if err != nil {
		return err
	}
	cmdstrs := strings.Split(rule, " ")
	if len(tables) >= 1 {
		for _, table := range tables {
			ok, err := checkroute(addr, table)
			if err != nil {
				return fmt.Errorf("checkroute error: %w\n", err)
			}
			if ok {
				return errno.ExistRoute
			}
			cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1], "table", table)
			_, err = cmd.Output()
			if err != nil {
				return fmt.Errorf("add %s route error: %w", addr, err)
			}
			fmt.Printf("add route %s for table %s successed\n", addr, table)
		}

	} else {
		ok, err := checkroute(addr, "main")
		if err != nil {
			return fmt.Errorf("checkroute error: %w\n", err)
		}
		if ok {
			return errno.ExistRoute
		}
		cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1])
		_, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("add %s route error: %w", addr, err)
		}
		fmt.Printf("add route %s for table main successed\n", addr)
	}
	return nil
}

// 删除路由, 用于删除过期的路由
func delroute(addr string) error {
	cmd := &exec.Cmd{}
	tables, err := parserRouteTables()
	if err != nil {
		return err
	}
	if len(tables) >= 1 {
		for _, table := range tables {
			ok, err := checkroute(addr, table)
			if err != nil {
				return fmt.Errorf("checkroute error: %w\n", err)
			}
			if ok {
				return errno.ExistRoute
			}
			cmd = exec.Command("/sbin/ip", "route", "del", addr, "table", table)
			_, err = cmd.Output()
			if err != nil {
				return fmt.Errorf("Delete %s route error: %w", addr, err)
			}
			fmt.Printf("Remove route for %s successed\n", addr)
		}
	} else {
		ok, err := checkroute(addr, "main")
		if err != nil {
			return fmt.Errorf("checkroute error: %w\n", err)
		}
		if ok {
			return errno.ExistRoute
		}
		cmd = exec.Command("/sbin/ip", "route", "del", addr)
		_, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("Delete %s route error: %w", addr, err)
		}
		fmt.Printf("Remove route for %s successed\n", addr)
	}
	return nil
}
