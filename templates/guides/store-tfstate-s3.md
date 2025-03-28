---
page_title: "Using CloudAvenue S3 to store Terraform State"
---

# Terraform Backend

This page describes how to configure a backend by adding the backend block to your configuration with the Terraform CloudAvenue Provider.

Terraform provides the option to set up a [“backend”](https://developer.hashicorp.com/terraform/language/backend) of the `state` data files.

This option allows you to handle the state and the way certain operations are executed.

Backends can store the state remotely and protect it with locks to prevent corruption;
it makes it possible for a team to work with ease, or, for instance, to run Terraform within a pipeline.

## Create your S3 bucket

You can create your S3 bucket resource using terraform itself .

If you have already one bucket you can step over to [Create your infrastructure with the cloudavenue provider](#create-your-infrastructure-with-the-cloudavenue-provider).

```hcl
terraform {
  required_providers {
    cloudavenue = {
      source  = "orange-cloudavenue/cloudavenue"
      version = "~> 0.30.0"
    }
  }
}

provider "cloudavenue" {
  # ...
}

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~#
# CREATE AN S3 BUCKET TO USE IT AS A TERRAFORM BACKEND #
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~#

resource "cloudavenue_s3_bucket" "example" {
  name = "example"
}

resource "cloudavenue_s3_bucket_policy" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action = [
          "s3:ListBucket",
        ]
        Resource = [
          "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}",
        ]
      },
      {
        Effect    = "Allow"
        Principal = "*"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = [
          "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}/path/to/your/terraform",
        ]
      },
      # If your following recommendation and your enable use_lockfile add the following policy
      {
        Effect    = "Allow"
        Principal = "*"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject"
        ]
        Resource = [
          "arn:aws:s3:::${cloudavenue_s3_bucket.example.name}/path/to/your/terraform/key.tflock",
        ]
      }
    ]
  })
}

# If you are enabling versioning on the bucket for the first time, it's recommends that you wait for 15 minutes after 
# enabling versioning before issuing write operations (PUT or DELETE) on objects in the bucket.

resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  status = "Enabled"
}

resource "cloudavenue_s3_credential" "example" {
  print_token = true
}
```

and deploy it:

```shell
terraform plan -out "planfile" ; terraform apply -input=false -auto-approve "planfile"
```

The resource `cloudavenue_s3_credential` will print a `access_key` and `secret_key` to access the bucket.

## Secrets

Hashicorp offers several methods to keep your secrets. Please check the Terraform [partial configuration](https://developer.hashicorp.com/terraform/language/backend#partial-configuration) for this topic.

## Create your infrastructure with the cloudavenue provider

```hcl
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# CREATE AN BACKEND TYPE S3
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
terraform {
  backend "s3" {
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    use_path_style              = true
    skip_s3_checksum            = true

    use_lockfile = true

    key    = "/path/to/your/terraform/prod.tfstate"
    region = "region01"
    bucket = "YourBucketName"
    endpoints = {
      s3 = "https://s3-region01.cloudavenue.orange-business.com"
    }

    # Please use a better approach with the flag -backend-config=PATH or a Vault configuration
    access_key = "YourAccessKey"
    secret_key = "Y0urS3cretK3y"
  }
}

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# CREATE YOUR INFRASTRUCTURE
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
resource "cloudavenue_vdc" "example" {
  name                  = "example"
  cpu_allocated         = 22000
  memory_allocated      = 30
  cpu_speed_in_mhz      = 2200
  billing_model         = "PAYG"
  disponibility_class   = "ONE-ROOM"
  service_class         = "STD"
  storage_billing_model = "PAYG"

  storage_profiles = [
    {
      class   = "gold"
      default = true
      limit   = 500
    }
  ]
}

# the rest of your configuration and resources to deploy
```

## Migrating the state

Considering you have already running infrastructure you want to use the `backend` option.

All we need to do is initialize Terraform passing the backend configuration.

Terraform should ask if you want to migrate from local to the new remote backend.

Answer the prompt `yes`, and your state will migrate.

```shell
terraform init  -migrate-state
```

## What about locking?

Most of the remote [backends](https://developer.hashicorp.com/terraform/language/backend#backend-types) natively support locking. To run terraform apply, Terraform will automatically acquire a lock;
if someone else is already running apply, they will already have the lock, and you will have to wait.
You can run apply with the `-lock-timeout=<TIME>` parameter to tell Terraform to wait up to TIME for a lock to be released (e.g., `-lock-timeout=10m` will wait for 10 minutes).

The Lock method prevents opening the state file while already in use.
