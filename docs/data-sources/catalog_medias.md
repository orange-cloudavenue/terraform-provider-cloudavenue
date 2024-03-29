---
page_title: "cloudavenue_catalog_medias Data Source - cloudavenue"
subcategory: "Catalog"
description: |-
  The Catalog medias allows you to retrieve information about a medias in Cloud Avenue.
---

# cloudavenue_catalog_medias (Data Source)

The Catalog medias allows you to retrieve information about a medias in Cloud Avenue.

## Example Usage

```terraform
data "cloudavenue_catalogs" "test" {}

data "cloudavenue_catalog_medias" "example" {
  catalog_name = data.cloudavenue_catalogs.test.catalogs["catalog-example"].name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `catalog_id` (String) The ID of the catalog. Ensure that one and only one attribute from this collection is set : `catalog_name`, `catalog_id`.
- `catalog_name` (String) The name of the catalog. Ensure that one and only one attribute from this collection is set : `catalog_name`, `catalog_id`.

### Read-Only

- `id` (String) The ID of the medias.
- `medias` (Attributes Map) The map of medias. (see [below for nested schema](#nestedatt--medias))
- `medias_name` (List of String) The list of medias name.

<a id="nestedatt--medias"></a>
### Nested Schema for `medias`

Read-Only:

- `catalog_id` (String) The ID of the catalog.
- `catalog_name` (String) The name of the catalog.
- `created_at` (String) The date and time when the media was created.
- `description` (String) The description of the media.
- `id` (String) The ID of the media.
- `is_iso` (Boolean) `True` if the media is an ISO.
- `is_published` (Boolean) `True` if the media is published.
- `name` (String) The name of the media.
- `owner_name` (String) The name of the owner of the media.
- `size` (Number) The size of the media in bytes.
- `status` (String) The status of the media.
- `storage_profile` (String) The storage profile of the media.

