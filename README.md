# domainroute
解析域名对应的服务器IP列表， 然后将这些IP路由到某条ISP线路上， 主要是服务于有多条公网线路时使用。

## route.ini
这个文件保存需要解析的域名及将这个域名指向下一跳网关
文件格式：
  domainname via [gateway|nexthop]
  例如:
  www.qq.com via 192.168.1.1

## iplist
该目录存放解析域名后的IP列表文件
