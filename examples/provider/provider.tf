# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    wandb = {
      source = "wandb/wandb"
    }
  }
}

provider "wandb" {
  base_url = "https://api.wandb.ai"
}
