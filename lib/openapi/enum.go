package openapi

import (
	"encoding/json"
)

//go:generate go run github.com/dmarkham/enumer@latest -type=Method -values
type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
	OPTIONS
	TRACE
)

type ParamIn int

//go:generate go run github.com/dmarkham/enumer@latest -type=ParamIn -values -transform=lower -json
const (
	PATH ParamIn = iota
	QUERY
	HEADER
)

//go:generate go run github.com/dmarkham/enumer@latest -type=StatusCode -values -trimprefix=Status
type StatusCode int

func (c StatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(c))
}

// Copied from [http] package.
const (
	StatusContinue           StatusCode = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols StatusCode = 101 // RFC 9110, 15.2.2
	StatusProcessing         StatusCode = 102 // RFC 2518, 10.1
	StatusEarlyHints         StatusCode = 103 // RFC 8297

	StatusOK                   StatusCode = 200 // RFC 9110, 15.3.1
	StatusCreated              StatusCode = 201 // RFC 9110, 15.3.2
	StatusAccepted             StatusCode = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo StatusCode = 203 // RFC 9110, 15.3.4
	StatusNoContent            StatusCode = 204 // RFC 9110, 15.3.5
	StatusResetContent         StatusCode = 205 // RFC 9110, 15.3.6
	StatusPartialContent       StatusCode = 206 // RFC 9110, 15.3.7
	StatusMultiStatus          StatusCode = 207 // RFC 4918, 11.1
	StatusAlreadyReported      StatusCode = 208 // RFC 5842, 7.1
	StatusIMUsed               StatusCode = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   StatusCode = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently  StatusCode = 301 // RFC 9110, 15.4.2
	StatusFound             StatusCode = 302 // RFC 9110, 15.4.3
	StatusSeeOther          StatusCode = 303 // RFC 9110, 15.4.4
	StatusNotModified       StatusCode = 304 // RFC 9110, 15.4.5
	StatusUseProxy          StatusCode = 305 // RFC 9110, 15.4.6
	_                       StatusCode = 306 // RFC 9110, 15.4.7 (Unused)
	StatusTemporaryRedirect StatusCode = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect StatusCode = 308 // RFC 9110, 15.4.9

	StatusBadRequest                   StatusCode = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                 StatusCode = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired              StatusCode = 402 // RFC 9110, 15.5.3
	StatusForbidden                    StatusCode = 403 // RFC 9110, 15.5.4
	StatusNotFound                     StatusCode = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed             StatusCode = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                StatusCode = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired            StatusCode = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout               StatusCode = 408 // RFC 9110, 15.5.9
	StatusConflict                     StatusCode = 409 // RFC 9110, 15.5.10
	StatusGone                         StatusCode = 410 // RFC 9110, 15.5.11
	StatusLengthRequired               StatusCode = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed           StatusCode = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        StatusCode = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong            StatusCode = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         StatusCode = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable StatusCode = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed            StatusCode = 417 // RFC 9110, 15.5.18
	StatusTeapot                       StatusCode = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest           StatusCode = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity          StatusCode = 422 // RFC 9110, 15.5.21
	StatusLocked                       StatusCode = 423 // RFC 4918, 11.3
	StatusFailedDependency             StatusCode = 424 // RFC 4918, 11.4
	StatusTooEarly                     StatusCode = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              StatusCode = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired         StatusCode = 428 // RFC 6585, 3
	StatusTooManyRequests              StatusCode = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   StatusCode = 451 // RFC 7725, 3

	StatusInternalServerError           StatusCode = 500 // RFC 9110, 15.6.1
	StatusNotImplemented                StatusCode = 501 // RFC 9110, 15.6.2
	StatusBadGateway                    StatusCode = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable            StatusCode = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout                StatusCode = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       StatusCode = 505 // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         StatusCode = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           StatusCode = 507 // RFC 4918, 11.5
	StatusLoopDetected                  StatusCode = 508 // RFC 5842, 7.2
	StatusNotExtended                   StatusCode = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired StatusCode = 511 // RFC 6585, 6
)
