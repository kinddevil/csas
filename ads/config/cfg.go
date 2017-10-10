package cfg

import "os"

var configs map[string]map[string]string = map[string]map[string]string{
	"local": nil,
	// default dev
	"dev": map[string]string{
		"port":       "6767",
		"appName":    "cas-ads",
		"eurekaAddr": "http://localhost:8761",
		"localIp":    "127.0.0.1",
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
