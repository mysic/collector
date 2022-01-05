package collector

func csvCollectorRunner(msg map[string]interface{}) error {
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

	//filePath := ""
	//file, err := os.Open(filePath)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//reader := csv.NewReader(file)
	return nil
}
