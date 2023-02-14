# CloudAvenue Plan Modifier String Helper

This helper is used to modify a string value in a plan.

## Helpers Available

### `Default`

This helper is used to set a default value for a string.

```go
// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Required:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.Default("default-name"),
                },
            },
```

### `DefaultEnvVar`

This helper is used to set a default value for a string from an environment variable.

```sh
export CAV_VAR_DEFAULT_NAME="default-name"
```

```go
// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Required:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.DefaultEnvVar("CAV_VAR_DEFAULT_NAME"),
                },
            },
```
