package apns

import (
	"fmt"
)

type ApnsStatus int8

const (
	ErrorNull ApnsStatus = iota
	ErrorProcessing
	ErrorMissingDevice
	ErrorMissingTopic
	ErrorMissingPayload
	ErrorInvalidTokenSize
	ErrorInvalidTopicSize
	ErrorInvalidPayloadSize
	ErrorInvalidToken
)

func (a ApnsStatus) String() string {
	switch a {
	case ErrorNull:
		return "No errors encountered"
	case ErrorProcessing:
		return "Processing error"
	case ErrorMissingDevice:
		return "Missing device token"
	case ErrorMissingTopic:
		return "Missing topic"
	case ErrorMissingPayload:
		return "Missing payload"
	case ErrorInvalidTokenSize:
		return "Invalid token size"
	case ErrorInvalidTopicSize:
		return "Invalid topic size"
	case ErrorInvalidPayloadSize:
		return "Invalid payload size"
	case ErrorInvalidToken:
		return "Invalid token"
	}
	return "Unknown"
}

type NotificationError struct {
	Command    uint8
	Status     ApnsStatus
	Identifier uint32

	OtherError error
}

// Make a new NotificationError with error response p and error err.
// If send in a 6-length p and non-nil err sametime, will ignore err and parse p.
func NewNotificationError(p []byte, err error) (e NotificationError) {
	if len(p) != 1+1+4 {
		if err != nil {
			e.OtherError = err
			return
		}
		e.OtherError = fmt.Errorf("Wrong data format, [%x]", p)
		return
	}
	e.Command = uint8(p[0])
	e.Status = ApnsStatus(p[1])
	e.Identifier = uint32(p[2])<<24 + uint32(p[3])<<16 + uint32(p[4])<<8 + uint32(p[5])
	return
}

func (e NotificationError) Error() string {
	if e.OtherError != nil {
		return e.OtherError.Error()
	}
	if e.Command != 8 {
		return fmt.Sprintf("Unknow error, command(%d), status(%d), id(%x)", e.Command, e.Status, e.Identifier)
	}
	return fmt.Sprintf("%s(%d): id(%x)", e.Status.String(), e.Status, e.Identifier)
}

func (e NotificationError) String() string {
	return e.Error()
}
