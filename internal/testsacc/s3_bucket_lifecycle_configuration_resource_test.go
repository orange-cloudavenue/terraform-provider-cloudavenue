package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &S3BucketLifecycleConfigurationResource{}

const (
	S3BucketLifecycleConfigurationResourceName = testsacc.ResourceName("cloudavenue_s3_bucket_lifecycle_configuration")
)

type S3BucketLifecycleConfigurationResource struct{}

func NewS3BucketLifecycleConfigurationResourceTest() testsacc.TestACC {
	return &S3BucketLifecycleConfigurationResource{}
}

// GetResourceName returns the name of the resource.
func (r *S3BucketLifecycleConfigurationResource) GetResourceName() string {
	return S3BucketLifecycleConfigurationResourceName.String()
}

func (r *S3BucketLifecycleConfigurationResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[S3BucketVersioningConfigurationResourceName]().GetDefaultConfig)
	return
}

func (r *S3BucketLifecycleConfigurationResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		// * Specifying a filter using key prefixes
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
						bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
					  
						rules = [{
						  id = {{ generate . "rule_id" }}
					  
						  filter = {
							prefix = "logs/"
						  }
					  
						  noncurrent_version_expiration = {
							noncurrent_days = 90
						  }
					  
						  status = "Enabled"
						}]
					  }`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "logs/"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.noncurrent_version_expiration.noncurrent_days", "90"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
							bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
						  
							rules = [{
							  id = {{ get . "rule_id" }}
						  
							  filter = {
								prefix = "config/"
							  }
						  
							  noncurrent_version_expiration = {
								noncurrent_days = 180
							  }
						  
							  status = "Enabled"
							}]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "config/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.noncurrent_version_expiration.noncurrent_days", "180"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
							bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
						  
							rules = [{
							  id = {{ get . "rule_id" }}
						  
							  filter = {
								prefix = "config/"
							  }
						  
							  noncurrent_version_expiration = {
								noncurrent_days = 180
							  }
						  
							  status = "Enabled"
							},
							{
								id = {{ generate . "rule_id_2" }}
							
								filter = {
								  prefix = "cache/"
								}
							
								noncurrent_version_expiration = {
								  noncurrent_days = 1
								}
							
								status = "Enabled"
							  }]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "config/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.noncurrent_version_expiration.noncurrent_days", "180"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.id", testsacc.GetValueFromTemplate(resourceName, "rule_id_2")),
							resource.TestCheckResourceAttr(resourceName, "rules.1.filter.prefix", "cache/"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.noncurrent_version_expiration.noncurrent_days", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.status", "Enabled"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				// ! Destroy
				Destroy: true,
			}
		},
		// * Example with a tag filter
		"example_filter_tag": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_filter_tag" {
						bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
					  
						rules = [{
						  id = {{ generate . "rule_id" }}
					  
						  filter = {
							tag = {
								key   = "tag1"
								value = "value1"
							}
						  }
					  
						  noncurrent_version_expiration = {
							noncurrent_days = 90
						  }
					  
						  status = "Enabled"
						}]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.tag.key", "tag1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.tag.value", "value1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.noncurrent_version_expiration.noncurrent_days", "90"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_filter_tag" {
							bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
						  
							rules = [{
							  id = {{ get . "rule_id" }}
						  
							  filter = {
								tag = {
									key   = "tag1-updated"
									value = "value1-updated"
								}
							  }
						  
							  noncurrent_version_expiration = {
								noncurrent_days = 180
							  }
						  
							  status = "Enabled"
							}]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.tag.key", "tag1-updated"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.tag.value", "value1-updated"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.noncurrent_version_expiration.noncurrent_days", "180"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				// ! Destroy
				Destroy: true,
			}
		},
		// * Example Specifying a filter based on tags and prefix
		"example_filter_tags_and_prefix": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_filter_tags_and_prefix" {
						bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
					  
						rules = [{
						  id = {{ generate . "rule_id" }}
					  
						  filter = {
							and = {
								prefix                   = "logs/"
								tags = [
									{
										key   = "tag1"
										value = "value1"
									}
								]
							}
						  }

						  noncurrent_version_expiration = {
							noncurrent_days = 180
						  }
					  
						  status = "Enabled"
						}]
					  }`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.key", "tag1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.value", "value1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.prefix", "logs/"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_filter_tags_and_prefix" {
							bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
						  
							rules = [{
							  id = {{ get . "rule_id" }}
						  
							  filter = {
								and = {
									prefix                   = "log/"
									tags = [
										{
											key   = "tag1"
											value = "value1"
										},
										{
											key   = "tag2"
											value = "value2"
										}
									]
								}
							  }

							  noncurrent_version_expiration = {
								noncurrent_days = 180
							  }
						  
							  status = "Enabled"
							}]
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.key", "tag1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.value", "value1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.1.key", "tag2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.1.value", "value2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.prefix", "log/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_filter_tags_and_prefix" {
							bucket = cloudavenue_s3_bucket.examplewithobjectlock.name
						  
							rules = [{
							  id = {{ get . "rule_id" }}
						  
							  filter = {
								and = {
									prefix                   = "log/"
									tags = [
										{
											key   = "tag1"
											value = "value1"
										},
										{
											key   = "tag2"
											value = "value2"
										}
									]
								}
							  }

							  noncurrent_version_expiration = {
								noncurrent_days = 180
							  }
							  
							  status = "Enabled"
							},
							{
								id = {{ generate . "rule_id_2" }}
							
								filter = {
									and = {
										prefix                   = "cache/"
										tags = [
											{
												key   = "tag1-cache"
												value = "value1-cache"
											}
										]
									}
								  }

								  noncurrent_version_expiration = {
									noncurrent_days = 180
								  }

								status = "Enabled"
							}]

							
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
							// ? 0
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.key", "tag1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.0.value", "value1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.1.key", "tag2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.tags.1.value", "value2"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.and.prefix", "log/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
							// ? 1
							resource.TestCheckResourceAttr(resourceName, "rules.1.id", testsacc.GetValueFromTemplate(resourceName, "rule_id_2")),
							resource.TestCheckResourceAttr(resourceName, "rules.1.filter.and.tags.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.filter.and.tags.0.key", "tag1-cache"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.filter.and.tags.0.value", "value1-cache"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.filter.and.prefix", "cache/"),
							resource.TestCheckResourceAttr(resourceName, "rules.1.status", "Enabled"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				// ! Destroy
				Destroy: true,
			}
		},
		// * Example Lifecycle Configuration for a bucket with versioning
		"example_with_versioning": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "bucket"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_with_versioning" {
							bucket = cloudavenue_s3_bucket_versioning_configuration.example.bucket
							
							rules = [{
							  id = {{ generate . "rule_id" }}
							
							  filter = {
								prefix = "logs/"
							  }
							
							  expiration = {
								  days = 90
							  }
						
							  status = "Enabled"
							}]
						}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
						resource.TestCheckResourceAttr(resourceName, "rules.0.expiration.days", "90"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "logs/"),
						resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_with_versioning" {
								bucket = cloudavenue_s3_bucket_versioning_configuration.example.bucket
								
								rules = [{
								  id = {{ generate . "rule_id" }}
								
								  filter = {
									prefix = "logs/"
								  }
								
								  expiration = {
									  date = "2027-01-01T00:00:00Z"
								  }
							
								  status = "Enabled"
								}]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.expiration.date", "2027-01-01T00:00:00Z"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "logs/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_s3_bucket_lifecycle_configuration" "example_with_versioning" {
								bucket = cloudavenue_s3_bucket_versioning_configuration.example.bucket
								
								rules = [{
								  id = {{ generate . "rule_id" }}
								
								  filter = {
									prefix = "logs/"
								  }
								
								  expiration = {
									expired_object_delete_marker = true
								  }
							
								  status = "Enabled"
								}]
							}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.id", testsacc.GetValueFromTemplate(resourceName, "rule_id")),
							resource.TestCheckResourceAttr(resourceName, "rules.0.expiration.expired_object_delete_marker", "true"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.filter.prefix", "logs/"),
							resource.TestCheckResourceAttr(resourceName, "rules.0.status", "Enabled"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"bucket"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				// ! Destroy
				Destroy: true,
			}
		},
	}
}

func TestAccS3BucketLifecycleConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&S3BucketLifecycleConfigurationResource{}),
	})
}
