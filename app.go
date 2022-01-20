package main

import (
	"encoding/json"
	"etl/collector"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	Network  = "unix"
	Timezone = "Asia/Shanghai"
	PidFile  = "/tmp/collector.pid"
	SockFile = "/tmp/collector.sock"
	LogFile  = "/var/log/collector.log"
)

func main() {
	_, err := os.Stat(PidFile)
	if err == nil {
		fmt.Println("collector is running...")
		os.Exit(0)
	}
	pidFile, err := os.OpenFile(PidFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
	}
	_, err = pidFile.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		return
	}
	err = pidFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		return
	}
	defer func() {
		err = os.Remove(PidFile)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	location, err := time.LoadLocation(Timezone)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	time.Local = location
	defer func() {
		log.Println("collector exiting...")
		_, err := os.Stat(SockFile)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			log.Println("unlink sock file err: " + err.Error())
			if os.IsExist(err) {
				err = syscall.Unlink(SockFile)
			}
		}
		log.Println("collector exited")
		fmt.Println("collector exited")
	}()
	exitSigs := make(chan os.Signal, 1)
	doExit := make(chan bool, 1)
	signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		close(exitSigs)
		close(doExit)
	}()
	go func() {
		<-exitSigs
		doExit <- true
	}()
	lf, _ := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	log.SetOutput(lf)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.Println("collector start")
	fmt.Println("collector start")
	socket, err := net.Listen(Network, SockFile)
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
			log.Println("request accept...")
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
					log.Println("request msg - " + msg)
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
	defer func() {
		if r := recover(); r != nil {
			errStr, ok := r.(string)
			if ok {
				log.Println("panic err: " + errStr)
			}
			err, ok = r.(error)
			if ok {
				log.Println("panic err: " + err.Error())
			}
		}
	}()
	if <-doExit {
		return
	}

}
