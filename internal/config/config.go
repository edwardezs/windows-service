package config

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type WindowsService struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	ParentExecPath    string   `json:"parentExecPath"`
	ChildExecPath     string   `json:"childExecPath"`
	ChildExecArgs     []string `json:"childExecArgs,omitempty"`
	LogFilePath       string   `json:"logFilePath,omitempty"`
	LogFileMaxSizeMB  int      `json:"logFileMaxSizeMB,omitempty"`
	LogFileMaxBackups int      `json:"logFileMaxBackups,omitempty"`
	LogFileMaxAgeDays int      `json:"logFileMaxAgeDays,omitempty"`
	LogFileCompress   bool     `json:"logFileCompress,omitempty"`
}

func New(filepath string) (cfg WindowsService, err error) {
	if err := configor.Load(&cfg, filepath); err != nil {
		return cfg, errors.Wrapf(err, "can not parse config file %s", filepath)
	}

	return cfg, nil
}
