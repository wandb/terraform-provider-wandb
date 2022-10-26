package provider

import (
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
				Optional:    false,
			},
			"organization_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The cloud storage bucket to use for the team",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage_bucket_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The cloud storage bucket to use for the team",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceWandbTeamCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*Client)
	err := client.CreateTeam(d.Get("organization_name").(string), d.Get("team_name").(string), d.Get("bucket_name").(string), d.Get("bucket_provider").(string))
	if err != nil {
		return diag.Errorf(err.Error())
	}
	return nil
}

func resourceWandbTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
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
