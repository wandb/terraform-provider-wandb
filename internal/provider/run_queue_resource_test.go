package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRunQueueResource_basic(t *testing.T) {
	resourceName := "wandb_run_queue.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		CheckDestroy:             testAccCheckRunQueueResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRunQueueResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRunQueueResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-queue"),
					resource.TestCheckResourceAttr(resourceName, "resource", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "prioritization_mode", "V0"),
					resource.TestCheckResourceAttr(resourceName, "external_links.label", "https://example.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("WANDB_API_KEY"); v == "" {
		t.Fatal("WANDB_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("WANDB_BASE_URL"); v == "" {
		t.Fatal("WANDB_BASE_URL must be set for acceptance tests")
	}
}

func testAccCheckRunQueueResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		// Optionally: Add any checks to verify the resource exists in the system
		return nil
	}
}

func testAccCheckRunQueueResourceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wandb_run_queue" {
			continue
		}

		client := newGraphQLClient()

		runQueue, err := readRunQueueHelper(rs.Primary.Attributes["entity_name"], rs.Primary.Attributes["name"], context.Background(), *client)
		if err != nil {
			if err.Error() == "queue not found" {
				continue
			}
		}

		if runQueue != nil {
			return fmt.Errorf("run queue still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccRunQueueResourceConfig() string {
	return `
terraform {
  required_providers {
    wandb = {
      source = "wandb/wandb"
    }
  }
}

provider "wandb" {
	version = "0.1.0"
}

resource "wandb_run_queue" "test" {
  name        = "example-queue"
  entity_name = "terraform-acceptance-test"

  resource = "kubernetes"

  resource_config = jsonencode({
      apiVersion = "batch/v1",
      kind       = "Job",
      metadata = {
        name = "{{exampleVariable}}"
      }
  })

  template_variables = jsonencode({
    exampleVariable = {
	  schema = {
        type = "string"
      }
    }
  })

  prioritization_mode = "V0"
  external_links = {
    "label" : "https://example.com",
    "label2" : "https://example2.com"
  }
}
`
}

func newGraphQLClient() *GraphQLClientWithHeaders {
	baseURL := os.Getenv("WANDB_BASE_URL")
	apiKey := os.Getenv("WANDB_API_KEY")
	headers := http.Header{}
	headers.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api:"+apiKey)))
	headers.Set("Content-Type", "application/json")
	return NewGraphQLClientWithHeaders(baseURL+"/graphql", headers)
}
