package stun

import (
	"errors"
	"fmt"
)

// ErrorCodeAttribute represents ERROR-CODE attribute.
type ErrorCodeAttribute struct {
	Code   ErrorCode
	Reason []byte
}

func (c ErrorCodeAttribute) String() string {
	return fmt.Sprintf("%d: %s", c.Code, c.Reason)
}

// constants for ERROR-CODE encoding.
const (
	errorCodeReasonStart = 4
	errorCodeClassByte   = 2
	errorCodeNumberByte  = 3
	errorCodeReasonMaxB  = 763
	errorCodeModulo      = 100
)

// ErrReasonLengthTooBig means that len(Reason) > 763 bytes.
var ErrReasonLengthTooBig = errors.New("reason for ERROR-CODE is too big")

// AddTo adds ERROR-CODE to m.
func (c *ErrorCodeAttribute) AddTo(m *Message) error {
	value := make([]byte, 0, errorCodeReasonMaxB)
	if len(c.Reason) > errorCodeReasonMaxB {
		return ErrReasonLengthTooBig
	}
	value = value[:errorCodeReasonStart+len(c.Reason)]
	number := byte(c.Code % errorCodeModulo) // error code modulo 100
	class := byte(c.Code / errorCodeModulo)  // hundred digit
	value[errorCodeClassByte] = class
	value[errorCodeNumberByte] = number
	copy(value[errorCodeReasonStart:], c.Reason)
	m.Add(AttrErrorCode, value)
	return nil
}

// GetFrom decodes ERROR-CODE from m. Reason is valid until m.Raw is valid.
func (c *ErrorCodeAttribute) GetFrom(m *Message) error {
	v, err := m.Get(AttrErrorCode)
	if err != nil {
		return err
	}
	var (
		class  = uint16(v[errorCodeClassByte])
		number = uint16(v[errorCodeNumberByte])
		code   = int(class*errorCodeModulo + number)
	)
	c.Code = ErrorCode(code)
	c.Reason = v[errorCodeReasonStart:]
	return nil
}

// ErrorCode is code for ERROR-CODE attribute.
type ErrorCode int

// ErrNoDefaultReason means that default reason for provided error code
// is not defined in RFC.
var ErrNoDefaultReason = errors.New("No default reason for ErrorCode")

// AddTo adds ERROR-CODE with default reason to m. If there
// is no default reason, returns ErrNoDefaultReason.
func (c ErrorCode) AddTo(m *Message) error {
	reason := errorReasons[c]
	if reason == nil {
		return ErrNoDefaultReason
	}
	a := &ErrorCodeAttribute{
		Code:   c,
		Reason: reason,
	}
	return a.AddTo(m)
}

// Possible error codes.
const (
	CodeTryAlternate     ErrorCode = 300
	CodeBadRequest       ErrorCode = 400
	CodeUnauthorised     ErrorCode = 401
	CodeUnknownAttribute ErrorCode = 420
	CodeStaleNonce       ErrorCode = 428
	CodeRoleConflict     ErrorCode = 478
	CodeServerError      ErrorCode = 500
)

var errorReasons = map[ErrorCode][]byte{
	CodeTryAlternate:     []byte("Try Alternate"),
	CodeBadRequest:       []byte("Bad Request"),
	CodeUnauthorised:     []byte("Unauthorised"),
	CodeUnknownAttribute: []byte("Unknown Attribute"),
	CodeStaleNonce:       []byte("Stale Nonce"),
	CodeServerError:      []byte("Server Error"),
	CodeRoleConflict:     []byte("Role Conflict"),
}