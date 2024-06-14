package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"wandb": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
