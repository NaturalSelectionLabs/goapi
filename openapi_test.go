package goapi_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
	"github.com/ysmood/vary"
)

type resGroup interface {
	goapi.Response
}

var iRes = vary.New(new(resGroup))

type res01 struct {
	goapi.Status200
	ID string `description:"response id"`
}

var _ = iRes.Add(res01{})

type res02 struct {
	goapi.Status403
	Error goapi.Error
}

func (res02) Description() string {
	return "returns 403"
}

var _ = iRes.Add(res02{})

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	r.GET("/test", func(p struct {
		goapi.InURL
		ID string
	}, h struct {
		goapi.InHeader
		UA string
	}) resGroup {
		return res01{}
	})

	g.Eq(g.JSON(r.OpenAPI(nil).JSON()), map[string]interface{} /* len=4 */ {
		"components": map[string]interface{} /* len=2 */ {
			"schemas": map[string]interface{} /* len=5 */ {
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
				"Status200": map[string]interface{} /* len=4 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.Status200`, /* len=47 */
					"title":                              "Status200",
					"type":                               "object",
				},
				"Status403": map[string]interface{} /* len=4 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.Status403`, /* len=47 */
					"title":                              "Status403",
					"type":                               "object",
				},
				"res01": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi_test.res01`, /* len=48 */
					"properties": map[string]interface{}{
						"ID": map[string]interface{} /* len=2 */ {
							"description": "response id",
							"type":        "string",
						},
					},
					"required": []interface{} /* len=1 cap=1 */ {
						"ID",
					},
					"title": "res01",
					"type":  "object",
				},
				"res02": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi_test.res02`, /* len=48 */
					"properties": map[string]interface{}{
						"Error": map[string]interface{}{
							"$ref": `#/components/schemas/Error`, /* len=26 */
						},
					},
					"required": []interface{} /* len=1 cap=1 */ {
						"Error",
					},
					"title": "res02",
					"type":  "object",
				},
			},
			"securitySchemes": nil,
		},
		"info": map[string]interface{} /* len=2 */ {
			"title":   "",
			"version": "",
		},
		"openapi": "3.1.0",
		"paths": map[string]interface{}{
			"/test": map[string]interface{}{
				"GET": map[string]interface{} /* len=2 */ {
					"parameters": []interface{} /* len=2 cap=2 */ {
						map[string]interface{} /* len=4 */ {
							"in":       "query",
							"name":     "id",
							"required": true,
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
					"responses": map[string]interface{} /* len=2 */ {
						"200": map[string]interface{}{
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": `#/components/schemas/res01`, /* len=26 */
									},
								},
							},
						},
						"403": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": `#/components/schemas/res02`, /* len=26 */
									},
								},
							},
							"description": "returns 403",
						},
					},
				},
			},
		},
	})

	g.Eq(1, 1)
}
