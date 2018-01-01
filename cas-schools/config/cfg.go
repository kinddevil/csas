package cfg

import "os"

const (
	Prefix = "/schools"
)

var configs map[string]map[string]string = map[string]map[string]string{
	"local": nil,
	// default dev
	"dev": map[string]string{
		"env":        "dev",
		"port":       "6768",
		"appName":    "cas-schools",
		"eurekaAddr": "http://localhost:8761",
		// "eurekaAddr": "http://10.88.215.157",

		"localIp": "127.0.0.1",

		"mysqlUrl": "user:pass@tcp(localhost:3306)/cas?charset=utf8&parseTime=false",
		// "mysqlUrl": "user:pass@tcp(docker.for.mac.localhost:3306)/cas?charset=utf8&parseTime=false",

		// "logPath":    "log.log",

		"eureka": "false",
		// "eureka": "true",
	},
	"prod": map[string]string{
		"env":        "prod",
		"port":       "6768",
		"appName":    "cas-schools",
		"eurekaAddr": "http://10.16.156.36:8761",
		// "eurekaAddr": "http://localhost:8761",

		"localIp": "10.16.156.36",
		// "localIp":  "127.0.0.1",

		"mysqlUrl": "user:pass@tcp(10.16.156.36:3306)/cas?charset=utf8&parseTime=false",
		// "mysqlUrl": "user:pass@tcp(localhost:3306)/cas?charset=utf8&parseTime=false",

		// "logPath":    "log.log",
		"eureka": "true",
	},
}

func GetEnv() map[string]string {
	env := os.Getenv("env")
	if envs, ok := configs[env]; ok {
		return envs
	} else {
		return configs["dev"]
	}
}
