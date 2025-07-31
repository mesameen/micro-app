package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	FieldService  = "service"
	FieldEndPoint = "endpoint"
)

var log *zap.Logger

func Init(serviceName string) error {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	log = log.With(zap.String(FieldService, serviceName))
	return nil
}

func Infof(format string, values ...interface{}) {
	log.Info(fmt.Sprintf(format, values...))
}

func Errorf(format string, values ...interface{}) {
	log.Error(fmt.Sprintf(format, values...))
}

func Panicf(format string, values ...interface{}) {
	log.Panic(fmt.Sprintf(format, values...))
}
