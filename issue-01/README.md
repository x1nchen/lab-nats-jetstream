# issue-01

主要测试特性：在不同的 nats 连接里，使用 Subscribe 方法订阅相同的 Jetstream consumer subject，是否会拷贝多份数据

TLDR: 这种情况下会拷贝多份相同的数据给不同的 subscription

PS：注意用本地网卡 IP 替换 nats-server-*.conf 配置文件中的 routes 的 IP

example 

```bash
# macos
find . -name "nats-server-*.conf" | xargs -L 1 sed -i '' -e "s/10.81.7.58/192.168.50.19/g"
```

启动测试环境

```bash
docker-compose up -d 
```

启动测试及结果

```bash
ζ go run issue-01/main.go
consumer[2] [2021-09-15 15:52:45] received message: [0]
consumer[1] [2021-09-15 15:52:45] received message: [0]
consumer[0] [2021-09-15 15:52:45] received message: [0]
consumer[1] [2021-09-15 15:52:46] received message: [1]
consumer[0] [2021-09-15 15:52:46] received message: [1]
consumer[2] [2021-09-15 15:52:46] received message: [1]
consumer[1] [2021-09-15 15:52:47] received message: [2]
consumer[2] [2021-09-15 15:52:47] received message: [2]
consumer[0] [2021-09-15 15:52:47] received message: [2]
consumer[0] [2021-09-15 15:52:48] received message: [3]
consumer[1] [2021-09-15 15:52:48] received message: [3]
consumer[2] [2021-09-15 15:52:48] received message: [3]
consumer[2] [2021-09-15 15:52:49] received message: [4]
consumer[0] [2021-09-15 15:52:49] received message: [4]
consumer[1] [2021-09-15 15:52:49] received message: [4]
consumer[1] [2021-09-15 15:52:50] received message: [5]
consumer[0] [2021-09-15 15:52:50] received message: [5]
consumer[2] [2021-09-15 15:52:50] received message: [5]
consumer[0] [2021-09-15 15:52:51] received message: [6]
consumer[1] [2021-09-15 15:52:51] received message: [6]
consumer[2] [2021-09-15 15:52:51] received message: [6]
consumer[1] [2021-09-15 15:52:52] received message: [7]
consumer[0] [2021-09-15 15:52:52] received message: [7]
consumer[2] [2021-09-15 15:52:52] received message: [7]
consumer[0] [2021-09-15 15:52:53] received message: [8]
consumer[2] [2021-09-15 15:52:53] received message: [8]
consumer[1] [2021-09-15 15:52:53] received message: [8]
consumer[0] [2021-09-15 15:52:54] received message: [9]
consumer[2] [2021-09-15 15:52:54] received message: [9]
consumer[1] [2021-09-15 15:52:54] received message: [9]
```