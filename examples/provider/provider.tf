terraform {
  required_providers {
    wandb = {
      source  = "wandb/wandb"
      version = "0.1.0"
    }
  }
}

provider "wandb" {
  base_url = "https://api.wandb.ai"
}
