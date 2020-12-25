package models

import "time"

type Address struct {
	Domainname string
	Ipaddr     string
	Created    time.Time
}

// NewAddress 返回一个Address的结构体
func NewAddress() *Address {
	return &Address{}
}

// Create 新增一个Address结构的数据
func (c *Address) Create(dname, ipaddr string, ctime time.Time) *Address {
	return &Address{Domainname: dname, Ipaddr: ipaddr, Created: ctime}
}

// Read 读取数据
func (c *Address) Read() {

}

// Save 保存以Address结构体形式保存数据
func (c *Address) Save() {

}

func (c *Address) Update() {

}

func (c *Address) Delete(addr string) {

}

func (c *Address) DeleteDomain(dname string) {

}
