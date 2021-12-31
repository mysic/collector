package collector

func csvCollectorRunner(msg map[string]interface{}) error {
	fileCollector := new(fileCollector)
	collector := csvCollector{fileCollector: fileCollector}
	err := collector.fileCollector.msgValid(msg)
	if err != nil {
		return err
	}
	collector.fileCollector.assignment(msg)
	collector.fileCollector.run()
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

func (csv *csvCollector) run() error {

	return nil
}
