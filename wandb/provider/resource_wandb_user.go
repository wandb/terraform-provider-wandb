package provider

import (
	"context"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWandbUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Resource for a user in Wandb",

		CreateContext: resourceWandbUserCreate,
		ReadContext:   resourceWandbUserRead,
		UpdateContext: resourceWandbUserUpdate,
		DeleteContext: resourceWandbUserDelete,

		Schema: map[string]*schema.Schema{
			"email": {
				// This description is used by the documentation generator and the language server.
				Description: "The email for the user",
				Type:        schema.TypeString,
				Required:    true,
			},
			"admin": {
				// This description is used by the documentation generator and the language server.
				Description: "True if the user should be an admin",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceWandbUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*Client)
	err := client.CreateUser(d.Get("email").(string), d.Get("admin").(bool))
	if err != nil {
		return diag.Errorf(err.Error())
	}

	d.SetId(d.Get("email").(string))
	return nil
}

func resourceWandbUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceWandbUserUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceWandbUserDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
