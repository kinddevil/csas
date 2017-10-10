package main

import (
	"ads/api"
	"ads/service"
	// "dbclient"
	"eureka"
	// "fmt"
	cfg "ads/config"
	"log"
	"sync"
	"webserver"
)

var (
	monce sync.Once
)

func main() {

	envs := cfg.GetEnv()
	port := envs["port"]
	appName := envs["appName"]
	eurekaAddr := envs["eurekaAddr"]
	ip := envs["localIp"]

	log.Printf("Starting %v\n", appName)

	initiallizeMysqlClient()

	go startService(port) // Starts HTTP service  (async)
	log.Println("Starting HTTP service at " + port)

	appId := eureka.Register(ip, port, appName, eurekaAddr) // Performs Eureka registration
	log.Println(appId)

	eureka.HandleSigterm(appName, appId) // Handle graceful shutdown on Ctrl+C or kill

	go eureka.StartHeartbeat(appName, appId) // Performs Eureka heartbeating (async)

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}

func startService(port string) {
	webserver.StartWebServer(port, &api.Routes)
}

func initiallizeMysqlClient() {
	monce.Do(func() {
		api.MysqlClient = &service.AdsClient{}
	})
	api.MysqlClient.Init("user:pass@tcp(localhost:3306)/db?charset=utf8&parseTime=true")
	api.MysqlClient.Seed()
}

func destructMysqlClient() {
	api.MysqlClient.Close()
}
