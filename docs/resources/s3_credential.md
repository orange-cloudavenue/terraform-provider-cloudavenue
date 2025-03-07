---
page_title: "cloudavenue_s3_credential Resource - cloudavenue"
subcategory: "S3 (Object Storage)"
description: |-
  The cloudavenue_s3_credential resource allows you to manage an access key and secret key for an S3 user.
---

# cloudavenue_s3_credential (Resource)

The `cloudavenue_s3_credential` resource allows you to manage an access key and secret key for an S3 user.

## Example Usage

```terraform
resource "cloudavenue_s3_credential" "example" {
  save_in_file = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `file_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of the file to store the API token. Value defaults to `token.json`.
- `print_token` (Boolean) <i style="color:red;font-weight: bold">(ForceNew)</i> If true, the API token will be printed in the console. Set this to true if you understand the security risks of using AccessKey/SecretKey and agree to creating them. This setting is only used when creating a new AccessKey/SecretKey and available only one time. Value defaults to `false`.
- `save_in_file` (Boolean) <i style="color:red;font-weight: bold">(ForceNew)</i> If true, the API token will be saved in a file. Set this to true if you understand the security risks of using AccessKey/SecretKey files and agree to creating them. This setting is only used when creating a new AccessKey/SecretKey and available only one time. Value defaults to `false`.
- `save_in_tfstate` (Boolean) <i style="color:red;font-weight: bold">(ForceNew)</i> If true, the SecretKey will be saved in the terraform state. Set this to true if you understand the security risks of using AccessKey/SecretKey and agree to creating them. This setting is only used when creating a new API token and available only one time. 

 !> **Warning:** This is a security risk and should only be used for testing purposes. Value defaults to `false`.

### Read-Only

- `access_key` (String, Sensitive) The Access Key.
- `id` (String) The ID of the credential. ID is a username and 4 first characters of the access key. (e.g. `username-1234`).
- `secret_key` (String, Sensitive) The Secret Key. Only Available if the `save_in_tfstate` is set to true.
- `username` (String) The username is configured at the provider level.

