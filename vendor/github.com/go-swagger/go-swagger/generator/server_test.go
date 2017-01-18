package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func testGenOpts() (g GenOpts) {
	g.Target = "."
	g.APIPackage = "operations"
	g.ModelPackage = "models"
	g.ServerPackage = "restapi"
	g.ClientPackage = "client"
	g.Principal = ""
	g.DefaultScheme = "http"
	g.IncludeModel = true
	g.IncludeValidator = true
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.IncludeMain = false
	g.IncludeSupport = true
	g.ExcludeSpec = true
	g.TemplateDir = ""
	g.WithContext = false
	g.DumpData = false
	g.EnsureDefaults(false)
	return
}

func testAppGenertor(t testing.TB, specPath, name string) (*appGenerator, error) {
	specDoc, err := loads.Spec(specPath)
	if !assert.NoError(t, err) {
		return nil, err
	}
	analyzed := analysis.New(specDoc.Spec())

	models, err := gatherModels(specDoc, nil)
	if !assert.NoError(t, err) {
		return nil, err
	}

	operations := gatherOperations(analyzed, nil)
	if len(operations) == 0 {
		return nil, errors.New("no operations were selected")
	}

	opts := testGenOpts()
	opts.Spec = specPath
	apiPackage := opts.LanguageOpts.MangleName(swag.ToFileName(opts.APIPackage), "api")

	return &appGenerator{
		Name:            appNameOrDefault(specDoc, name, "swagger"),
		Receiver:        "o",
		SpecDoc:         specDoc,
		Analyzed:        analyzed,
		Models:          models,
		Operations:      operations,
		Target:          ".",
		DumpData:        opts.DumpData,
		Package:         apiPackage,
		APIPackage:      apiPackage,
		ModelsPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ModelPackage), "definitions"),
		ServerPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ServerPackage), "server"),
		ClientPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ClientPackage), "client"),
		Principal:       opts.Principal,
		DefaultScheme:   "http",
		DefaultProduces: runtime.JSONMime,
		DefaultConsumes: runtime.JSONMime,
		GenOpts:         &opts,
	}, nil
}

func TestServer_UrlEncoded(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	gen, err := testAppGenertor(t, "../fixtures/codegen/simplesearch.yml", "search")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("search_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assert.Regexp(t, "UrlformConsumer:\\s+runtime\\.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
			buf = bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverConfigureapi").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("configure_search_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, "api.UrlformConsumer = runtime.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_MultipartForm(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	gen, err := testAppGenertor(t, "../fixtures/codegen/shipyard.yml", "shipyard")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assert.Regexp(t, "MultipartformConsumer:\\s+runtime\\.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
			buf = bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverConfigureapi").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("configure_shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, "api.MultipartformConsumer = runtime.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_InvalidSpec(t *testing.T) {
	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/825/swagger.yml"
	opts.ValidateSpec = true
	assert.Error(t, GenerateServer("foo", nil, nil, &opts))
}
