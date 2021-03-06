package collector

import (
	"errors"
	"fmt"
	scp "github.com/hnakamur/go-scp"
	"golang.org/x/crypto/ssh"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
)

func fileCollectorRunner(msg map[string]interface{}) error {
	collector := new(fileCollector)
	err := collector.msgValid(msg)
	if err != nil {
		return err
	}
	collector.assignment(msg)
	return collector.run()
}

type fileCollector struct {
	data file
}

func (f *fileCollector) run() error {
	log.Println("fileCollector init...")
	for _, val := range f.data.paths {
		go func(val *filePath) {
			defer func() {
				err := recover()
				log.Println(err)
				log.Println(val)
			}()
			val.account = strings.Trim(val.account, " \r\n\t")
			val.pwd = strings.Trim(val.pwd, " \r\n\t")
			val.host = strings.Trim(val.host, " \r\n\t")
			val.path = strings.Trim(val.path, " \r\n\t")
			val.kind = strings.Trim(val.kind, " \r\n\t")
			val.saveTo = strings.Trim(val.saveTo, " \r\n\t")
			hostParams, err := url.Parse(val.host)
			log.Println("host parsed - ", hostParams)
			if err != nil {
				fmt.Println(err.Error())
			}
			config := &ssh.ClientConfig{
				User: val.account,
				Auth: []ssh.AuthMethod{
					ssh.Password(val.pwd), //todo 管理端后台前端表单提交前需要做crypto加密
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(), //todo 学习ssh host_key,client_key
			}
			sshClient, err := ssh.Dial(hostParams.Scheme, hostParams.Host, config)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer sshClient.Close()
			log.Println("ssh client created")
			// scp
			scpClient := scp.NewSCP(sshClient)
			// save path exist detect
			savePath := val.saveTo + string(os.PathSeparator) + f.data.sourceId
			isExist, err := f.pathExist(savePath)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if !isExist {
				log.Println("save path not exist,creating...")
				log.Println(savePath)
				err := os.MkdirAll(savePath, os.ModePerm)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
			log.Println("save path created - ", savePath)
			if val.kind == "dir" {
				err := scpClient.ReceiveDir(val.path, savePath, nil)
				if err != nil {
					log.Println(err.Error())
					fmt.Println(err.Error())
				}
				return
			}
			savePath += string(os.PathSeparator) + path.Base(val.path)
			err = scpClient.ReceiveFile(val.path, savePath)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}(val)
	}
	return nil
}

func (f *fileCollector) pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (f *fileCollector) assignment(msg map[string]interface{}) {
	var fileData file
	fileData.date = msg["date"].(string)
	fileData.sourceId = msg["source_id"].(string)
	fileData.sourceType = msg["source_type"].(string)

	paths := msg["paths"].([]interface{})
	for _, item := range paths {
		m := item.(map[string]interface{})
		filePathData := new(filePath)
		filePathData.host = m["host"].(string)
		filePathData.path = m["path"].(string)
		filePathData.kind = m["kind"].(string)
		filePathData.saveTo = m["saveTo"].(string)
		filePathData.account = m["account"].(string)
		filePathData.pwd = m["pwd"].(string)
		fileData.paths = append(fileData.paths, filePathData)
	}
	f.data = fileData
}

func (f *fileCollector) msgValid(msg map[string]interface{}) error {
	_, ok := msg["date"]
	if !ok {
		return errors.New("date field not exist")
	}
	if reflect.TypeOf(msg["date"]).String() != "string" {
		return errors.New("date field not string type")
	}

	_, ok = msg["source_id"]
	if !ok {
		return errors.New("source_id field not exist")
	}
	if reflect.TypeOf(msg["source_id"]).String() != "string" {
		return errors.New("source_id field not string type")
	}

	_, ok = msg["source_type"]
	if !ok {
		return errors.New("source_type field not exist")
	}
	if reflect.TypeOf(msg["source_type"]).String() != "string" {
		return errors.New("source_type field not string type")
	}

	_, ok = msg["paths"]
	if !ok {
		return errors.New("paths field not exist")
	}
	if reflect.TypeOf(msg["paths"]).String() != "[]interface {}" {
		return errors.New("source_type field not []interface {} type")
	}
	paths := msg["paths"].([]interface{})
	for _, item := range paths {
		p := item.(map[string]interface{})
		_, ok = p["host"]
		if !ok {
			return errors.New("paths.host field not exist")
		}
		if reflect.TypeOf(p["host"]).String() != "string" {
			return errors.New("paths.host field not string type")
		}
		_, ok = p["kind"]
		if !ok {
			return errors.New("paths.kind field not exist")
		}
		if reflect.TypeOf(p["kind"]).String() != "string" {
			return errors.New("paths.kind field not string type")
		}
		_, ok = p["path"]
		if !ok {
			return errors.New("paths.uri field not exist")
		}
		if reflect.TypeOf(p["path"]).String() != "string" {
			return errors.New("paths.path field not string type")
		}
		_, ok = p["saveTo"]
		if !ok {
			return errors.New("paths.saveTo field not exist")
		}
		if reflect.TypeOf(p["saveTo"]).String() != "string" {
			return errors.New("paths.saveTo field not string type")
		}
		_, ok = p["account"]
		if !ok {
			return errors.New("paths.account field not exist")
		}
		if reflect.TypeOf(p["account"]).String() != "string" {
			return errors.New("paths.account field not string type")
		}
		_, ok = p["pwd"]
		if !ok {
			return errors.New("paths.pwd field not exist")
		}
		if reflect.TypeOf(p["pwd"]).String() != "string" {
			return errors.New("paths.pwd field not string type")
		}
	}
	return nil
}
