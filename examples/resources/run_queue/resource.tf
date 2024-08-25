resource "wandb_run_queue" "tf_example" {
  name        = "example-queue"
  entity_name = "<entity-name>"

  resource = "kubernetes"

  resource_config = jsonencode({
    apiVersion = "batch/v1",
    kind       = "Job",
    metadata = {
      name = "{{example_variable}}"
    },
    spec = {
      template = {
        spec = {
          containers = [{
            name = "example-container",
          }],
          restartPolicy = "Never"
        }
      }
    }
  })

  template_variables = jsonencode({
    example_variable = {
      description = "An example variable",
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