package openapi

//go:generate go run github.com/dmarkham/enumer@latest -type=Method -values -json
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

//go:generate go run github.com/dmarkham/enumer@latest -type=StatusCode -values -trimprefix=Status -json
type StatusCode int

const (
	Status200 StatusCode = 200
	Status201 StatusCode = 201
	Status202 StatusCode = 202
	Status204 StatusCode = 204
	Status301 StatusCode = 301
	Status302 StatusCode = 302
	Status303 StatusCode = 303
	Status304 StatusCode = 304
	Status305 StatusCode = 305
	Status307 StatusCode = 307
	Status400 StatusCode = 400
	Status401 StatusCode = 401
	Status402 StatusCode = 402
	Status403 StatusCode = 403
	Status404 StatusCode = 404
	Status405 StatusCode = 405
	Status406 StatusCode = 406
	Status407 StatusCode = 407
	Status408 StatusCode = 408
	Status409 StatusCode = 409
	Status410 StatusCode = 410
	Status411 StatusCode = 411
	Status412 StatusCode = 412
	Status413 StatusCode = 413
	Status414 StatusCode = 414
	Status415 StatusCode = 415
	Status416 StatusCode = 416
	Status417 StatusCode = 417
	Status418 StatusCode = 418
	Status421 StatusCode = 421
	Status422 StatusCode = 422
	Status423 StatusCode = 423
	Status424 StatusCode = 424
	Status425 StatusCode = 425
	Status426 StatusCode = 426
	Status428 StatusCode = 428
	Status429 StatusCode = 429
	Status431 StatusCode = 431
	Status451 StatusCode = 451
	Status500 StatusCode = 500
	Status501 StatusCode = 501
	Status502 StatusCode = 502
	Status503 StatusCode = 503
	Status504 StatusCode = 504
	Status505 StatusCode = 505
	Status506 StatusCode = 506
	Status507 StatusCode = 507
	Status508 StatusCode = 508
	Status510 StatusCode = 510
	Status511 StatusCode = 511
)
