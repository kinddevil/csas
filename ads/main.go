package main

import (
	"ads/service"
	"dbclient"
	"eureka"
	"fmt"
	"log"
	"sync"
	"webserver"
)

// https://golang.org/pkg/strconv/
// i, err := strconv.Atoi("-42")
// s := strconv.Itoa(-42)
// b, err := strconv.ParseBool("true")
// f, err := strconv.ParseFloat("3.1415", 64)
// i, err := strconv.ParseInt("-42", 10, 64)
// u, err := strconv.ParseUint("42", 10, 64)
// s := "2147483647" // biggest int32
// i64, err := strconv.ParseInt(s, 10, 32)
// ...
// i := int32(i64)
// s := strconv.FormatBool(true)
// s := strconv.FormatFloat(3.1415, 'E', -1, 64)
// s := strconv.FormatInt(-42, 16)
// s := strconv.FormatUint(42, 16)
// q := Quote("Hello, 世界")
// q := QuoteToASCII("Hello, 世界")

func main() {

	port := "6767"
	appName := "cas-ads"
	eurekaAddr := "http://localhost:8761"
	ip := "127.0.0.1"
	fmt.Printf("Starting %v\n", appName)
	// initializeBoltClient()
	initiallizeMysqlClient()

	// handleSigterm()       // Handle graceful shutdown on Ctrl+C or kill
	go startService(port) // Starts HTTP service  (async)
	log.Println("Starting HTTP service at " + port)
	appId := eureka.Register(ip, port, appName, eurekaAddr) // Performs Eureka registration
	log.Println(appId)
	eureka.HandleSigterm(appName, appId) // Handle graceful shutdown on Ctrl+C or kill
	// handleSigterm("a")
	go eureka.StartHeartbeat(appName, appId) // Performs Eureka heartbeating (async)
	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()

	// startService("6767")
}

func startService(port string) {
	log.Println(service.Routes)
	webserver.StartWebServer(port, &service.Routes)
}

func initializeBoltClient() {
	service.DBClient = dbclient.GetDbClient()
	// service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}

func initiallizeMysqlClient() {
	service.MysqlClient = dbclient.InitMysql()
	service.MysqlClient.Open()
	service.MysqlClient.Seed()
}

func destructMysqlClient() {
	service.MysqlClient.Close()
}

// func Register() {
// 	instanceId = util.GetUUID() // Create a unique ID for this instance

// 	dir, _ := os.Getwd()
// 	data, _ := ioutil.ReadFile(dir + "/templates/regtpl.json") // Read registration JSON template file

// 	tpl := string(data)
// 	tpl = strings.Replace(tpl, "${ipAddress}", util.GetLocalIP(), -1) // Replace some placeholders
// 	tpl = strings.Replace(tpl, "${port}", "8080", -1)
// 	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)

// 	// Register.
// 	registerAction := HttpAction{ // Build a HttpAction struct
// 		Url:         "http://192.168.99.100:8761/eureka/apps/vendor", // Note hard-coded path to Eureka...
// 		Method:      "POST",
// 		ContentType: "application/json",
// 		Body:        tpl,
// 	}
// 	var result bool
// 	for {
// 		result = DoHttpRequest(registerAction) // Execute the HTTP request. result == true if req went OK
// 		if result {
// 			break // Success, end registration loop
// 		} else {
// 			time.Sleep(time.Second * 5) // Registration failed (usually, Eureka isn't up yet),
// 		} // retry in 5 seconds.
// 	}
// }
