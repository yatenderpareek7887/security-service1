// models/models.go
package models

import (
	"fmt"

	logDataModel "github.com/yatender-pareek/threat-analyzer-service/src/models/log-data-model"
	threatentity "github.com/yatender-pareek/threat-analyzer-service/src/models/threat-model"
)

func GetAllModels() []interface{} {
	models := []interface{}{
		&logDataModel.LogData{},
		&threatentity.Threat{},
	}
	fmt.Printf("Models: %+v\n", models)
	return models
}
