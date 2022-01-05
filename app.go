package main

import (
	"encoding/json"
	"etl/collector"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	Timezone = "Asia/Shanghai"
	Address  = "/tmp/collector.sock"
	Network  = "unix"
	LogFile  = "/var/log/etl/collector.log"
)

func main() {
	location, err := time.LoadLocation(Timezone)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	time.Local = location
	defer func() {
		log.Println("collector exiting...")
		_, err := os.Stat(Address)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			log.Println("unlink sock file err: " + err.Error())
			if os.IsExist(err) {
				err = syscall.Unlink(Address)
			}
		}
		log.Println("collector exited")
	}()
	exitSigs := make(chan os.Signal, 1)
	doExit := make(chan bool, 1)
	signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exitSigs
		doExit <- true
	}()
	lf, _ := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	log.SetOutput(lf)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	/**
	logics start
	*/
	log.Println("collector start")
	socket, err := net.Listen(Network, Address)
	if err != nil {
		log.Println("socket listen err: " + err.Error())
		return
	}
	defer func(socket net.Listener) {
		log.Println("socket closing...")
		err := socket.Close()
		if err != nil {
			log.Println("socket closing error: " + err.Error())
			return
		}
		log.Println("socket closed")
	}(socket)
	go func() {
		log.Println("socket waiting for request...")
		for {
			clientSock, err := socket.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				continue
			}
			go func() {
				for {
					buf := make([]byte, 1024)
					dataLen, err := clientSock.Read(buf)
					if err != nil {
						if err.Error() == "EOF" {
							break
						}
						wrong := []byte("error: invalid data.")
						_, _ = clientSock.Write(wrong)
						break
					}
					data := buf[0:dataLen]
					msg := string(data)
					log.Println(msg)
					if json.Valid([]byte(msg)) != true {
						wrong := []byte("error: invalid json format")
						_, _ = clientSock.Write(wrong)
						break
					}
					go func() {
						err := collector.Run(msg)
						if err != nil {
							log.Println("init collector err: " + err.Error())
						}
					}()
					response := []byte("ok")
					_, _ = clientSock.Write(response)

				}
				err = clientSock.Close()
				if err != nil {
					log.Println("client sock close err: " + err.Error())
				}
			}()
		}
	}()
	/**
	logics end
	*/

	if <-doExit {
		return
	}

}
