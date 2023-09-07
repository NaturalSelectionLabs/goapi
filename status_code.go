package goapi

// StatusContinue 100.
type StatusContinue struct{}

func (StatusContinue) statusCode() int {
	return 100
}

// StatusSwitchingProtocols 101.
type StatusSwitchingProtocols struct{}

func (StatusSwitchingProtocols) statusCode() int {
	return 101
}

// StatusProcessing 102.
type StatusProcessing struct{}

func (StatusProcessing) statusCode() int {
	return 102
}

// StatusEarlyHints 103.
type StatusEarlyHints struct{}

func (StatusEarlyHints) statusCode() int {
	return 103
}

// StatusOK 200.
type StatusOK struct{}

func (StatusOK) statusCode() int {
	return 200
}

// StatusCreated 201.
type StatusCreated struct{}

func (StatusCreated) statusCode() int {
	return 201
}

// StatusAccepted 202.
type StatusAccepted struct{}

func (StatusAccepted) statusCode() int {
	return 202
}

// StatusNonAuthoritativeInfo 203.
type StatusNonAuthoritativeInfo struct{}

func (StatusNonAuthoritativeInfo) statusCode() int {
	return 203
}

// StatusNoContent 204.
type StatusNoContent struct{}

func (StatusNoContent) statusCode() int {
	return 204
}

// StatusResetContent 205.
type StatusResetContent struct{}

func (StatusResetContent) statusCode() int {
	return 205
}

// StatusPartialContent 206.
type StatusPartialContent struct{}

func (StatusPartialContent) statusCode() int {
	return 206
}

// StatusMultiStatus 207.
type StatusMultiStatus struct{}

func (StatusMultiStatus) statusCode() int {
	return 207
}

// StatusAlreadyReported 208.
type StatusAlreadyReported struct{}

func (StatusAlreadyReported) statusCode() int {
	return 208
}

// StatusIMUsed 226.
type StatusIMUsed struct{}

func (StatusIMUsed) statusCode() int {
	return 226
}

// StatusMultipleChoices 300.
type StatusMultipleChoices struct{}

func (StatusMultipleChoices) statusCode() int {
	return 300
}

// StatusMovedPermanently 301.
type StatusMovedPermanently struct{}

func (StatusMovedPermanently) statusCode() int {
	return 301
}

// StatusFound 302.
type StatusFound struct{}

func (StatusFound) statusCode() int {
	return 302
}

// StatusSeeOther 303.
type StatusSeeOther struct{}

func (StatusSeeOther) statusCode() int {
	return 303
}

// StatusNotModified 304.
type StatusNotModified struct{}

func (StatusNotModified) statusCode() int {
	return 304
}

// StatusUseProxy 305.
type StatusUseProxy struct{}

func (StatusUseProxy) statusCode() int {
	return 305
}

// StatusTemporaryRedirect 307.
type StatusTemporaryRedirect struct{}

func (StatusTemporaryRedirect) statusCode() int {
	return 307
}

// StatusPermanentRedirect 308.
type StatusPermanentRedirect struct{}

func (StatusPermanentRedirect) statusCode() int {
	return 308
}

// StatusBadRequest 400.
type StatusBadRequest struct{}

func (StatusBadRequest) statusCode() int {
	return 400
}

// StatusUnauthorized 401.
type StatusUnauthorized struct{}

func (StatusUnauthorized) statusCode() int {
	return 401
}

// StatusPaymentRequired 402.
type StatusPaymentRequired struct{}

func (StatusPaymentRequired) statusCode() int {
	return 402
}

// StatusForbidden 403.
type StatusForbidden struct{}

func (StatusForbidden) statusCode() int {
	return 403
}

// StatusNotFound 404.
type StatusNotFound struct{}

func (StatusNotFound) statusCode() int {
	return 404
}

// StatusMethodNotAllowed 405.
type StatusMethodNotAllowed struct{}

func (StatusMethodNotAllowed) statusCode() int {
	return 405
}

// StatusNotAcceptable 406.
type StatusNotAcceptable struct{}

func (StatusNotAcceptable) statusCode() int {
	return 406
}

// StatusProxyAuthRequired 407.
type StatusProxyAuthRequired struct{}

func (StatusProxyAuthRequired) statusCode() int {
	return 407
}

// StatusRequestTimeout 408.
type StatusRequestTimeout struct{}

func (StatusRequestTimeout) statusCode() int {
	return 408
}

// StatusConflict 409.
type StatusConflict struct{}

func (StatusConflict) statusCode() int {
	return 409
}

// StatusGone 410.
type StatusGone struct{}

func (StatusGone) statusCode() int {
	return 410
}

// StatusLengthRequired 411.
type StatusLengthRequired struct{}

func (StatusLengthRequired) statusCode() int {
	return 411
}

// StatusPreconditionFailed 412.
type StatusPreconditionFailed struct{}

func (StatusPreconditionFailed) statusCode() int {
	return 412
}

// StatusRequestEntityTooLarge 413.
type StatusRequestEntityTooLarge struct{}

func (StatusRequestEntityTooLarge) statusCode() int {
	return 413
}

// StatusRequestURITooLong 414.
type StatusRequestURITooLong struct{}

func (StatusRequestURITooLong) statusCode() int {
	return 414
}

// StatusUnsupportedMediaType 415.
type StatusUnsupportedMediaType struct{}

func (StatusUnsupportedMediaType) statusCode() int {
	return 415
}

// StatusRequestedRangeNotSatisfiable 416.
type StatusRequestedRangeNotSatisfiable struct{}

func (StatusRequestedRangeNotSatisfiable) statusCode() int {
	return 416
}

// StatusExpectationFailed 417.
type StatusExpectationFailed struct{}

func (StatusExpectationFailed) statusCode() int {
	return 417
}

// StatusTeapot 418.
type StatusTeapot struct{}

func (StatusTeapot) statusCode() int {
	return 418
}

// StatusMisdirectedRequest 421.
type StatusMisdirectedRequest struct{}

func (StatusMisdirectedRequest) statusCode() int {
	return 421
}

// StatusUnprocessableEntity 422.
type StatusUnprocessableEntity struct{}

func (StatusUnprocessableEntity) statusCode() int {
	return 422
}

// StatusLocked 423.
type StatusLocked struct{}

func (StatusLocked) statusCode() int {
	return 423
}

// StatusFailedDependency 424.
type StatusFailedDependency struct{}

func (StatusFailedDependency) statusCode() int {
	return 424
}

// StatusTooEarly 425.
type StatusTooEarly struct{}

func (StatusTooEarly) statusCode() int {
	return 425
}

// StatusUpgradeRequired 426.
type StatusUpgradeRequired struct{}

func (StatusUpgradeRequired) statusCode() int {
	return 426
}

// StatusPreconditionRequired 428.
type StatusPreconditionRequired struct{}

func (StatusPreconditionRequired) statusCode() int {
	return 428
}

// StatusTooManyRequests 429.
type StatusTooManyRequests struct{}

func (StatusTooManyRequests) statusCode() int {
	return 429
}

// StatusRequestHeaderFieldsTooLarge 431.
type StatusRequestHeaderFieldsTooLarge struct{}

func (StatusRequestHeaderFieldsTooLarge) statusCode() int {
	return 431
}

// StatusUnavailableForLegalReasons 451.
type StatusUnavailableForLegalReasons struct{}

func (StatusUnavailableForLegalReasons) statusCode() int {
	return 451
}

// StatusInternalServerError 500.
type StatusInternalServerError struct{}

func (StatusInternalServerError) statusCode() int {
	return 500
}

// StatusNotImplemented 501.
type StatusNotImplemented struct{}

func (StatusNotImplemented) statusCode() int {
	return 501
}

// StatusBadGateway 502.
type StatusBadGateway struct{}

func (StatusBadGateway) statusCode() int {
	return 502
}

// StatusServiceUnavailable 503.
type StatusServiceUnavailable struct{}

func (StatusServiceUnavailable) statusCode() int {
	return 503
}

// StatusGatewayTimeout 504.
type StatusGatewayTimeout struct{}

func (StatusGatewayTimeout) statusCode() int {
	return 504
}

// StatusHTTPVersionNotSupported 505.
type StatusHTTPVersionNotSupported struct{}

func (StatusHTTPVersionNotSupported) statusCode() int {
	return 505
}

// StatusVariantAlsoNegotiates 506.
type StatusVariantAlsoNegotiates struct{}

func (StatusVariantAlsoNegotiates) statusCode() int {
	return 506
}

// StatusInsufficientStorage 507.
type StatusInsufficientStorage struct{}

func (StatusInsufficientStorage) statusCode() int {
	return 507
}

// StatusLoopDetected 508.
type StatusLoopDetected struct{}

func (StatusLoopDetected) statusCode() int {
	return 508
}

// StatusNotExtended 510.
type StatusNotExtended struct{}

func (StatusNotExtended) statusCode() int {
	return 510
}

// StatusNetworkAuthenticationRequired 511.
type StatusNetworkAuthenticationRequired struct{}

func (StatusNetworkAuthenticationRequired) statusCode() int {
	return 511
}
