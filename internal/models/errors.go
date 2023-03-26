package models

import "fmt"

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

func ErrUrlParamNotFound(paramNames ...string) InternalError {
	return InternalError(fmt.Sprintf("cannot parse %v from URL", paramNames))
}

func ErrCannotParseJSON(dataTypeName string) InternalError {
	return InternalError(fmt.Sprintf("cannot parse %s from received JSON", dataTypeName))
}

func ErrDatabase() InternalError {
	return "internal database error"
}

func ErrDatabaseWrap(err error) InternalError {
	return InternalError(fmt.Errorf("internal database error: %w", err).Error())
}

func ErrNotEnoughMoney() InternalError {
	return "insufficient funds in the account balance"
}

func ErrSendingMsgWrap(err error) InternalError {
	return InternalError(fmt.Errorf("error while sending message: %w", err).Error())
}

func ErrReceivingMsgWrap(err error) InternalError {
	return InternalError(fmt.Errorf("error while receiving message: %w", err).Error())
}

func ErrEncodingMsgWrap(err error) InternalError {
	return InternalError(fmt.Errorf("error while encoding message: %w", err).Error())
}

func ErrDecodingMsgWrap(err error) InternalError {
	return InternalError(fmt.Errorf("error while decoding message: %w", err).Error())
}

func ErrCreateTopicWrap(err error) InternalError {
	return InternalError(fmt.Errorf("error while creating topic: %w", err).Error())
}
