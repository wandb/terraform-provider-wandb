package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

// WANDB_BASE_URL name of env var for base URL
const WandbBaseURLEnvName = "WANDB_BASE_URL"

// WANDB_API_KEY name of env var for API key
const WandbAPIKeyEnvName = "WANDB_API_KEY"

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"wandb_user": resourceWandbUser(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"wandb_team": resourceWandbTeam(),
			},
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(WandbBaseURLEnvName, "api.wandb.ai"),
					Description: "Wandb API URL. This can also be set via the WANDB_BASE_URL environment variable.",
				},
				"api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(WandbAPIKeyEnvName, nil),
					Description: "(Required unless validate is false) Wandb API key. This can also be set via the WANDB_API_KEY environment variable.",
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent
		defaultTimeout := time.Second * 10
		apiClient := NewClient(d.Get("host").(string), d.Get("api_key").(string), defaultTimeout)
		return apiClient, nil
	}
}
