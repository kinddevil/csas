package eureka

import (
	"fmt"
	"github.com/twinj/uuid"
	"net"
	// "github.com/eriklupander/goeureka/util"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var eurekaAddr string

// https://github.com/Netflix/eureka/wiki/Eureka-REST-operations
func Register(ip string, port string, appName string, euredaAdd string) string {
	eurekaAddr = euredaAdd
	instanceId := GetUUID()

	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(dir + "/templates/regtpl.json")

	tpl := string(data)
	// tpl = strings.Replace(tpl, "${ipAddress}", GetLocalIP(), -1)
	// tpl = strings.Replace(tpl, "${port}", "6767", -1)
	// tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)
	tpl = strings.Replace(tpl, "${ipAddress}", ip, -1)
	tpl = strings.Replace(tpl, "${vipAddress}", ip, -1)
	tpl = strings.Replace(tpl, "${port}", port, -1)
	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)
	tpl = strings.Replace(tpl, "${appName}", appName, -1)
	fmt.Println(tpl)

	// Register.
	registerAction := HttpAction{
		Url:         eurekaAddr + "/eureka/apps/" + appName,
		Method:      "POST",
		ContentType: "application/json",
		Body:        tpl,
	}
	var result bool
	for {
		fmt.Println("register eureka...")
		result = DoHttpRequest(registerAction)
		fmt.Println(result)
		if result {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
	return ip + ":" + appName + ":" + instanceId
}

func StartHeartbeat(appName, appId string) {
	isRuning := true

	c := make(chan os.Signal, 1) // Create a channel accepting os.Signal
	// Bind a given os.Signal to the channel we just created
	signal.Notify(c, os.Interrupt)    // Register os.Interrupt
	signal.Notify(c, syscall.SIGTERM) // Register syscall.SIGTERM

	go func() { // Start an anonymous func running in a goroutine
		sig := <-c // that will block until a message is recieved on
		log.Println("system shutdown in eureka heartbeat by", sig.String(), "...")
		isRuning = false // deregistration and exit program.
	}()

	for isRuning {
		time.Sleep(time.Second * 30)
		heartbeat(appName, appId)
	}
}

func heartbeat(appName, appId string) {
	heartbeatAction := HttpAction{
		// Url:    "http://127.0.0.1:8761/eureka/apps/gotest/" + GetLocalIP() + ":gotest:" + instanceId,
		Url:    eurekaAddr + "/eureka/apps/" + appName + "/" + appId,
		Method: "PUT",
	}
	fmt.Println("eureka heartbeat...")
	ret := DoHttpRequest(heartbeatAction)
	fmt.Println(ret)
}

func Deregister(appName, appId string) {
	fmt.Println("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		// Url:    "http://127.0.0.1:8761/eureka/apps/gotest/" + GetLocalIP() + ":gotest:" + instanceId,
		Url:    eurekaAddr + "/eureka/apps/" + appName + "/" + appId,
		Method: "DELETE",
	}
	log.Println("system shutdown in eureka deregister...")
	ret := DoHttpRequest(deregisterAction)
	fmt.Println(ret)
	fmt.Println("Deregistered application, exiting. Check Eureka...")
}

func HandleSigterm(appName, appId string) {
	c := make(chan os.Signal, 1) // Create a channel accepting os.Signal
	// Bind a given os.Signal to the channel we just created
	signal.Notify(c, os.Interrupt)    // Register os.Interrupt
	signal.Notify(c, syscall.SIGTERM) // Register syscall.SIGTERM

	go func() { // Start an anonymous func running in a goroutine
		sig := <-c                 // that will block until a message is recieved on
		Deregister(appName, appId) // the channel. When that happens, perform Eureka
		log.Println("system shutdown by", sig.String(), "...")
		os.Exit(1) // deregistration and exit program.
	}()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetUUID() string {
	return uuid.NewV4().String()
}
