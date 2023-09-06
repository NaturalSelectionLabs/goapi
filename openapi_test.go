package goapi_test

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
	"github.com/ysmood/vary"
)

type Res interface {
	goapi.Response
}

var iRes = vary.New(new(Res))

type Res01 struct {
	goapi.StatusOK

	Data string

	Header struct {
		SetCookie string
	}
}

var _ = iRes.Add(Res01{})

type Res02 struct {
	goapi.StatusForbidden
	Error goapi.Error
}

func (Res02) Description() string {
	return "returns 403"
}

var _ = iRes.Add(Res02{})

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

	r.Use(goapi.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}))

	r.GET("/one", func(p struct {
		goapi.InURL
		ID string `description:"id"`
	}, h struct {
		goapi.InHeader
		UA string
	}, b struct {
		goapi.InBody
		Data string `json:"data"`
	}) Res {
		return Res01{}
	}, r.Meta(goapi.OperationMeta{
		Summary:     "test",
		Description: "test endpoint",
		OperationID: "test",
		Tags:        []string{"test"},
	}))

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
					"description":                        `github.com/NaturalSelectionLabs/goapi.Error`, /* len=43 */
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
				"get": map[string]interface{} /* len=7 */ {
					"description": "test endpoint",
					"operationId": "test",
					"parameters": []interface{} /* len=2 cap=2 */ {
						map[string]interface{} /* len=5 */ {
							"description": "id",
							"in":          "query",
							"name":        "id",
							"required":    true,
							"schema": map[string]interface{}{
								"type": "string",
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
							"description": "",
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
							"description": "",
						},
					},
				},
			},
		},
	})
}
