package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"wandb": func() (*schema.Provider, error) {
		return initAccProvider(), nil
	},
}

func initAccProvider() *schema.Provider {
	p := New("dev")()
	p.ConfigureContextFunc = configure("dev", p)

	return p
}

func testAccProvider(t *testing.T, accProviders map[string]func() (*schema.Provider, error)) func() (*schema.Provider, error) {
	accProvider, ok := accProviders["wandb"]
	if !ok {
		t.Fatal("could not find wandb provider")
	}
	return accProvider
}

func TestProvider(t *testing.T) {
	provider := New("dev")()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
