package goapi_test

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
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

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	r.Use(middlewares.Identity)

	r.GET("/override", func(w http.ResponseWriter, r *http.Request) {})

	r.GET("/one", func(p struct {
		goapi.InURL
		ID string `default:"\"123\"" description:"id" example:"\"456\""`
	}, h struct {
		goapi.InHeader
		UA string
	}, b struct {
		goapi.InBody
		Data string `json:"data"`
	}) Res {
		return Res01{}
	},
		goapi.Summary("test"),
		goapi.Description("test endpoint"),
		goapi.OperationID("test"),
		goapi.Tags("test"),
		goapi.Security(map[string][]string{"auth": {"read"}}),
	)

	r.GET("/two/{id}", func(struct {
		goapi.InURL
		ID string
	}) Res03 {
		return Res03{}
	})

	doc := r.OpenAPI(nil).JSON()

	// Ensure you have nodejs installed
	{
		g.E(os.WriteFile("tmp/openapi.json", []byte(doc), 0666))
		out, err := exec.Command("npx", strings.Split("rdme openapi:validate tmp/openapi.json", " ")...).CombinedOutput()
		g.Desc("%s", out).E(err)
	}

	g.Eq(g.JSON(doc), map[string]interface{} /* len=4 */ {
		"components": map[string]interface{}{
			"schemas": map[string]interface{} /* len=2 */ {
				"Error": map[string]interface{} /* len=5 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi/lib/openapi.Error`, /* len=55 */
					"properties": map[string]interface{} /* len=5 */ {
						"code": map[string]interface{}{
							"type": "string",
						},
						"details": map[string]interface{} /* len=2 */ {
							"items": map[string]interface{}{
								"$ref": `#/components/schemas/Error`, /* len=26 */
							},
							"type": "array",
						},
						"innererror": map[string]interface{}{
							"type": "object",
						},
						"message": map[string]interface{}{
							"type": "string",
						},
						"target": map[string]interface{}{
							"type": "string",
						},
					},
					"title": "Error",
					"type":  "object",
				},
				"InBody": map[string]interface{} /* len=4 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.InBody`, /* len=44 */
					"title":                              "InBody",
					"type":                               "object",
				},
			},
		},
		"info": map[string]interface{} /* len=2 */ {
			"title":   "",
			"version": "",
		},
		"openapi": "3.1.0",
		"paths": map[string]interface{} /* len=2 */ {
			"/one": map[string]interface{}{
				"get": map[string]interface{} /* len=8 */ {
					"description": "test endpoint",
					"operationId": "test",
					"parameters": []interface{} /* len=2 cap=2 */ {
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
