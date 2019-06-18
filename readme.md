# golang websocket C10K 测试

## 环境配置

* 系统：ubuntu desktop
* 硬件：8C4G

### 参数限制调整
*  最大文件描述符的限制

    临时：ulimit -n (desktop版本，默认是1024)，临时设置最大只能到ulimit -n 2049
    永久：
    网上的解决方案是修改： /etc/security/limits.conf，然后重启，但是并不生效
```
* soft nofile 655350
* hard nofile 655350
* soft nproc 655350
* hard nproc 655350
```
    桌面版本:
    但是ubuntu（desktop)的配置是在/etc/systemd/system.conf，找到#DefaultLimitNOFILE=，去掉注释，等号后面输入最大数，
    推荐设置为65536

* 端口数限制

    一个ip地址，最多可以分配65535个端口，要达到10k的数量，可以在网卡上绑定多个ip，如：
    for i in `seq 1 9`; do sudo ifconfig eth0:$i 192.168.190.15$i up ; done
    这样就可以给eth0这个网卡绑定9个ip，从192.168.190.151-192.168.190.159

## 测试
### Server
```
go run ./server
```

### Client
```
# 分别不同的shell下执行, 每个client ，10000个连接
go run ./client 192.168.190.151:0
go run ./client 192.168.190.152:0
go run ./client 192.168.190.153:0
go run ./client 192.168.190.154:0
go run ./client 192.168.190.155:0
```
