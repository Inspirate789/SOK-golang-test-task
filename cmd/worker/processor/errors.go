package processor

import "fmt"

type ProcessorError string

func (e ProcessorError) Error() string {
	return string(e)
}

func ErrFinishWrap(err error) ProcessorError {
	return ProcessorError(fmt.Errorf("kafka processor finished with error: %w", err).Error())
}
