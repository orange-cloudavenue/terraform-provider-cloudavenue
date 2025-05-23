---
page_title: "cloudavenue_s3_bucket_versioning_configuration Resource - cloudavenue"
subcategory: "S3 (Object Storage)"
description: |-
  The cloudavenue_s3_bucket_versioning_configuration resource allows you to manage the versioning configuration of an S3 bucket. Provides a resource for controlling versioning on an S3 bucket. Deleting this resource will either suspend versioning on the associated S3 bucket or simply remove the resource from Terraform state if the associated S3 bucket is unversioned. For more information, see How S3 versioning works https://docs.aws.amazon.com/AmazonS3/latest/userguide/manage-versioning-examples.html.
---

# cloudavenue_s3_bucket_versioning_configuration (Resource)

The `cloudavenue_s3_bucket_versioning_configuration` resource allows you to manage the versioning configuration of an S3 bucket. Provides a resource for controlling versioning on an S3 bucket. Deleting this resource will either suspend versioning on the associated S3 bucket or simply remove the resource from Terraform state if the associated S3 bucket is unversioned. For more information, see [How S3 versioning works](https://docs.aws.amazon.com/AmazonS3/latest/userguide/manage-versioning-examples.html).

 ~> **NOTE:** If you are enabling versioning on the bucket for the first time, it's recommends that you wait for 15 minutes after enabling versioning before issuing write operations (`PUT` or `DELETE`) on objects in the bucket.

## Examples Usage

### With Versioning Enabled
```hcl
resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  status = "Enabled"
}
```

### With Versioning Suspended
```hcl
resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  status = "Suspended"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of the bucket.
- `status` (String) Versioning state of the bucket. Value must be one of : `Enabled`, `Suspended`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID is a bucket name.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_s3_bucket_versioning_configuration.example bucket-name
```