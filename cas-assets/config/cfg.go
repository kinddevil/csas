package cfg

import "os"

const (
	Prefix = "/assets/v1"
)

var configs map[string]map[string]string = map[string]map[string]string{
	"local": nil,
	// default dev
	"dev": map[string]string{
		"env":        "dev",
		"port":       "5000",
		"appName":    "cas-assets",
		"eurekaAddr": "http://localhost:8761",
		// "eurekaAddr": "http://10.88.215.157",

		"localIp": "127.0.0.1",

		"mysqlUrl": "user:pass@tcp(localhost:3306)/cas?charset=utf8&parseTime=false",
		// "mysqlUrl": "user:pass@tcp(docker.for.mac.localhost:3306)/cas?charset=utf8&parseTime=false",

		// "logPath":    "log.log",

		"eureka": "false",
		// "eureka": "true",

		"ossBucket": "cas-assets",
	},
	"prod": map[string]string{
		"env":        "prod",
		"port":       "5000",
		"appName":    "cas-assets",
		"eurekaAddr": "http://10.16.156.36:8761",
		// "eurekaAddr": "http://localhost:8761",

		"localIp": "10.16.156.36",
		// "localIp":  "127.0.0.1",

		"mysqlUrl": "user:pass@tcp(10.16.156.36:3306)/cas?charset=utf8&parseTime=false",
		// "mysqlUrl": "user:pass@tcp(localhost:3306)/cas?charset=utf8&parseTime=false",

		// "logPath":    "log.log",
		"eureka": "true",

		"ossBucket": "cas-assets",
	},
}

func GetEnv() map[string]string {
	env := os.Getenv("env")
	port := os.Getenv("port")
	localIP := os.Getenv("localip")
	ossKey := os.Getenv("osskey")
	ossSec := os.Getenv("osssec")
	config := configs["dev"]
	if envs, ok := configs[env]; ok {
		config = envs
	}

	if port != "" {
		config["port"] = port
	}
	if localIP != "" {
		config["localIp"] = localIP
	}

	if ossKey == "" {
		panic("OSS key is empty")
	} else {
		config["ossKey"] = ossKey
	}

	if ossSec == "" {
		panic("OSS sec is empty")
	} else {
		config["ossSecret"] = ossSec
	}
	return config
}
