provider "wandb" {
  host = "https://t4l.wandb.ml"
  api_key = "local-bb3a44320434bd75aa88725906cf51e8b1f541ed"
}

resource "wandb_team" "example" {
  team_name = "team-tmp2"
  organization_name = ""
  storage_bucket_name = "hackweek-intense-kitten"
  storage_bucket_provider = "GCP"
}