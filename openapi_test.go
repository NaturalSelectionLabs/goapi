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
	g.Eq(g.JSON(doc), map[string]interface{} /* len=4 */ {
		"components": map[string]interface{}{
			"schemas": map[string]interface{} /* len=3 */ {
				"Code": map[string]interface{} /* len=3 */ {
					"description": `github.com/NaturalSelectionLabs/goapi/lib/openapi.Code`, /* len=54 */
					"enum": []interface{} /* len=3 cap=4 */ {
						"internal_error",
						"invalid_param",
						"not_found",
					},
					"title": "Code",
				},
				"CommonError": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi/lib/openapi.CommonError[github.com/NaturalSelectionLabs/goapi/lib/openapi.Code]`, /* len=117 */
					"properties": map[string]interface{} /* len=5 */ {
						"code": map[string]interface{}{
							"$ref": `#/components/schemas/Code`, /* len=25 */
						},
						"details": map[string]interface{} /* len=2 */ {
							"items": map[string]interface{}{
								"$ref": `#/components/schemas/CommonError`, /* len=32 */
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
					"required": []interface{} /* len=1 cap=1 */ {
						"code",
					},
					"title": `CommonError[github.com/NaturalSelectionLabs/goapi/lib/openapi.Code]`, /* len=67 */
					"type":  "object",
				},
				"Error": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi/lib/openapi.Error`, /* len=55 */
					"properties": map[string]interface{} /* len=5 */ {
						"code": map[string]interface{}{
							"$ref": `#/components/schemas/Code`, /* len=25 */
						},
						"details": map[string]interface{} /* len=2 */ {
							"items": map[string]interface{}{
								"$ref": `#/components/schemas/CommonError`, /* len=32 */
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
					"required": []interface{} /* len=1 cap=1 */ {
						"code",
					},
					"title": "Error",
					"type":  "object",
				},
			},
		},
		"info": map[string]interface{} /* len=2 */ {
			"title":   "",
			"version": "",
		},
		"openapi": "3.1.0",
		"paths": map[string]interface{} /* len=4 */ {
			"/four": map[string]interface{}{
				"get": map[string]interface{}{
					"responses": map[string]interface{}{
						"200": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								`application/octet-stream` /* len=24 */ : map[string]interface{}{
									"schema": map[string]interface{} /* len=2 */ {
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
				"get": map[string]interface{} /* len=8 */ {
					"description": "test endpoint",
					"operationId": "one",
					"parameters": []interface{} /* len=3 cap=4 */ {
						map[string]interface{} /* len=4 */ {
							"description": "id",
							"in":          "query",
							"name":        "id",
							"schema": map[string]interface{} /* len=3 */ {
								"default": "123",
								"example": "456",
								"type":    "string",
							},
						},
						map[string]interface{} /* len=3 */ {
							"in":   "query",
							"name": "type",
							"schema": map[string]interface{}{
								"$ref": `#/components/schemas/Code`, /* len=25 */
							},
						},
						map[string]interface{} /* len=4 */ {
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
							"application/json" /* len=16 */ : map[string]interface{}{
								"schema": map[string]interface{} /* len=4 */ {
									`additionalProperties` /* len=20 */ : false,
									"properties": map[string]interface{}{
										"data": map[string]interface{}{
											"type": "string",
										},
									},
									"required": []interface{} /* len=1 cap=1 */ {
										"data",
									},
									"type": "object",
								},
							},
						},
					},
					"responses": map[string]interface{} /* len=2 */ {
						"200": map[string]interface{} /* len=3 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{} /* len=4 */ {
										`additionalProperties` /* len=20 */ : false,
										"properties": map[string]interface{}{
											"data": map[string]interface{}{
												"type": "string",
											},
										},
										"required": []interface{} /* len=1 cap=1 */ {
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
						"403": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{} /* len=4 */ {
										`additionalProperties` /* len=20 */ : false,
										"properties": map[string]interface{}{
											"error": map[string]interface{}{
												"$ref": `#/components/schemas/Error`, /* len=26 */
											},
										},
										"required": []interface{} /* len=1 cap=1 */ {
											"error",
										},
										"type": "object",
									},
								},
							},
							"description": "returns 403",
						},
					},
					"security": []interface{} /* len=1 cap=1 */ {
						map[string]interface{}{
							"auth": []interface{} /* len=1 cap=1 */ {
								"read",
							},
						},
					},
					"summary": "test",
					"tags": []interface{} /* len=1 cap=1 */ {
						"test",
					},
				},
			},
			"/three": map[string]interface{}{
				"get": map[string]interface{}{
					"responses": map[string]interface{}{
						"200": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
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
				"get": map[string]interface{} /* len=2 */ {
					"parameters": []interface{} /* len=1 cap=1 */ {
						map[string]interface{} /* len=4 */ {
							"in":       "path",
							"name":     "id",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{} /* len=4 */ {
										`additionalProperties` /* len=20 */ : false,
										"properties": map[string]interface{} /* len=2 */ {
											"data": map[string]interface{}{
												"type": "string",
											},
											"meta": map[string]interface{}{
												"type": "string",
											},
										},
										"required": []interface{} /* len=2 cap=2 */ {
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
