---
page_title: "cloudavenue_s3_bucket_lifecycle_configuration Resource - cloudavenue"
subcategory: "S3 (Object Storage)"
description: |-
  The cloudavenue_s3_bucket_lifecycle_configuration resource allows you to manage lifecycle configuration of an S3 bucket.
---

# cloudavenue_s3_bucket_lifecycle_configuration (Resource)

Provides an independent configuration resource for S3 bucket [lifecycle configuration](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lifecycle-mgmt.html).

The `cloudavenue_s3_bucket_lifecycle_configuration` resource allows you to manage lifecycle configuration of an S3 bucket.

An S3 Lifecycle configuration consists of one or more Lifecycle rules. Each rule consists of the following:

* An ID that identifies the rule. The ID must be unique within the configuration.
* A Status that indicates whether the rule is currently being applied.
* A Filter that identifies a subset of objects to which the rule applies.
* One or more Lifecycle actions that you want S3 to perform on the objects identified by the Filter.

For more information about Lifecycle configuration, see [Lifecycle Configuration Elements](https://docs.aws.amazon.com/AmazonS3/latest/userguide/intro-lifecycle-rules.html).

 ~> **NOTE** S3 Buckets only support a single lifecycle configuration. Declaring multiple `cloudavenue_s3_bucket_lifecycle_configuration` resources to the same S3 Bucket will cause a perpetual difference in configuration.

 ~> **NOTE** Lifecycle configurations may take some time to fully propagate to all CloudAvenue S3 systems. Running Terraform operations shortly after creating a lifecycle configuration may result in changes that affect configuration idempotence. See the S3 User Guide on [setting lifecycle configuration on a bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/how-to-set-lifecycle-configuration-intro.html).

## Examples Usage

### Specifying a filter using key prefixes

The Lifecycle rule applies to a subset of objects based on the key name prefix (`logs/`).

```hcl
resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
	bucket = cloudavenue_s3_bucket.example.name
  
	rules = [{
	  id = "rule_id_1"
  
	  filter = {
		prefix = "logs/"
	  }
  
	  noncurrent_version_expiration = {
		noncurrent_days = 90
	  }
  
	  status = "Enabled"
	}]
}
```

If you want to apply a Lifecycle action to a subset of objects based on different key name prefixes, specify separate rules.

```hcl
resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
	bucket = cloudavenue_s3_bucket.example.name
  
	rules = [
	{
	  id = "rule_id_1"
  
	  filter = {
		prefix = "config/"
	  }
  
	  noncurrent_version_expiration = {
		noncurrent_days = 180
	  }
  
	  status = "Enabled"
	},
	{
	  id = "rule_id_2"
	
	  filter = {
		prefix = "cache/"
	  }
	
	  noncurrent_version_expiration = {
		noncurrent_days = 10
	  }
	
	  status = "Enabled"
	}]
}
```

### Specifying a filter based on tag

The Lifecycle rule applies to a subset of objects based on the tag key and value (`tag1` and `value1`).

```hcl
resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
	bucket = cloudavenue_s3_bucket.example.name
  
	rules = [{
	  id = "rule_id_1"
  
	  filter = {
		tag = {
			key   = "tag1"
			value = "value1"
		}
	  }
  
	  expiration = {
	    days = 90
	  }
  
	  status = "Enabled"
	}]
}
```

### Specifying a filter based on tags range and prefix

The Lifecycle rule applies to a subset of objects based on the tag key and value (`tag1` and `value1`) and the key name prefix (`logs/`).

```hcl
resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
	bucket = cloudavenue_s3_bucket.example.name
	
	rules = [{
	  id = "rule_id_1"
	
	  filter = {
		and {
			prefix = "logs/"
			tags = [
				{
					key   = "tag1"
					value = "value1"
				}
			]
		}
	  }

	  expiration = {
	    days = 90
	  }
	
	  status = "Enabled"
	}]
}
```

### Creating a Lifecycle Configuration for a bucket with versioning

```hcl
resource "cloudavenue_s3_bucket" "example" {
	name = "example"
}

resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  status = "Enabled"
}

resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
	bucket = cloudavenue_s3_bucket_versioning_configuration.example.bucket
	
	rules = [{
	  id = "rule_id_1"
	
	  filter = {
		prefix = "logs/"
	  }
	
	  expiration {
	  	days = 90
	  }

	  status = "Enabled"
	}]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket` (String) (ForceNew) The name of the bucket.
- `rules` (Attributes List) Rules that define lifecycle configuration. List must contain at least 1 elements. (see [below for nested schema](#nestedatt--rules))

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID is a bucket name.

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `filter` (Attributes) Configuration block used to identify objects that a Lifecycle Rule applies to. (see [below for nested schema](#nestedatt--rules--filter))
- `id` (String) Unique identifier for the rule. String length must be between 1 and 255.

Optional:

- `abort_incomplete_multipart_upload` (Attributes) Configuration block that specifies the days since the initiation of an incomplete multipart upload that S3 will wait before permanently removing all parts of the upload. (see [below for nested schema](#nestedatt--rules--abort_incomplete_multipart_upload))
- `expiration` (Attributes) Configuration block that specifies the expiration for the lifecycle of the object in the form of date, days and, whether the object has a delete marker. Ensure that if an attribute is set, these are not set: "[<.expiration,<.noncurrent_version_expiration]". (see [below for nested schema](#nestedatt--rules--expiration))
- `noncurrent_version_expiration` (Attributes) Configuration block that specifies when noncurrent object versions expire. Ensure that if an attribute is set, these are not set: "[<.expiration,<.noncurrent_version_expiration]". (see [below for nested schema](#nestedatt--rules--noncurrent_version_expiration))
- `status` (String) Whether the rule is currently being applied. Value must be one of : `Enabled`, `Disabled`. Value defaults to `Enabled`.

<a id="nestedatt--rules--filter"></a>
### Nested Schema for `rules.filter`

Optional:

- `and` (Attributes) Configuration block used to apply a logical AND to two or more predicates. The Lifecycle Rule will apply to any object matching all the predicates configured inside the and block. (see [below for nested schema](#nestedatt--rules--filter--and))
- `prefix` (String) Match objects with this prefix.
- `tag` (Attributes) Specifies object tag key and value. (see [below for nested schema](#nestedatt--rules--filter--tag))

<a id="nestedatt--rules--filter--and"></a>
### Nested Schema for `rules.filter.and`

Optional:

- `prefix` (String) Match objects with this prefix.
- `tags` (Attributes List) Specifies object tag key and value. (see [below for nested schema](#nestedatt--rules--filter--and--tags))

<a id="nestedatt--rules--filter--and--tags"></a>
### Nested Schema for `rules.filter.and.tags`

Required:

- `key` (String) Object tag key.
- `value` (String) Object tag value.



<a id="nestedatt--rules--filter--tag"></a>
### Nested Schema for `rules.filter.tag`

Required:

- `key` (String) Object tag key.
- `value` (String) Object tag value.



<a id="nestedatt--rules--abort_incomplete_multipart_upload"></a>
### Nested Schema for `rules.abort_incomplete_multipart_upload`

Optional:

- `days_after_initiation` (Number) Number of days after which S3 aborts an incomplete multipart upload.


<a id="nestedatt--rules--expiration"></a>
### Nested Schema for `rules.expiration`

Optional:

- `date` (String) Date the object is to be moved or deleted. The date value must be in [RFC3339 full-date format](https://datatracker.ietf.org/doc/html/rfc3339#section-5.6) e.g. `2023-10-10T00:00:00Z`. Ensure that one and only one attribute from this collection is set : `date`, `days`, `expired_object_delete_marker`.
- `days` (Number) Lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer. Ensure that one and only one attribute from this collection is set : `date`, `days`, `expired_object_delete_marker`.
- `expired_object_delete_marker` (Boolean) Indicates whether S3 will remove a delete marker with no noncurrent versions. If set to `true`, the delete marker will be expired, if set to `false` the policy takes no action. Ensure that one and only one attribute from this collection is set : `date`, `days`, `expired_object_delete_marker`. Ensure that if an attribute is set, these are not set: "[<.<.filter.tag,<.<.filter.and.tags]".


<a id="nestedatt--rules--noncurrent_version_expiration"></a>
### Nested Schema for `rules.noncurrent_version_expiration`

Optional:

- `newer_noncurrent_versions` (Number) Number of noncurrent versions S3 will retain. Value must be at least 0.
- `noncurrent_days` (Number) Number of days an object is noncurrent before S3 can perform the associated action. Must be a positive integer. Value must be at least 1.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

 -> **Timeout** Default timeout is **5 minutes**.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_s3_bucket_lifecycle_configuration.example bucketName
```