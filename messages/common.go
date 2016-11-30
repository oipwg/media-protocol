package messages

import "errors"

var ErrBadSignature = errors.New("Bad signature")
var ErrInvalidAddress = errors.New("Not a valid address")
var ErrInvalidReference = errors.New("Invalid reference transaction")
var ErrNotJSON = errors.New("Not a JSON string")
var ErrTooEarly = errors.New("Too early for a valid message")
var ErrWrongPrefix = errors.New("Wrong prefix for message type")
