package goapi

// StatusContinue represents http status code 100.
type StatusContinue struct{}

func (StatusContinue) statusCode() int {
	return 100
}

// StatusSwitchingProtocols represents http status code 101.
type StatusSwitchingProtocols struct{}

func (StatusSwitchingProtocols) statusCode() int {
	return 101
}

// StatusProcessing represents http status code 102.
type StatusProcessing struct{}

func (StatusProcessing) statusCode() int {
	return 102
}

// StatusEarlyHints represents http status code 103.
type StatusEarlyHints struct{}

func (StatusEarlyHints) statusCode() int {
	return 103
}

// StatusOK represents http status code 200.
type StatusOK struct{}

func (StatusOK) statusCode() int {
	return 200
}

// StatusCreated represents http status code 201.
type StatusCreated struct{}

func (StatusCreated) statusCode() int {
	return 201
}

// StatusAccepted represents http status code 202.
type StatusAccepted struct{}

func (StatusAccepted) statusCode() int {
	return 202
}

// StatusNonAuthoritativeInfo represents http status code 203.
type StatusNonAuthoritativeInfo struct{}

func (StatusNonAuthoritativeInfo) statusCode() int {
	return 203
}

// StatusNoContent represents http status code 204.
type StatusNoContent struct{}

func (StatusNoContent) statusCode() int {
	return 204
}

// StatusResetContent represents http status code 205.
type StatusResetContent struct{}

func (StatusResetContent) statusCode() int {
	return 205
}

// StatusPartialContent represents http status code 206.
type StatusPartialContent struct{}

func (StatusPartialContent) statusCode() int {
	return 206
}

// StatusMultiStatus represents http status code 207.
type StatusMultiStatus struct{}

func (StatusMultiStatus) statusCode() int {
	return 207
}

// StatusAlreadyReported represents http status code 208.
type StatusAlreadyReported struct{}

func (StatusAlreadyReported) statusCode() int {
	return 208
}

// StatusIMUsed represents http status code 226.
type StatusIMUsed struct{}

func (StatusIMUsed) statusCode() int {
	return 226
}

// StatusMultipleChoices represents http status code 300.
type StatusMultipleChoices struct{}

func (StatusMultipleChoices) statusCode() int {
	return 300
}

// StatusMovedPermanently represents http status code 301.
type StatusMovedPermanently struct{}

func (StatusMovedPermanently) statusCode() int {
	return 301
}

// StatusFound represents http status code 302.
type StatusFound struct{}

func (StatusFound) statusCode() int {
	return 302
}

// StatusSeeOther represents http status code 303.
type StatusSeeOther struct{}

func (StatusSeeOther) statusCode() int {
	return 303
}

// StatusNotModified represents http status code 304.
type StatusNotModified struct{}

func (StatusNotModified) statusCode() int {
	return 304
}

// StatusUseProxy represents http status code 305.
type StatusUseProxy struct{}

func (StatusUseProxy) statusCode() int {
	return 305
}

// StatusTemporaryRedirect represents http status code 307.
type StatusTemporaryRedirect struct{}

func (StatusTemporaryRedirect) statusCode() int {
	return 307
}

// StatusPermanentRedirect represents http status code 308.
type StatusPermanentRedirect struct{}

func (StatusPermanentRedirect) statusCode() int {
	return 308
}

// StatusBadRequest represents http status code 400.
type StatusBadRequest struct{}

func (StatusBadRequest) statusCode() int {
	return 400
}

// StatusUnauthorized represents http status code 401.
type StatusUnauthorized struct{}

func (StatusUnauthorized) statusCode() int {
	return 401
}

// StatusPaymentRequired represents http status code 402.
type StatusPaymentRequired struct{}

func (StatusPaymentRequired) statusCode() int {
	return 402
}

// StatusForbidden represents http status code 403.
type StatusForbidden struct{}

func (StatusForbidden) statusCode() int {
	return 403
}

// StatusNotFound represents http status code 404.
type StatusNotFound struct{}

func (StatusNotFound) statusCode() int {
	return 404
}

// StatusMethodNotAllowed represents http status code 405.
type StatusMethodNotAllowed struct{}

func (StatusMethodNotAllowed) statusCode() int {
	return 405
}

// StatusNotAcceptable represents http status code 406.
type StatusNotAcceptable struct{}

func (StatusNotAcceptable) statusCode() int {
	return 406
}

// StatusProxyAuthRequired represents http status code 407.
type StatusProxyAuthRequired struct{}

func (StatusProxyAuthRequired) statusCode() int {
	return 407
}

// StatusRequestTimeout represents http status code 408.
type StatusRequestTimeout struct{}

func (StatusRequestTimeout) statusCode() int {
	return 408
}

// StatusConflict represents http status code 409.
type StatusConflict struct{}

func (StatusConflict) statusCode() int {
	return 409
}

// StatusGone represents http status code 410.
type StatusGone struct{}

func (StatusGone) statusCode() int {
	return 410
}

// StatusLengthRequired represents http status code 411.
type StatusLengthRequired struct{}

func (StatusLengthRequired) statusCode() int {
	return 411
}

// StatusPreconditionFailed represents http status code 412.
type StatusPreconditionFailed struct{}

func (StatusPreconditionFailed) statusCode() int {
	return 412
}

// StatusRequestEntityTooLarge represents http status code 413.
type StatusRequestEntityTooLarge struct{}

func (StatusRequestEntityTooLarge) statusCode() int {
	return 413
}

// StatusRequestURITooLong represents http status code 414.
type StatusRequestURITooLong struct{}

func (StatusRequestURITooLong) statusCode() int {
	return 414
}

// StatusUnsupportedMediaType represents http status code 415.
type StatusUnsupportedMediaType struct{}

func (StatusUnsupportedMediaType) statusCode() int {
	return 415
}

// StatusRequestedRangeNotSatisfiable represents http status code 416.
type StatusRequestedRangeNotSatisfiable struct{}

func (StatusRequestedRangeNotSatisfiable) statusCode() int {
	return 416
}

// StatusExpectationFailed represents http status code 417.
type StatusExpectationFailed struct{}

func (StatusExpectationFailed) statusCode() int {
	return 417
}

// StatusTeapot represents http status code 418.
type StatusTeapot struct{}

func (StatusTeapot) statusCode() int {
	return 418
}

// StatusMisdirectedRequest represents http status code 421.
type StatusMisdirectedRequest struct{}

func (StatusMisdirectedRequest) statusCode() int {
	return 421
}

// StatusUnprocessableEntity represents http status code 422.
type StatusUnprocessableEntity struct{}

func (StatusUnprocessableEntity) statusCode() int {
	return 422
}

// StatusLocked represents http status code 423.
type StatusLocked struct{}

func (StatusLocked) statusCode() int {
	return 423
}

// StatusFailedDependency represents http status code 424.
type StatusFailedDependency struct{}

func (StatusFailedDependency) statusCode() int {
	return 424
}

// StatusTooEarly represents http status code 425.
type StatusTooEarly struct{}

func (StatusTooEarly) statusCode() int {
	return 425
}

// StatusUpgradeRequired represents http status code 426.
type StatusUpgradeRequired struct{}

func (StatusUpgradeRequired) statusCode() int {
	return 426
}

// StatusPreconditionRequired represents http status code 428.
type StatusPreconditionRequired struct{}

func (StatusPreconditionRequired) statusCode() int {
	return 428
}

// StatusTooManyRequests represents http status code 429.
type StatusTooManyRequests struct{}

func (StatusTooManyRequests) statusCode() int {
	return 429
}

// StatusRequestHeaderFieldsTooLarge represents http status code 431.
type StatusRequestHeaderFieldsTooLarge struct{}

func (StatusRequestHeaderFieldsTooLarge) statusCode() int {
	return 431
}

// StatusUnavailableForLegalReasons represents http status code 451.
type StatusUnavailableForLegalReasons struct{}

func (StatusUnavailableForLegalReasons) statusCode() int {
	return 451
}

// StatusInternalServerError represents http status code 500.
type StatusInternalServerError struct{}

func (StatusInternalServerError) statusCode() int {
	return 500
}

// StatusNotImplemented represents http status code 501.
type StatusNotImplemented struct{}

func (StatusNotImplemented) statusCode() int {
	return 501
}

// StatusBadGateway represents http status code 502.
type StatusBadGateway struct{}

func (StatusBadGateway) statusCode() int {
	return 502
}

// StatusServiceUnavailable represents http status code 503.
type StatusServiceUnavailable struct{}

func (StatusServiceUnavailable) statusCode() int {
	return 503
}

// StatusGatewayTimeout represents http status code 504.
type StatusGatewayTimeout struct{}

func (StatusGatewayTimeout) statusCode() int {
	return 504
}

// StatusHTTPVersionNotSupported represents http status code 505.
type StatusHTTPVersionNotSupported struct{}

func (StatusHTTPVersionNotSupported) statusCode() int {
	return 505
}

// StatusVariantAlsoNegotiates represents http status code 506.
type StatusVariantAlsoNegotiates struct{}

func (StatusVariantAlsoNegotiates) statusCode() int {
	return 506
}

// StatusInsufficientStorage represents http status code 507.
type StatusInsufficientStorage struct{}

func (StatusInsufficientStorage) statusCode() int {
	return 507
}

// StatusLoopDetected represents http status code 508.
type StatusLoopDetected struct{}

func (StatusLoopDetected) statusCode() int {
	return 508
}

// StatusNotExtended represents http status code 510.
type StatusNotExtended struct{}

func (StatusNotExtended) statusCode() int {
	return 510
}

// StatusNetworkAuthenticationRequired represents http status code 511.
type StatusNetworkAuthenticationRequired struct{}

func (StatusNetworkAuthenticationRequired) statusCode() int {
	return 511
}
