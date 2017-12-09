package cfg

import "os"

var configs map[string]map[string]string = map[string]map[string]string{
	"local": nil,
	// default dev
	"dev": map[string]string{
		"env":        "dev",
		"port":       "6767",
		"appName":    "cas-ads",
		"eurekaAddr": "http://localhost:8761",
		// "eurekaAddr": "http://10.88.215.157",
		"localIp":    "127.0.0.1",
		"mysqlUrl":   "user:pass@tcp(localhost:3306)/cas?charset=utf8&parseTime=false",
		// "logPath":    "log.log",
		"eureka": "false",
	},
	"prod": nil,
}

func GetEnv() map[string]string {
	env := os.Getenv("env")
	if envs, ok := configs[env]; ok {
		return envs
	} else {
		return configs["dev"]
	}
}
