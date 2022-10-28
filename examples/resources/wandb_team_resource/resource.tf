resource "wandb_team" "example" {
  team_name = "foo"
  organization_name = "my org"
  storage_bucket_name = "my-bucket"
  storage_bucket_provider = "gcs"
}