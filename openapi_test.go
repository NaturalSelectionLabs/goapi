package goapi_test

import (
	"context"
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

type Three struct {
	goapi.StatusOK
	Data string `response:"direct"`
}

type Four struct {
	goapi.StatusOK
	Data goapi.DataBinary
}

func fnFour() Four {
	return Four{}
}

type Five struct {
	goapi.StatusOK
	Data goapi.DataBinary
}

func (Five) ContentType() string {
	return "image/png"
}

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	r.Use(middlewares.Identity)

	r.GET("/override", func(w http.ResponseWriter, r *http.Request) {})

	r.GET("/one", func(_ context.Context, p struct {
		goapi.InURL
		ID   string        `default:"123" description:"id" examples:"[\"456\"]"`
		Type *openapi.Code `description:"type code"`
	}, h struct {
		goapi.InHeader
		UA string
	}, b struct {
		Data string `json:"data"`
	}) Res {
		return Res01{}
	}).OpenAPI(func(doc *openapi.Operation) {
		doc.OperationID = "one"
		doc.Summary = "test"
		doc.Description = "test endpoint"
		doc.Tags = []string{"test"}
		doc.Security = []map[string][]string{{"auth": {"read"}}}
	})

	r.GET("/two/{id}", func(struct {
		goapi.InURL
		ID string
	}) Res03 {
		return Res03{}
	})

	r.GET("/three", func() Three {
		return Three{}
	})

	r.POST("/three", func() Three {
		return Three{}
	})

	r.GET("/four", fnFour)

	r.GET("/five", func() Five {
		return Five{}
	})

	doc := r.OpenAPI().JSON()

	// Ensure you have nodejs installed
	{
		g.E(os.WriteFile("tmp/openapi.json", []byte(doc), 0666))
		out, err := exec.Command("npx", strings.Split("rdme openapi:validate tmp/openapi.json", " ")...).CombinedOutput()
		g.Desc("%s", out).Nil(err)
	}

	g.Snapshot("openapi", g.JSON(doc))
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
