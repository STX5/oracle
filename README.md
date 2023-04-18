* `LIST ALL KEYS`: etcdctl  --endpoints=localhost:2379 get --prefix --keys-only ''
* `PUT JOB`: etcdctl  --endpoints=localhost:2379 put 011 '{"url":"http://www.baidu.com","pattern":" 1"}'