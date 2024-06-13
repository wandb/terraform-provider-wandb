package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"wandb": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6(New("test")())(), nil
		},
	}
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
