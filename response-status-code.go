package goapi

type StatusContinue struct{}

func (StatusContinue) statusCode() int {
	return 100
}

type StatusSwitchingProtocols struct{}

func (StatusSwitchingProtocols) statusCode() int {
	return 101
}

type StatusProcessing struct{}

func (StatusProcessing) statusCode() int {
	return 102
}

type StatusEarlyHints struct{}

func (StatusEarlyHints) statusCode() int {
	return 103
}

type StatusOK struct{}

func (StatusOK) statusCode() int {
	return 200
}

type StatusCreated struct{}

func (StatusCreated) statusCode() int {
	return 201
}

type StatusAccepted struct{}

func (StatusAccepted) statusCode() int {
	return 202
}

type StatusNonAuthoritativeInfo struct{}

func (StatusNonAuthoritativeInfo) statusCode() int {
	return 203
}

type StatusNoContent struct{}

func (StatusNoContent) statusCode() int {
	return 204
}

type StatusResetContent struct{}

func (StatusResetContent) statusCode() int {
	return 205
}

type StatusPartialContent struct{}

func (StatusPartialContent) statusCode() int {
	return 206
}

type StatusMultiStatus struct{}

func (StatusMultiStatus) statusCode() int {
	return 207
}

type StatusAlreadyReported struct{}

func (StatusAlreadyReported) statusCode() int {
	return 208
}

type StatusIMUsed struct{}

func (StatusIMUsed) statusCode() int {
	return 226
}

type StatusMultipleChoices struct{}

func (StatusMultipleChoices) statusCode() int {
	return 300
}

type StatusMovedPermanently struct{}

func (StatusMovedPermanently) statusCode() int {
	return 301
}

type StatusFound struct{}

func (StatusFound) statusCode() int {
	return 302
}

type StatusSeeOther struct{}

func (StatusSeeOther) statusCode() int {
	return 303
}

type StatusNotModified struct{}

func (StatusNotModified) statusCode() int {
	return 304
}

type StatusUseProxy struct{}

func (StatusUseProxy) statusCode() int {
	return 305
}

type StatusTemporaryRedirect struct{}

func (StatusTemporaryRedirect) statusCode() int {
	return 307
}

type StatusPermanentRedirect struct{}

func (StatusPermanentRedirect) statusCode() int {
	return 308
}

type StatusBadRequest struct{}

func (StatusBadRequest) statusCode() int {
	return 400
}

type StatusUnauthorized struct{}

func (StatusUnauthorized) statusCode() int {
	return 401
}

type StatusPaymentRequired struct{}

func (StatusPaymentRequired) statusCode() int {
	return 402
}

type StatusForbidden struct{}

func (StatusForbidden) statusCode() int {
	return 403
}

type StatusNotFound struct{}

func (StatusNotFound) statusCode() int {
	return 404
}

type StatusMethodNotAllowed struct{}

func (StatusMethodNotAllowed) statusCode() int {
	return 405
}

type StatusNotAcceptable struct{}

func (StatusNotAcceptable) statusCode() int {
	return 406
}

type StatusProxyAuthRequired struct{}

func (StatusProxyAuthRequired) statusCode() int {
	return 407
}

type StatusRequestTimeout struct{}

func (StatusRequestTimeout) statusCode() int {
	return 408
}

type StatusConflict struct{}

func (StatusConflict) statusCode() int {
	return 409
}

type StatusGone struct{}

func (StatusGone) statusCode() int {
	return 410
}

type StatusLengthRequired struct{}

func (StatusLengthRequired) statusCode() int {
	return 411
}

type StatusPreconditionFailed struct{}

func (StatusPreconditionFailed) statusCode() int {
	return 412
}

type StatusRequestEntityTooLarge struct{}

func (StatusRequestEntityTooLarge) statusCode() int {
	return 413
}

type StatusRequestURITooLong struct{}

func (StatusRequestURITooLong) statusCode() int {
	return 414
}

type StatusUnsupportedMediaType struct{}

func (StatusUnsupportedMediaType) statusCode() int {
	return 415
}

type StatusRequestedRangeNotSatisfiable struct{}

func (StatusRequestedRangeNotSatisfiable) statusCode() int {
	return 416
}

type StatusExpectationFailed struct{}

func (StatusExpectationFailed) statusCode() int {
	return 417
}

type StatusTeapot struct{}

func (StatusTeapot) statusCode() int {
	return 418
}

type StatusMisdirectedRequest struct{}

func (StatusMisdirectedRequest) statusCode() int {
	return 421
}

type StatusUnprocessableEntity struct{}

func (StatusUnprocessableEntity) statusCode() int {
	return 422
}

type StatusLocked struct{}

func (StatusLocked) statusCode() int {
	return 423
}

type StatusFailedDependency struct{}

func (StatusFailedDependency) statusCode() int {
	return 424
}

type StatusTooEarly struct{}

func (StatusTooEarly) statusCode() int {
	return 425
}

type StatusUpgradeRequired struct{}

func (StatusUpgradeRequired) statusCode() int {
	return 426
}

type StatusPreconditionRequired struct{}

func (StatusPreconditionRequired) statusCode() int {
	return 428
}

type StatusTooManyRequests struct{}

func (StatusTooManyRequests) statusCode() int {
	return 429
}

type StatusRequestHeaderFieldsTooLarge struct{}

func (StatusRequestHeaderFieldsTooLarge) statusCode() int {
	return 431
}

type StatusUnavailableForLegalReasons struct{}

func (StatusUnavailableForLegalReasons) statusCode() int {
	return 451
}

type StatusInternalServerError struct{}

func (StatusInternalServerError) statusCode() int {
	return 500
}

type StatusNotImplemented struct{}

func (StatusNotImplemented) statusCode() int {
	return 501
}

type StatusBadGateway struct{}

func (StatusBadGateway) statusCode() int {
	return 502
}

type StatusServiceUnavailable struct{}

func (StatusServiceUnavailable) statusCode() int {
	return 503
}

type StatusGatewayTimeout struct{}

func (StatusGatewayTimeout) statusCode() int {
	return 504
}

type StatusHTTPVersionNotSupported struct{}

func (StatusHTTPVersionNotSupported) statusCode() int {
	return 505
}

type StatusVariantAlsoNegotiates struct{}

func (StatusVariantAlsoNegotiates) statusCode() int {
	return 506
}

type StatusInsufficientStorage struct{}

func (StatusInsufficientStorage) statusCode() int {
	return 507
}

type StatusLoopDetected struct{}

func (StatusLoopDetected) statusCode() int {
	return 508
}

type StatusNotExtended struct{}

func (StatusNotExtended) statusCode() int {
	return 510
}

type StatusNetworkAuthenticationRequired struct{}

func (StatusNetworkAuthenticationRequired) statusCode() int {
	return 511
}
