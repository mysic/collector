package collector

import (
	"errors"
	"reflect"
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
	//todo 通过SFTP文件从远程拉取到本地

	// todo my blog
	// ssh
	//config := &ssh.ClientConfig{
	//	User: "root",
	//	Auth: []ssh.AuthMethod{
	//		ssh.Password("eia3920@e208tset"),
	//	},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//}
	//sshClient, err := ssh.Dial("tcp", "47.119.146.28:22", config)
	//if err != nil {
	//	return err
	//}

	//session, err := client.NewSession()
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

func (f *fileCollector) assignment(msg map[string]interface{}) {
	var fileData file
	fileData.date = msg["date"].(string)
	fileData.sourceId = msg["source_id"].(string)
	fileData.sourceType = msg["source_type"].(string)
	// todo blog interface转换过程 需记录
	paths := msg["paths"].([]interface{})
	for _, item := range paths {
		m := item.(map[string]interface{})
		filePathData := new(filePath)
		filePathData.path = m["path"].(string)
		folder := m["folder"].([]interface{})
		for _, folderItem := range folder {
			filePathData.folder = append(filePathData.folder, folderItem.(string))
		}
		suffix := m["suffix"].([]interface{})
		for _, suffixItem := range suffix {
			filePathData.suffix = append(filePathData.suffix, suffixItem.(string))
		}
		filePathData.saveTo = m["saveTo"].(string)
		filePathData.account = m["account"].(string)
		filePathData.pwd = m["pwd"].(string)
		fileData.paths = append(fileData.paths, filePathData)
	}
	f.data = fileData

}

func (f *fileCollector) msgValid(msg map[string]interface{}) error {
	// todo blog interface转其他类型的一些经验
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
		_, ok = p["path"]
		if !ok {
			return errors.New("paths.path field not exist")
		}
		if reflect.TypeOf(p["path"]).String() != "string" {
			return errors.New("paths.path field not string type")
		}
		_, ok = p["folder"]
		if !ok {
			return errors.New("paths.folder field not exist")
		}
		if reflect.TypeOf(p["folder"]).String() != "[]interface {}" {
			return errors.New("paths.folder field not []interface {} type")
		}
		f := p["folder"].([]interface{})
		if len(f) > 0 {
			for _, item = range f {
				if reflect.TypeOf(item).String() != "string" {
					return errors.New("paths.folder value not string type")
				}
			}
		}
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
