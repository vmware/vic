package generator

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// func init() {
// 	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
// }

func TestBuildDiscriminatorMap(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		di := discriminatorInfo(analysis.New(specDoc.Spec()))
		assert.Len(t, di.Discriminators, 1)
		assert.Len(t, di.Discriminators["#/definitions/Pet"].Children, 2)
		assert.Len(t, di.Discriminated, 2)
	}
}

func opts() *GenOpts {
	var opts GenOpts
	opts.IncludeValidator = true
	opts.IncludeModel = true
	if err := opts.EnsureDefaults(false); err != nil {
		panic(err)
	}
	return &opts
}

func TestGenerateModel_DiscriminatorSlices(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Kennel"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) {
			assert.True(t, genModel.HasBaseType)
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("has_discriminator.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "type Kennel struct {", res)
					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "Pets []Pet `json:\"pets\"`", res)
					assertInCode(t, "if err := m.Pets[i].Validate(formats); err != nil {", res)
					assertInCode(t, "m.validatePet", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenerateModel_Discriminators(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions

		for _, k := range []string{"cat", "Dog"} {
			schema := definitions[k]
			opts := opts()
			genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
			if assert.NoError(t, err) {
				assert.True(t, genModel.IsComplexObject)
				assert.Equal(t, "petType", genModel.DiscriminatorField)
				assert.Equal(t, k, genModel.DiscriminatorValue)
				buf := bytes.NewBuffer(nil)
				err := templates.MustGet("model").Execute(buf, genModel)
				if assert.NoError(t, err) {
					b, err := opts.LanguageOpts.FormatContent("discriminated.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(b)
						assertNotInCode(t, "m.Pet.Validate", res)
						assertInCode(t, "if err := m.validateName(formats); err != nil {", res)
						if k == "Dog" {
							assertInCode(t, "func (m *Dog) validatePackSize(formats strfmt.Registry) error {", res)
							assertInCode(t, "if err := m.validatePackSize(formats); err != nil {", res)
							assertInCode(t, "data.PackSize = m.PackSize", res)
							assertInCode(t, "validate.Required(\"packSize\", \"body\", m.PackSize)", res)
						} else {
							assertInCode(t, "func (m *Cat) validateHuntingSkill(formats strfmt.Registry) error {", res)
							assertInCode(t, "if err := m.validateHuntingSkill(formats); err != nil {", res)
							assertInCode(t, "if err := m.validateHuntingSkillEnum(\"huntingSkill\", \"body\", *m.HuntingSkill); err != nil {", res)
							assertInCode(t, "data.HuntingSkill = m.HuntingSkill", res)
						}
						assertInCode(t, "Name *string `json:\"name\"`", res)
						assertInCode(t, "PetType string `json:\"petType\"`", res)

						assertInCode(t, "data.Name = m.nameField", res)
						assertInCode(t, "data.PetType = \""+k+"\"", res)

						kk := swag.ToGoName(k)
						assertInCode(t, "func (m *"+kk+") Name() *string", res)
						assertInCode(t, "func (m *"+kk+") SetName(val *string)", res)
						assertInCode(t, "func (m *"+kk+") PetType() string", res)
						assertInCode(t, "func (m *"+kk+") SetPetType(val string)", res)
						assertInCode(t, "validate.Required(\"name\", \"body\", m.Name())", res)
					}
				}
			}
		}

		k := "Pet"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "petType", genModel.DiscriminatorField)
			assert.Len(t, genModel.Discriminates, 2)
			assert.Len(t, genModel.ExtraSchemas, 0)
			assert.Equal(t, "Cat", genModel.Discriminates["cat"])
			assert.Equal(t, "Dog", genModel.Discriminates["Dog"])
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("with_discriminator.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "type Pet interface {", res)
					assertInCode(t, "runtime.Validatable", res)
					assertInCode(t, "Name() *string", res)
					assertInCode(t, "SetName(*string)", res)
					assertInCode(t, "PetType() string", res)
					assertInCode(t, "SetPetType(string)", res)
					assertInCode(t, "UnmarshalPet(reader io.Reader, consumer runtime.Consumer) (Pet, error)", res)
					assertInCode(t, "PetType string `json:\"petType\"`", res)
					assertInCode(t, "validate.RequiredString(\"petType\"", res)
					assertInCode(t, "switch getType.PetType {", res)
					assertInCode(t, "var result Cat", res)
					assertInCode(t, "var result Dog", res)
				}
			}
		}

	}
}

func TestGenerateModel_UsesDiscriminator(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "WithPet"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.True(t, genModel.HasBaseType) {

			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("has_discriminator.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "type WithPet struct {", res)
					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "Pet Pet `json:\"-\"`", res)
					assertInCode(t, "if err := m.Pet.Validate(formats); err != nil {", res)
					assertInCode(t, "m.validatePet", res)
				}
			}
		}
	}
}

func TestGenerateClient_OKResponseWithDiscriminator(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		method, path, op, ok := analysis.New(specDoc.Spec()).OperationForName("modelOp")
		if assert.True(t, ok) {
			bldr := codeGenOpBuilder{
				Name:          "modelOp",
				Method:        method,
				Path:          path,
				APIPackage:    "restapi",
				ModelsPackage: "models",
				Principal:     "",
				Target:        ".",
				Doc:           specDoc,
				Analyzed:      analysis.New(specDoc.Spec()),
				Operation:     *op,
				Authed:        false,
				DefaultScheme: "http",
				ExtraSchemas:  make(map[string]GenSchema),
			}
			genOp, err := bldr.MakeOperation()
			if assert.NoError(t, err) {
				assert.True(t, genOp.Responses[0].Schema.IsBaseType)
				var buf bytes.Buffer
				err := templates.MustGet("clientResponse").Execute(&buf, genOp)
				if assert.NoError(t, err) {
					res := buf.String()
					assertInCode(t, "Payload models.Pet", res)
					assertInCode(t, "o.Payload = payload", res)
					assertInCode(t, "payload, err := models.UnmarshalPet(response.Body(), consumer)", res)
				}
			}
		}
	}
}

func TestGenerateServer_Parameters(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.discriminators.yml")
	if assert.NoError(t, err) {
		method, path, op, ok := analysis.New(specDoc.Spec()).OperationForName("modelOp")
		if assert.True(t, ok) {
			bldr := codeGenOpBuilder{
				Name:          "modelOp",
				Method:        method,
				Path:          path,
				APIPackage:    "restapi",
				ModelsPackage: "models",
				Principal:     "",
				Target:        ".",
				Doc:           specDoc,
				Analyzed:      analysis.New(specDoc.Spec()),
				Operation:     *op,
				Authed:        false,
				DefaultScheme: "http",
				ExtraSchemas:  make(map[string]GenSchema),
			}
			genOp, err := bldr.MakeOperation()
			if assert.NoError(t, err) {
				assert.True(t, genOp.Responses[0].Schema.IsBaseType)
				var buf bytes.Buffer
				err := templates.MustGet("serverParameter").Execute(&buf, genOp)
				if assert.NoError(t, err) {
					res := buf.String()
					assertInCode(t, "Pet models.Pet", res)
					assertInCode(t, "body, err := models.UnmarshalPet(r.Body, route.Consumer)", res)
					assertInCode(t, "o.Pet = body", res)
				}
			}
		}
	}
}

func TestGenerateModel_Discriminator_Billforward(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/billforward.discriminators.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "FlatPricingComponent"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.True(t, genModel.IsSubType) {

			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("has_discriminator.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					//assertInCode(t, "err", res)
					assertNotInCode(t, "for i := 0; i < len(m.PriceExplanation()); i++", res)
					// assertInCode(t, "swag.IsZero(m.priceExplanationField[i])", res)
				}
			}
		}
	}
}

func TestGenerateModel_Bitbucket_Repository(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/bitbucket.json")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "repository"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsNullable)
			for _, gm := range genModel.AllOf {
				for _, p := range gm.Properties {
					if p.Name == "parent" {
						assert.True(t, p.IsNullable)
					}
				}
			}
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("repository.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "Parent *Repository", res)
					assertNotInCode(t, "Parent Repository", res)
				}
			}
		}
	}
}

func TestGenerateModel_Bitbucket_WebhookSubscription(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/bitbucket.json")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "webhook_subscription"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("webhook_subscription.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "result.subjectField", res)
					assertInCode(t, "Subject: m.subjectField", res)
				}
			}
		}
	}
}

func TestGenerateModel_Issue319(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/319/swagger.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Container"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.Equal(t, "map[string]Base", genModel.Properties[0].GoType) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("ifacedmap.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "MapNoWorky map[string]Base", res)
				}
			}
		}
	}
}

func TestGenerateModel_Issue541(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/541/swagger.json")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Lion"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.AllOf) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("lion.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					// fmt.Println(res)
					assertInCode(t, "Cat", res)
					assertInCode(t, "m.Cat.Validate(formats)", res)
				}
			}
		}
	}
}

func TestGenerateModel_Issue436(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/436/swagger.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Image"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.AllOf) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("lion.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					// fmt.Println(res)
					assertInCode(t, "Links", res)
					assertInCode(t, "m.Links.Validate(formats)", res)
					assertInCode(t, "Created *strfmt.DateTime `json:\"created\"`", res)
					assertInCode(t, "ImageID *string `json:\"imageId\"`", res)
					assertInCode(t, "Size *int64 `json:\"size\"`", res)
				}
			}
		}
	}
}

func TestGenerateModel_Issue740(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/740/swagger.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Bar"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.AllOf) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("bar.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					assertInCode(t, "Foo", res)
					assertInCode(t, "m.Foo.Validate(formats)", res)
				}
			}
		}
	}
}

func TestGenerateModel_Issue743(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/743/swagger.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Awol"
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.AllOf) {
			buf := bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			if assert.NoError(t, err) {
				b, err := opts.LanguageOpts.FormatContent("awol.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(b)
					// fmt.Println(res)
					assertInCode(t, "Foo", res)
					assertInCode(t, "Bar", res)
					assertInCode(t, "m.Foo.Validate(formats)", res)
					assertInCode(t, "m.Bar.Validate(formats)", res)
					assertInCode(t, "swag.WriteJSON(m.Foo)", res)
					assertInCode(t, "swag.WriteJSON(m.Bar)", res)
					assertInCode(t, "swag.ReadJSON(raw, &aO0)", res)
					assertInCode(t, "swag.ReadJSON(raw, &aO1)", res)
					assertInCode(t, "m.Foo = aO0", res)
					assertInCode(t, "m.Bar = aO1", res)

				}
			}
		}
	}
}
