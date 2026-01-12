package utils

import (
	"encoding/json"

	"go.uber.org/zap"
)

func InitLogger() (*zap.SugaredLogger, error) {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "info.log"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	logger := zap.Must(cfg.Build())

	return logger.Sugar(), nil
}
