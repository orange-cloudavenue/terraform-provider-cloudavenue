---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudavenue_catalog_media Data Source - cloudavenue"
subcategory: ""
description: |-
  The catalog_media datasource provides a CloudAvenue Catalog media data source.
---

# cloudavenue_catalog_media (Data Source)

The `catalog_media` datasource provides a CloudAvenue Catalog media data source.

## Example Usage

```terraform
data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_media" "example" {
  catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].catalog_name
  name         = "debian-9.9.0-amd64-netinst.iso"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the media.

### Optional

- `catalog_id` (String) The ID of the catalog to which media file belongs.
- `catalog_name` (String) The name of the catalog to which media file belongs.

### Read-Only

- `created_at` (String) The creation date of the media.
- `description` (String) The description of the media.
- `id` (String) The ID of the catalog media.
- `is_iso` (Boolean) True if this media file is an Iso.
- `is_published` (Boolean) True if this media file is in a published catalog.
- `owner_name` (String) The name of the owner.
- `size` (Number) The size of the media in bytes.
- `status` (String) The media status.
- `storage_profile` (String) The name of the storage profile.


