package models

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func HandleRoute(addr, rule string) error {
	if checkroute(addr, rule) {
		return fmt.Errorf("The address %s has been existing in route table\n", addr)
	}
	err := addroute(addr, rule)
	return err
}

func checkroute(addr, rule string) bool {
	cmdstr := "ip route list"
	if strings.Contains(rule, "table") {
		cmdstr = cmdstr + " " + strings.SplitAfter(rule, "table")[1]
	}
	out, err := run(cmdstr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	content := bytes.Split(out, []byte("\n"))
	for _, v := range content {
		if bytes.Contains(v, []byte(addr)) {
			return true
		}
	}
	return false

}

func addroute(addr, rule string) error {
	cmdstr := "ip route add " + addr + " " + rule
	if strings.Contains(rule, "table") {
		cmdstr = cmdstr + " " + strings.SplitAfter(rule, "table")[1]
	}
	_, err := run(cmdstr)
	if err != nil {
		return fmt.Errorf("add %s route error: %w", addr, err)
	}
	fmt.Printf("add route for %s successed\n", addr)
	return nil
}

func run(cmdstr string) ([]byte, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
