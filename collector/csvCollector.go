package collector

import "log"

func csvCollectorRunner(msg map[string]interface{}) error {
	log.Println("csvCollector init...")
	fileCollector := new(fileCollector)
	collector := csvCollector{fileCollector: fileCollector}
	err := collector.fileCollector.msgValid(msg)
	if err != nil {
		return err
	}
	collector.fileCollector.assignment(msg)
	err = collector.fileCollector.run()
	if err != nil {
		return err
	}
	err = collector.run()
	if err != nil {
		return err
	}
	return nil
}

type csvCollector struct {
	*fileCollector
	data file
}

func (c *csvCollector) run() error {
	log.Println("csvCollector run...")
	//filePath := ""
	//file, err := os.Open(filePath)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//reader := csv.NewReader(file)
	return nil
}
