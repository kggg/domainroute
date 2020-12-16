package resolv

import (
	"fmt"
	"net"
	"strings"
)

func Resolv(dname string) ([]string, error) {
	dname = strings.TrimSpace(dname)
	iplist, err := net.LookupHost(dname)
	if err != nil {
		return nil, fmt.Errorf("Resolv %s error: %w", dname, err)
	}
	var newslice []string
	for i := 0; i < len(iplist); i++ {
		if strings.Contains(iplist[i], ":") {
			continue
		}
		newslice = append(newslice, iplist[i])
	}
	return newslice, nil
}
