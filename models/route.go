package models

import (
	"bytes"
	"domainroute/myerrors"
	"fmt"
	"os/exec"
	"strings"
)

func HandleRoute(addr, rule string) error {
	ok, err := checkroute(addr, rule)
	if err != nil {
		return fmt.Errorf("checkroute error: %w\n", err)
	}
	if ok {
		return myerrors.ExistRoute
	}
	err = addroute(addr, rule)
	return err
}

func checkroute(addr, rule string) (bool, error) {
	cmd := &exec.Cmd{}
	if strings.Contains(rule, "table") {
		table := strings.Split(rule, "table")[1]
		table = strings.TrimSpace(table)
		cmd = exec.Command("/sbin/ip", "route", "list", "table", table)
	} else {
		cmd = exec.Command("/sbin/ip", "route", "list")
	}
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
	cmdstrs := strings.Split(rule, " ")
	if strings.Contains(rule, "table") {
		cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1], "table", cmdstrs[3])
	} else {
		cmd = exec.Command("/sbin/ip", "route", "add", addr, "via", cmdstrs[1])
	}
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("add %s route error: %w", addr, err)
	}
	fmt.Printf("add route for %s successed\n", addr)
	return nil
}

// 删除路由, 用于删除过期的路由
func delroute(addr string) error {
	return nil
}
