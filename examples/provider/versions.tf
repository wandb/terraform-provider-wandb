terraform {
  required_providers {
    wandb = {
      source  = "local/wandb"
      version = "0.0.1"
    }
  }
  required_version = ">=0.0.1"
}