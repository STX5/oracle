# 一、预言机架构图
![img.png](./img.png)
# 二、格式约定
* `LIST ALL KEYS`: etcdctl  --endpoints=localhost:2379 get --prefix --keys-only ''
* `PUT JOB`: etcdctl  --endpoints=localhost:2379 put 011 '{"url":"http://www.baidu.com","pattern":" 1"}'