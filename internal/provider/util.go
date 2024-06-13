package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/machinebox/graphql"
)

func convertExternalLinksMapToInputType(externalLinksMap map[string]attr.Value) (*string, error) {
	var result ExternalLinks

	for label, v := range externalLinksMap {
		if v.IsNull() {
			continue
		}

		// Assert that the value is a string
		url, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("unexpected type for external link URL, expected types.String, got %T", v)
		}

		// Append the link to the result
		result.Links = append(result.Links, ExternalLink{
			Label: label,
			URL:   url.ValueString(),
		})
	}

	if len(result.Links) == 0 {
		return nil, nil
	}

	linksBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	linksString := string(linksBytes)
	return &linksString, nil
}

func convertExternalLinksListToMap(input ExternalLinks) (types.Map, diag.Diagnostics) {
	result := make(map[string]attr.Value)
	for _, link := range input.Links {
		result[link.Label] = types.StringValue(link.URL)
	}
	return types.MapValue(types.StringType, result)
}

// generateCompositeID generates a composite ID from entityName and queueName.
func generateCompositeID(entityName, queueName string) string {
	return fmt.Sprintf("%s:%s", entityName, queueName)
}

// parseCompositeID parses a composite ID into entityName and queueName.
func parseCompositeID(compositeID string) (string, string, error) {
	parts := strings.SplitN(compositeID, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid composite ID: %s", compositeID)
	}
	return parts[0], parts[1], nil
}

func readRunQueueHelper(entityName, queueName string, ctx context.Context, client GraphQLClientWithHeaders) (*RunQueue, error) {
	if entityName == "" || queueName == "" {
		return nil, fmt.Errorf("entity_name and name must be specified")
	}

	gqlReq := graphql.NewRequest(`
		query GetRunQueueByName($entityName:String!, $projectName: String!, $queueName: String!) {
			project(entityName: $entityName, name: $projectName) {
				runQueue(name: $queueName) {
					id
					name
					entityName
					defaultResourceConfig {
						id
						resource
						config
						templateVariables {
							name
							description
							schema
						}
					}
					prioritizationMode
					externalLinks
					createdAt
					updatedAt
				}
			}
		}
	`)
	gqlReq.Var("entityName", entityName)
	gqlReq.Var("queueName", queueName)
	gqlReq.Var("projectName", "model-registry")
	var result struct {
		Project struct {
			RunQueue *RunQueue `json:"runQueue"`
		} `json:"project"`
	}

	if err := client.Run(ctx, gqlReq, &result); err != nil {
		return nil, err
	}

	if result.Project.RunQueue == nil {
		return nil, fmt.Errorf("run queue not found")
	}

	return result.Project.RunQueue, nil
}

func templateVarsWithNamesListToMap(tvList []TemplateVariableWithName) map[string]TemplateVariable {
	result := make(map[string]TemplateVariable)
	for _, tv := range tvList {
		result[tv.Name] = TemplateVariable{
			Description: tv.Description,
			Schema:      tv.Schema,
		}
	}
	return result
}

func injectResourceArgsAndResourceFields(resourceConfig string, resourceType string) (string, error) {
	var resourceArgs map[string]interface{}
	if err := json.Unmarshal([]byte(resourceConfig), &resourceArgs); err != nil {
		return "", err
	}

	if _, ok := resourceArgs["resource_args"]; !ok {
		newResourceArgs := map[string]interface{}{
			"resource_args": map[string]interface{}{
				resourceType: resourceArgs,
			},
		}
		resourceBytes, err := json.Marshal(newResourceArgs)
		if err != nil {
			return "", err
		}
		return string(resourceBytes), nil
	}
	return resourceConfig, nil
}

func stripResourceArgsAndResourceFields(resourceConfig, resourceType string) (string, error) {
	var resourceArgs map[string]interface{}
	if err := json.Unmarshal([]byte(resourceConfig), &resourceArgs); err != nil {
		return "", err
	}

	if _, ok := resourceArgs["resource_args"]; ok {
		resourceArgs = resourceArgs["resource_args"].(map[string]interface{})
		if _, ok := resourceArgs[resourceType]; ok {
			resourceBytes, err := json.Marshal(resourceArgs[resourceType])
			if err != nil {
				return "", err
			}
			return string(resourceBytes), nil
		} else {
			return "", fmt.Errorf("resource type %s not found in resource_args", resourceType)
		}
	} else {
		return "", fmt.Errorf("resource_args not found in resource config")
	}
}

func hasResourceArgsAndResourceFields(resourceConfig, resourceType string) bool {
	var resourceArgs map[string]interface{}
	if err := json.Unmarshal([]byte(resourceConfig), &resourceArgs); err != nil {
		return false
	}

	if _, ok := resourceArgs["resource_args"]; ok {
		resourceArgs = resourceArgs["resource_args"].(map[string]interface{})
		if _, ok := resourceArgs[resourceType]; ok {
			return true
		}
	}
	return false
}
