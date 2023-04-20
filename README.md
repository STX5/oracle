# 一、预言机架构图
![](doc/img.png)
# 二、常用命令
* `LIST ALL KEYS`: etcdctl  --endpoints=localhost:2379 get --prefix --keys-only ''
* `PUT JOB`: etcdctl  --endpoints=localhost:2379 put 011 '{"url":"http://www.baidu.com","pattern":" 1"}'
# 三、Quick Start
## 1. run a ETCD cluster
```sh
# your etcd executable directory
./etcd
```
## 2. run jobDeamon
```sh
cd cmd/jobDeamon
go run .
```
## 3. run worker
```sh
cd cmd/worker
go run .
```
## 4. put a job
```sh
# your etcd executable directory
./etcdctl  --endpoints=localhost:2379 put 011 '{"url":"http://www.baidu.com","pattern":" 1"}'

./etcdctl  --endpoints=localhost:2379 get --prefix --keys-only ''
```
This project is bulit for fault tolerance. You can have multiple `wokers` and `jobDeamons` running, and feel free to shut them down at any time(eg. shut a worker down when its working; shut a jobDeamon down when it detects a lock release event).

As long as there are at least one woker and one jobDeamon, every job will be done exactly once.