package model

const (
	StatusOK             = 0
	StatusDataNotFound   = 10
	StatusDataDuplicated = 19
	StatusDataNotMatched = 21

	StatusBadRequest        = 40
	StatusUnauthorized      = 41 // RFC 7235, 3.1
	StatusPaymentRequired   = 42 // RFC 7231, 6.5.2
	StatusForbidden         = 43 // RFC 7231, 6.5.3
	StatusNotFound          = 44 // RFC 7231, 6.5.4
	StatusMethodNotAllowed  = 45 // RFC 7231, 6.5.5
	StatusNotAcceptable     = 46 // RFC 7231, 6.5.6
	StatusProxyAuthRequired = 47 // RFC 7235, 3.2
	StatusRequestTimeout    = 48 // RFC 7231, 6.5.7
	StatusConflict          = 49 // RFC 7231, 6.5.8

	StatusInternalServerError = 50
	StatusNotImplemented      = 51
	StatusBadGateway          = 52
	StatusServiceUnavailable  = 53
	StatusGatewayTimeout      = 54

	StatusEmailDuplicated = 191
)
