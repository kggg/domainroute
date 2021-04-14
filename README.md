# domainroute
解析域名对应的服务器IP列表， 然后将这些IP路由到某条ISP线路上， 主要是服务于有多条公网线路时使用。

## 配置文件在目录conf下
### app.ini
   设置app， 包含存储在文件或者mysql数据库
### route.ini
   这个文件保存需要解析的域名及将这个域名指向下一跳网关
   文件格式：<br>
      domainname via [gateway|nexthop]<br>
      例如:<br>
           www.qq.com via 192.168.1.1

####  存储在文件模式，  iplist目录保存解析到的域名地址和IP列表
该目录存放解析域名后的IP列表文件
