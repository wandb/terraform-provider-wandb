// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/machinebox/graphql"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RunQueueResource{}
var _ resource.ResourceWithConfigure = &RunQueueResource{}
var _ resource.ResourceWithImportState = &RunQueueResource{}

func NewRunQueueResource() resource.Resource {
	return &RunQueueResource{}
}

type RunQueueResource struct {
	client *GraphQLClientWithHeaders
}

type RunQueueResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	EntityName         types.String `tfsdk:"entity_name"`
	Resource           types.String `tfsdk:"resource"`
	ResourceConfig     types.String `tfsdk:"resource_config"`
	TemplateVariables  types.String `tfsdk:"template_variables"`
	PrioritizationMode types.String `tfsdk:"prioritization_mode"`
	ExternalLinks      types.Map    `tfsdk:"external_links"`
}

func (r *RunQueueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "wandb_run_queue"
}

func (r *RunQueueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RunQueue resource used with W&B Launch. See: https://docs.wandb.ai/guides/launch. See [here](https://github.com/wandb/terraform-provider-wandb/blob/main/examples/resources/run_queue/resource.tf) for an example",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the run queue. This is a composite ID of the entity name and the queue name, separated by a ':'",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the run queue. This is unique within the entity.",
			},
			"entity_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the entity that this run queue belongs to.",
			},
			"resource": schema.StringAttribute{
				Required:    true,
				Description: "The resource type for this queue, options include: 'local-container', 'kubernetes', 'vertex', 'sagemaker'",
			},
			"resource_config": schema.StringAttribute{
				Optional:    true,
				Description: "The configuration for the resource type. This is a JSON string that will be passed to the resource. For more information about the resource configuration see: https://docs.wandb.ai/guides/launch/setup-launch",
			},
			"template_variables": schema.StringAttribute{
				Optional:    true,
				Description: "The template variables for the resource configuration. This is a JSON string that will be passed to the resource. For more information about the template variables see: https://docs.wandb.ai/guides/launch/setup-queue-advanced#configure-queue-template",
			},
			"prioritization_mode": schema.StringAttribute{
				Optional:    true,
				Description: "The prioritization mode for the run queue. Options include: disabled and V0. V0 allows users to specify priority when launching items. Once a queue specifies V0, it can not be disabled.",
			},
			"external_links": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A map of external links for the run queue. Provided as a map with the key being the label, and the value being the URL.",
			},
		},
	}
}

func (r *RunQueueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*GraphQLClientWithHeaders)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RunQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RunQueueResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	externalLinks, err := convertExternalLinksMapToInputType(data.ExternalLinks.Elements())
	if err != nil {
		resp.Diagnostics.AddError("Error converting external links", err.Error())
		return
	}

	prioritizationMode := data.PrioritizationMode.ValueStringPointer()
	if prioritizationMode == nil {
		defaultPrioritizationMode := "V0"
		prioritizationMode = &defaultPrioritizationMode
	}

	normalizedTemplateVariables, err := normalizeTemplateVariables(data.TemplateVariables.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error normalizing template variables", err.Error())
		return
	}
	if normalizedTemplateVariables != nil {
		data.TemplateVariables = types.StringValue(*normalizedTemplateVariables)
	}

	// Inject resource args and fields into the resource config backend expects wrapped in these fields
	resourceConfig, err := injectResourceArgsAndResourceFields(data.ResourceConfig.ValueString(), data.Resource.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error injecting resource args and fields", err.Error())
		return
	}

	input := UpsertRunQueueInput{
		QueueName:          data.Name.ValueString(),
		EntityName:         data.EntityName.ValueString(),
		ProjectName:        "model-registry",
		ResourceType:       data.Resource.ValueString(),
		ResourceConfig:     resourceConfig,
		TemplateVariables:  data.TemplateVariables.ValueStringPointer(),
		PrioritizationMode: prioritizationMode,
		ExternalLinks:      externalLinks,
	}

	result, err := upsertRunQueue(ctx, input, r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating run queue",
			"Could not create run queue, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.UpsertRunQueue.Success {
		resp.Diagnostics.AddError(
			"Failed to create run queue",
			"The API did not confirm the creation of the run queue.",
		)
		return
	}

	if len(result.UpsertRunQueue.Errors) > 0 {
		resp.Diagnostics.AddWarning("Config schema validation errors", strings.Join(result.UpsertRunQueue.Errors, ", "))
	}

	id := generateCompositeID(data.EntityName.ValueString(), data.Name.ValueString())
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a run queue resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RunQueueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	entityName, queueName, err := parseCompositeID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing composite ID", err.Error())
		return
	}
	runQueue, err := readRunQueueHelper(entityName, queueName, ctx, *r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading run queue",
			"Could not read run queue, unexpected error: "+err.Error(),
		)
		return
	}

	// Map API response to the Terraform state model
	data.Name = types.StringValue(runQueue.Name)
	data.EntityName = types.StringValue(runQueue.EntityName)
	data.Resource = types.StringValue(runQueue.DefaultResourceConfig.Resource)
	data.PrioritizationMode = types.StringValue(runQueue.PrioritizationMode)
	byteConfig, err := json.Marshal(runQueue.DefaultResourceConfig.Config)
	if err != nil {
		resp.Diagnostics.AddError("Error marshalling resource config", err.Error())
		return
	}

	// Handle issue where empty configs are saved as a resource_args.<resource> map
	emptyConfigString := fmt.Sprintf(`{"resource_args":{"%s":{}}}`, data.Resource.ValueString())
	if string(byteConfig) != "{}" && string(byteConfig) != emptyConfigString {
		config, err := stripResourceArgsAndResourceFields(string(byteConfig), runQueue.DefaultResourceConfig.Resource)
		if err != nil {
			resp.Diagnostics.AddError("Error stripping resource args and fields", err.Error())
			return
		}
		data.ResourceConfig = types.StringValue(config)
	}

	externalLinks, externalLinksDiags := convertExternalLinksListToMap(runQueue.ExternalLinks)
	resp.Diagnostics.Append(externalLinksDiags...)

	if len(runQueue.DefaultResourceConfig.TemplateVariables) > 0 {
		tvMap, err := templateVarsWithNamesListToMap(runQueue.DefaultResourceConfig.TemplateVariables)
		if err != nil {
			resp.Diagnostics.AddError("Error converting template variables", err.Error())
			return
		}
		tvBytes, err := json.Marshal(tvMap)
		if err != nil {
			resp.Diagnostics.AddError("Error marshalling template variables", err.Error())
			return
		}
		stringTemplateVariables := string(tvBytes)
		normalizedTemplateVariables, err := normalizeTemplateVariables(&stringTemplateVariables)
		if err != nil {
			resp.Diagnostics.AddError("Error normalizing template variables", err.Error())
			return
		}
		if normalizedTemplateVariables != nil {
			data.TemplateVariables = types.StringValue(*normalizedTemplateVariables)
		}
	}
	data.ExternalLinks = externalLinks

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "updated a run queue resource")
}

func (r *RunQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RunQueueResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	externalLinks, err := convertExternalLinksMapToInputType(data.ExternalLinks.Elements())
	if err != nil {
		resp.Diagnostics.AddError("Error converting external links", err.Error())
		return
	}

	prioritizationMode := data.PrioritizationMode.ValueStringPointer()
	if prioritizationMode == nil {
		defaultPrioritizationMode := "V0"
		prioritizationMode = &defaultPrioritizationMode
	}

	// Inject resource args and fields into the resource config backend expects wrapped in these fields
	resourceConfig, err := injectResourceArgsAndResourceFields(data.ResourceConfig.ValueString(), data.Resource.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error injecting resource args and fields", err.Error())
		return
	}

	normalizedTemplateVariables, err := normalizeTemplateVariables(data.TemplateVariables.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error normalizing template variables", err.Error())
		return
	}
	if normalizedTemplateVariables != nil {
		data.TemplateVariables = types.StringValue(*normalizedTemplateVariables)
	}

	input := UpsertRunQueueInput{
		QueueName:          data.Name.ValueString(),
		EntityName:         data.EntityName.ValueString(),
		ProjectName:        "model-registry",
		ResourceType:       data.Resource.ValueString(),
		ResourceConfig:     resourceConfig,
		TemplateVariables:  data.TemplateVariables.ValueStringPointer(),
		PrioritizationMode: prioritizationMode,
		ExternalLinks:      externalLinks,
	}

	result, err := upsertRunQueue(ctx, input, r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating run queue",
			"Could not create run queue, unexpected error: "+err.Error(),
		)
		return
	}

	if len(result.UpsertRunQueue.Errors) > 0 {
		resp.Diagnostics.AddWarning("Config schema validation errors", strings.Join(result.UpsertRunQueue.Errors, ", "))
	}
	id := generateCompositeID(data.EntityName.ValueString(), data.Name.ValueString())
	data.Id = types.StringValue(id)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunQueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RunQueueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	entityName, queueName, err := parseCompositeID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing composite ID", err.Error())
		return
	}
	runQueue, err := readRunQueueHelper(entityName, queueName, ctx, *r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading run queue id for delete",
			"Could not read run queue id, unexpected error: "+err.Error(),
		)
		return
	} else if runQueue == nil {
		resp.Diagnostics.AddError(
			"Error reading run queue id for delete",
			"Could not read run queue id, unexpected error: run queue not found",
		)
		return
	}

	// Create the GraphQL request to delete the run_queue using the ID
	gqlReq := graphql.NewRequest(`
		mutation DeleteRunQueues($queueIDs: [ID!]!) {
			deleteRunQueues(input:{queueIDs: $queueIDs}) {
				success
			}
		}
	`)

	ids := []string{runQueue.ID}
	gqlReq.Var("queueIDs", ids)

	var result struct {
		DeleteRunQueues struct {
			Success bool `json:"success"`
		} `json:"deleteRunQueues"`
	}

	if err := r.client.Run(ctx, gqlReq, &result); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting run queue",
			"Could not delete run queue, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.DeleteRunQueues.Success {
		resp.Diagnostics.AddError(
			"Failed to delete run queue",
			"The API did not confirm the deletion of the run queue.",
		)
		return
	}

	tflog.Trace(ctx, "deleted a run queue resource")

	resp.State.RemoveResource(ctx)
}

func (r *RunQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func upsertRunQueue(ctx context.Context, input UpsertRunQueueInput, client *GraphQLClientWithHeaders) (UpsertRunQueueResponse, error) {
	gqlReq := graphql.NewRequest(`
		mutation UpsertRunQueue(
			$entityName: String!,
			$projectName: String!,
			$queueName: String!,
			$resourceType: String!,
			$resourceConfig: JSONString!,
			$templateVariables: JSONString,
			$prioritizationMode: RunQueuePrioritizationMode,
			$externalLinks: JSONString,
		) { 
			upsertRunQueue(input: {
				entityName: $entityName,
				projectName: $projectName,
				queueName: $queueName,
				resourceType: $resourceType,
				resourceConfig: $resourceConfig,
				templateVariables: $templateVariables,
				prioritizationMode: $prioritizationMode,
				externalLinks: $externalLinks,
			}) {
				success
				configSchemaValidationErrors
			}
		}
	`)

	gqlReq.Var("entityName", input.EntityName)
	gqlReq.Var("projectName", "model-registry")
	gqlReq.Var("queueName", input.QueueName)
	gqlReq.Var("resourceType", input.ResourceType)
	gqlReq.Var("resourceConfig", input.ResourceConfig)
	gqlReq.Var("templateVariables", input.TemplateVariables)
	gqlReq.Var("prioritizationMode", input.PrioritizationMode)
	gqlReq.Var("externalLinks", input.ExternalLinks)

	var result UpsertRunQueueResponse

	err := client.Run(ctx, gqlReq, &result)
	return result, err
}
