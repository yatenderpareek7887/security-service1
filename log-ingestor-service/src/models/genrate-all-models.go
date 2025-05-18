// models/models.go
package models

import (
	"fmt"

	logDataentity "github.com/yatender-pareek/log-ingestor-service/src/models/log-data-model"
	userentity "github.com/yatender-pareek/log-ingestor-service/src/models/user-model"
)

func GetAllModels() []interface{} {
	models := []interface{}{
		&logDataentity.LogData{},
		&userentity.User{},
	}
	fmt.Printf("Models: %+v\n", models)
	return models
}
