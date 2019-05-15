## rproxy

Redis proxy server to be used during redis migration.

### Usage

**Help**
```sh
$ rproxy --help
Usage of rproxy:
  -addr string
    	Address on which rproxy will listen on (default ":6379")
  -dst string
    	Address of destination redis: Example: redis-clusterdev.example.com:6379
  -phase int
    	Migration phase. Possible values:
    		(0: Write to both redis; Read from src redis)(default)
    		(1: Write to both redis; Read from dst redis)
    		(2: Write to and read from dst redis only)
  -src string
    	Address of source redis: Example: redis-nonclusterdev.example.com:6379
```

**Run the proxy**
```sh
$ rproxy -src=127.0.0.1:6380 -dst=127.0.0.1:6381 -phase=0

2019/05/15 11:39:08 Source redis reachable at 127.0.0.1:6380
2019/05/15 11:39:08 Destination redis reachable at 127.0.0.1:6381
2019/05/15 11:39:08 Chosen migration phase 0: Write to both redis; Read from src redis
2019/05/15 11:39:08 Serving at :6379
```
