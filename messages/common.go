package messages

import "errors"

var ErrBadSignature = errors.New("Bad signature")
var ErrInvalidAddress = errors.New("Not a valid address")
var ErrNotJSON = errors.New("Not a JSON string")
