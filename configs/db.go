package configs

import "time"


var minConn = 2
var maxConn = 10
var maxConnLifetime = 1 * time.Minute

func GetMinConns() int {
	return minConn
}

func GetMaxConns() int {
	return maxConn
}

func GetMaxConnLifetime() time.Duration {
	return maxConnLifetime
}
