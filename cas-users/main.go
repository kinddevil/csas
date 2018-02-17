package main

import (
	"cas-users/api"
	"cas-users/service"
	// "os"
	// "dbclient"
	cfg "cas-users/config"
	"eureka"
	"fmt"
	"localLog"
	"log"
	"sync"
	// "time"
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

	go localLog.RegesterLog(envs["logPath"])

	log.Printf("Starting %v\n", appName)

	initiallizeMysqlClient(envs["mysqlUrl"])

	go startService(port) // Starts HTTP service  (async)
	log.Println("Starting HTTP service at " + port)
	fmt.Println("Starting HTTP service at " + port)

	if envs["eureka"] == "true" {

		appId := eureka.Register(ip, port, appName, eurekaAddr) // Performs Eureka registration
		log.Println(appId)

		eureka.HandleSigterm(appName, appId) // Handle graceful shutdown on Ctrl+C or kill

		go eureka.StartHeartbeat(appName, appId) // Performs Eureka heartbeating (async)
	}

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}

func startService(port string) {
	webserver.StartWebServerWithPrefix(port, &api.Routes, cfg.Prefix)
	// webserver.StartWebServer(port, &api.Routes)
}

func initiallizeMysqlClient(dbUrl string) {
	monce.Do(func() {
		api.MysqlClient = &service.UsersClient{}
	})
	api.MysqlClient.Init(dbUrl)
	api.MysqlClient.Seed()
}

func destructMysqlClient() {
	api.MysqlClient.Close()
}
