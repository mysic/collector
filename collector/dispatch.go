package collector

import (
	"encoding/json"
	"errors"
	"log"
)

func Run(m string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("paniking: %s\r\n", r)
		}
	}()
	log.Println("dispatch request, init collector...")
	msg := make(map[string]interface{})
	err := json.Unmarshal([]byte(m), &msg)
	if err != nil {
		return err
	}
	switch msg["source_type"] {
	case "file":
		return fileCollectorRunner(msg)
	case "api":
		return apiCollectorRunner(msg)
	case "mysql":
		return mysqlCollectorRunner(msg)
	case "mongodb":
		return mongoDbCollectorRunner(msg)
	case "csv":
		return csvCollectorRunner(msg)
	default:
		str, _ := msg["source_type"].(string)
		return errors.New("source type isn't exist: " + str)
	}
}
