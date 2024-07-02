package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertExternalLinksMapToInputType(t *testing.T) {
	externalLinksMap := map[string]attr.Value{
		"example-label": types.StringValue("https://example.com"),
	}

	expected := `{"links":[{"label":"example-label","url":"https://example.com"}]}`

	result, err := convertExternalLinksMapToInputType(externalLinksMap)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, *result)
}

func TestConvertExternalLinksListToMap(t *testing.T) {
	externalLinks := ExternalLinks{
		Links: []ExternalLink{
			{Label: "example-label", URL: "https://example.com"},
		},
	}

	expected := map[string]attr.Value{
		"example-label": types.StringValue("https://example.com"),
	}

	result, diags := convertExternalLinksListToMap(externalLinks)
	assert.False(t, diags.HasError())
	assert.Equal(t, expected, result.Elements())
}

func TestGenerateCompositeID(t *testing.T) {
	entityName := "example-entity"
	queueName := "example-queue"

	expected := "example-entity:example-queue"
	result := generateCompositeID(entityName, queueName)
	assert.Equal(t, expected, result)
}

func TestParseCompositeID(t *testing.T) {
	compositeID := "example-entity:example-queue"

	entityName, queueName, err := parseCompositeID(compositeID)
	assert.NoError(t, err)
	assert.Equal(t, "example-entity", entityName)
	assert.Equal(t, "example-queue", queueName)
}

func TestParseCompositeID_InvalidFormat(t *testing.T) {
	compositeID := "invalid-composite-id"

	_, _, err := parseCompositeID(compositeID)
	assert.Error(t, err)
}

func TestTemplateVarsWithNamesListToMap(t *testing.T) {
	schema := TemplateVariableSchema{
		Type:    "string",
		Default: "default1",
		Enum:    []string{"default1", "default2"},
	}
	schemaString := `{"type":"string", "default":"default1", "enum":["default1","default2"]}`
	description := "desc1"
	tvList := []TemplateVariableWithName{
		{
			Name:        "var1",
			Description: &description,
			Schema:      schemaString,
		},
	}

	expected := map[string]TemplateVariable{
		"var1": {
			Description: &description,
			Schema:      schema,
		},
	}

	result, err := templateVarsWithNamesListToMap(tvList)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestInjectResourceArgsAndResourceFieldsEmpty(t *testing.T) {
	resourceConfig := ""
	resourceType := "kubernetes"

	expected := `{"resource_args":{"kubernetes":{}}}`

	result, err := injectResourceArgsAndResourceFields(resourceConfig, resourceType)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, result)
}

func TestInjectResourceArgsAndResourceFields(t *testing.T) {
	resourceConfig := `{"apiVersion":"batch/v1","kind":"Job","metadata":{"name":"example-job"}}`
	resourceType := "kubernetes"

	expected := `{"resource_args":{"kubernetes":{"apiVersion":"batch/v1","kind":"Job","metadata":{"name":"example-job"}}}}`

	result, err := injectResourceArgsAndResourceFields(resourceConfig, resourceType)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, result)
}

func TestInjectResourceArgsAndResourceFields_AlreadyPresent(t *testing.T) {
	resourceConfig := `{"resource_args":{"kubernetes":{"apiVersion":"batch/v1","kind":"Job","metadata":{"name":"example-job"}}}}`
	resourceType := "kubernetes"

	result, err := injectResourceArgsAndResourceFields(resourceConfig, resourceType)
	assert.NoError(t, err)
	assert.JSONEq(t, resourceConfig, result)
}
