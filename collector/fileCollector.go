package collector

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net/url"
	"os"
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
			val.account = strings.Trim(val.account, " \r\n\t")
			val.pwd = strings.Trim(val.pwd, " \r\n\t")
			val.host = strings.Trim(val.host, " \r\n\t")
			val.uri = strings.Trim(val.uri, " \r\n\t")
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
			log.Println("ssh config - ", config)
			sshClient, err := ssh.Dial(hostParams.Scheme, hostParams.Host, config)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer sshClient.Close()
			log.Println("ssh client created")
			// scp

		}(val)

	}
	// todo my blog

	//
	//session, err := sshClient.NewSession()
	//if err != nil {
	//	return err
	//}
	//defer func(session *ssh.Session) {
	//	err := session.Close()
	//	if err != nil && err.Error() != "EOF" {
	//	}
	//}(session)
	//var buf bytes.Buffer
	//session.Stdout = &buf
	//if err := session.Run("ls -al"); err != nil {
	//	return err
	//}
	//fmt.Println(buf.String())

	// sftp
	//newClient, err := sftp.NewClient(client)
	//if err != nil {
	//	return err
	//}

	// scp
	//var scpClient scp.Client
	//scpClient, err = scp.NewClientBySSH(sshClient)
	//if err != nil {
	//	return err
	//}
	//defer scpClient.Close()
	//var localFile os.File
	//localFile, err = os.Open(localFile)
	//defer localFile.Close()

	//err = scpClient.CopyFromRemote()

	return nil
}

func (f *fileCollector) isFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !fileInfo.IsDir(), nil
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
		filePathData.uri = m["uri"].(string)
		filePathData.saveTo = m["saveTo"].(string)
		filePathData.account = m["account"].(string)
		filePathData.pwd = m["pwd"].(string)
		suffix := m["suffix"].([]interface{})
		for _, suffixItem := range suffix {
			filePathData.suffix = append(filePathData.suffix, suffixItem.(string))
		}

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
		//_, ok = p["folder"]
		//if !ok {
		//	return errors.New("paths.folder field not exist")
		//}
		//if reflect.TypeOf(p["folder"]).String() != "[]interface {}" {
		//	return errors.New("paths.folder field not []interface {} type")
		//}
		//f := p["folder"].([]interface{})
		//if len(f) > 0 {
		//	for _, item = range f {
		//		if reflect.TypeOf(item).String() != "string" {
		//			return errors.New("paths.folder value not string type")
		//		}
		//	}
		//}
		_, ok = p["suffix"]
		if !ok {
			return errors.New("paths.suffix field not exist")
		}
		if reflect.TypeOf(p["suffix"]).String() != "[]interface {}" {
			return errors.New("paths.suffix field not []interface {} type")
		}
		s := p["suffix"].([]interface{})
		if len(s) > 0 {
			for _, item = range s {
				if reflect.TypeOf(item).String() != "string" {
					return errors.New("paths.suffix value not string type")
				}
			}
		}
		_, ok = p["uri"]
		if !ok {
			return errors.New("paths.uri field not exist")
		}
		if reflect.TypeOf(p["uri"]).String() != "string" {
			return errors.New("paths.uri field not string type")
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
