provider "wandb" {
  host = "api.wandb.ai"
  api_key = "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913"
}

resource "wandb_team" "example" {
  team_name = "foo"
  organization_name = "my org"
  storage_bucket_name = "my-bucket"
  storage_bucket_provider = "gcs"
}