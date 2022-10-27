package provider

import (
	"context"
	"time"

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

	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	return diag.Errorf("not implemented")
}

func resourceWandbTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	client := NewClient("https://api.wandb.ai", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913", time.Second*10)
	client.ReadTeam(d.Get("team_name").(string))

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
