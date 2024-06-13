terraform {
  required_providers {
    wandb = {
      source  = "wandb/wandb"
      version = "0.1.0"
    }
  }
}

provider "wandb" {
}

resource "wandb_run_queue" "tf_example" {
  name        = "example-queue33"
  entity_name = "kylegoyette"

  resource = "kubernetes"

  resource_config = jsonencode({
    apiVersion = "batch/v1",
    kind       = "Job",
    metadata = {
      name = "{{example-variable}}"
    },
    spec = {
      template = {
        spec = {
          containers = [{
            name    = "example-container",
            image   = "example-image",
            command = ["example-command"]
          }],
          restartPolicy = "Never"
        }
      }
    }
  })

  template_variables = jsonencode({
    variable1 = {
      name        = "example-variable",
      description = "An example variable",
      schema = {
        type    = "string",
        minimum = 1,
        maximum = 10,
        enum    = ["option1", "option2"]
        default = "option1"
      }
    }
  })

  prioritization_mode = "V0"
  external_links = {
    "label" : "https://example.com",
    "label2" : "https://example2.com"
  }

}
