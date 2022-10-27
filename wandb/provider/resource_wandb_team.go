package provider

import (
	"fmt"
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWandbTeam() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Resource for a team in Wandb",

		CreateContext: resourceWandbTeamCreate,
		ReadContext:   resourceWandbTeamRead,
		UpdateContext: resourceWandbTeamUpdate,
		DeleteContext: resourceWandbTeamDelete,

		Schema: map[string]*schema.Schema{
			"team_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The name for the team",
				Type:        schema.TypeString,
				Required:    true,
			},
			"organization_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The organization name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage_bucket_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The cloud storage bucket to use for the team",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage_bucket_provider": {
				// This description is used by the documentation generator and the language server.
				Description: "The cloud storage bucket provider",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"id": {
				Description: "The ID of the team",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "When the team was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "When the team was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceWandbTeamCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*Client)
	team, err := client.CreateTeam(d.Get("organization_name").(string), d.Get("team_name").(string), d.Get("storage_bucket_name").(string), d.Get("storage_bucket_provider").(string))
	tflog.Trace(ctx, team.Name)
	if err != nil {
		return diag.Errorf(err.Error())
	}
	d.SetId(team.Id)
	d.Set("team_name", team.Name)
	d.Set("created_at", team.CreatedAt)
	d.Set("updated_at", team.UpdatedAt)
	d.Set("organization_name", d.Get("organization_name").(string))
	d.Set("storage_bucket_name", d.Get("storage_bucket_name").(string))
	d.Set("storage_bucket_provider", d.Get("storage_bucket_provider").(string))
	return nil
}

func resourceWandbTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	client := meta.(*Client)
	team, err := client.ReadTeam(d.Get("team_name").(string))
	if err != nil {
		return diag.Errorf(err.Error())
	}
	fmt.Println(string(team.Id))
	d.Set("id", team.Id)
	d.Set("created_at", team.CreatedAt)
	d.Set("updated_at", team.UpdatedAt)
	d.Set("organization_name", d.Get("organization_name").(string))
	d.Set("storage_bucket_name", d.Get("storage_bucket_name").(string))
	d.Set("storage_bucket_provider", d.Get("storage_bucket_provider").(string))

	return nil
}

func resourceWandbTeamUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceWandbTeamDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*Client)
	err := client.DeleteTeam(d.Get("team_name").(string))
	if err != nil {
		return diag.Errorf(err.Error())
	}
	return nil
}
