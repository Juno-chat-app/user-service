package dto

import "github.com/Juno-chat-app/user-service/infra/logger"

type Converter interface {
	Convert(input interface{}) (output interface{}, err error)
}

func NewConverter(tInput interface{}, tOutput interface{}, logger logger.ILogger) Converter {
	panic("Not implemented")
}
