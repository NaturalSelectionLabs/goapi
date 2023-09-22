package goapi_test

import (
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/naturalselectionlabs/vary"
	"github.com/ysmood/got"
)

type Res interface {
	goapi.Response
}

var _ = goapi.Interface(new(Res), Res01{}, Res02{})

type Res01 struct {
	goapi.StatusOK

	Data string

	Header struct {
		SetCookie string
	}
}

type Res02 struct {
	goapi.StatusForbidden
	Error openapi.Error
}

func (Res02) Description() string {
	return "returns 403"
}

type Res03 struct {
	goapi.StatusOK

	Data string
	Meta string
}

type One struct {
}

func (One) OpenAPI(doc openapi.Operation) openapi.Operation {
	doc.Summary = "test"
	doc.Description = "test endpoint"
	doc.Tags = []string{"test"}
	doc.Security = []map[string][]string{{"auth": {"read"}}}

	return doc
}

func (One) Handle(p struct {
	goapi.InURL
	ID   string `default:"\"123\"" description:"id" example:"\"456\""`
	Type *openapi.Code
}, h struct {
	goapi.InHeader
	UA string
}, b struct {
	goapi.InBody
	Data string `json:"data"`
}) Res {
	return Res01{}
}

type Three struct {
	goapi.StatusOK
	Data string `response:"direct"`
}

type Four struct {
	goapi.StatusOK
	Data goapi.DataBinary
}

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	r.Use(middlewares.Identity)

	r.GET("/override", func(w http.ResponseWriter, r *http.Request) {})

	r.GET("/one", One{})

	r.GET("/two/{id}", func(struct {
		goapi.InURL
		ID string
	}) Res03 {
		return Res03{}
	})

	r.GET("/three", func() Three {
		return Three{}
	})

	r.GET("/four", func() Four {
		return Four{}
	})

	doc := r.OpenAPI().JSON()

	// Ensure you have nodejs installed
	{
		g.E(os.WriteFile("tmp/openapi.json", []byte(doc), 0666))
		out, err := exec.Command("npx", strings.Split("rdme openapi:validate tmp/openapi.json", " ")...).CombinedOutput()
		g.Desc("%s", out).Nil(err)
	}

	//nolint: lll
	g.Eq(g.JSON(doc), map[string]interface{}{
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"Code": map[string]interface{}{
					"description": "github.com/NaturalSelectionLabs/goapi/lib/openapi.Code",
					"enum": []interface{}{
						"internal_error",
						"invalid_param",
						"not_found",
					},
					"title": "Code",
				},
				"CommonError": map[string]interface{}{
					"additionalProperties": false,
					"description":          "github.com/NaturalSelectionLabs/goapi/lib/openapi.CommonError[github.com/NaturalSelectionLabs/goapi/lib/openapi.Code]",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{
							"$ref": "#/components/schemas/Code",
						},
						"details": map[string]interface{}{
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/CommonError",
							},
							"type": "array",
						},
						"innererror": map[string]interface{}{},
						"message": map[string]interface{}{
							"type": "string",
						},
						"target": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []interface{}{
						"code",
					},
					"title": "CommonError[github.com/NaturalSelectionLabs/goapi/lib/openapi.Code]",
					"type":  "object",
				},
				"Error": map[string]interface{}{
					"additionalProperties": false,
					"description":          "github.com/NaturalSelectionLabs/goapi/lib/openapi.Error",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{
							"$ref": "#/components/schemas/Code",
						},
						"details": map[string]interface{}{
							"items": map[string]interface{}{
								"$ref": "#/components/schemas/CommonError",
							},
							"type": "array",
						},
						"innererror": map[string]interface{}{},
						"message": map[string]interface{}{
							"type": "string",
						},
						"target": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []interface{}{
						"code",
					},
					"title": "Error",
					"type":  "object",
				},
			},
		},
		"info": map[string]interface{}{
			"title":   "",
			"version": "",
		},
		"openapi": "3.1.0",
		"paths": map[string]interface{}{
			"/four": map[string]interface{}{
				"get": map[string]interface{}{
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"content": map[string]interface{}{
								"application/octet-stream": map[string]interface{}{
									"schema": map[string]interface{}{
										"format": "binary",
										"type":   "string",
									},
								},
							},
							"description": "OK",
						},
					},
				},
			},
			"/one": map[string]interface{}{
				"get": map[string]interface{}{
					"description": "test endpoint",
					"operationId": "one",
					"parameters": []interface{}{
						map[string]interface{}{
							"description": "id",
							"in":          "query",
							"name":        "id",
							"schema": map[string]interface{}{
								"default": "123",
								"example": "456",
								"type":    "string",
							},
						},
						map[string]interface{}{
							"in":   "query",
							"name": "type",
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Code",
							},
						},
						map[string]interface{}{
							"in":       "header",
							"name":     "ua",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"requestBody": map[string]interface{}{
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"additionalProperties": false,
									"properties": map[string]interface{}{
										"data": map[string]interface{}{
											"type": "string",
										},
									},
									"required": []interface{}{
										"data",
									},
									"type": "object",
								},
							},
						},
						"required": true,
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"additionalProperties": false,
										"properties": map[string]interface{}{
											"data": map[string]interface{}{
												"type": "string",
											},
										},
										"required": []interface{}{
											"data",
										},
										"type": "object",
									},
								},
							},
							"description": "OK",
							"headers": map[string]interface{}{
								"set-cookie": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
						"403": map[string]interface{}{
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"additionalProperties": false,
										"properties": map[string]interface{}{
											"error": map[string]interface{}{
												"$ref": "#/components/schemas/Error",
											},
										},
										"required": []interface{}{
											"error",
										},
										"type": "object",
									},
								},
							},
							"description": "returns 403",
						},
					},
					"security": []interface{}{
						map[string]interface{}{
							"auth": []interface{}{
								"read",
							},
						},
					},
					"summary": "test",
					"tags": []interface{}{
						"test",
					},
				},
			},
			"/three": map[string]interface{}{
				"get": map[string]interface{}{
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "string",
									},
								},
							},
							"description": "OK",
						},
					},
				},
			},
			"/two/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"parameters": []interface{}{
						map[string]interface{}{
							"in":       "path",
							"name":     "id",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"additionalProperties": false,
										"properties": map[string]interface{}{
											"data": map[string]interface{}{
												"type": "string",
											},
											"meta": map[string]interface{}{
												"type": "string",
											},
										},
										"required": []interface{}{
											"data",
											"meta",
										},
										"type": "object",
									},
								},
							},
							"description": "OK",
						},
					},
				},
			},
		},
	})
}

func TestAddInterfaces(t *testing.T) {
	g := got.T(t)

	set := vary.NewInterfaces()

	type AddInterfaces interface{}

	set.New(new(AddInterfaces))

	goapi.AddInterfaces(set)

	g.Eq(goapi.Interfaces[vary.ID(reflect.TypeOf(new(AddInterfaces)).Elem())].ID(),
		"github.com/NaturalSelectionLabs/goapi_test.AddInterfaces")
}
