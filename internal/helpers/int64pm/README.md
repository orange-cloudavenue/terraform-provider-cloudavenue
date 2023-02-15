# CloudAvenue Plan Modifier Bool Helper

This helper is used to modify a boolean value in a plan.

## Helpers Available

### `SetDefault`

This helper is used to set a default value for a boolean.

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "disk_size": schema.Int64Attribute{
                Optional:            true,
                MarkdownDescription: "The size of the disk in MB.",
                PlanModifiers: []planmodifier.Int64{
                    int64pm.SetDefault(100),
                },
            },
```

### `SetDefaultEnvVar`

This helper is used to set a default value for a boolean from an environment variable.

```sh
export CAV_VAR_DEFAULT_DISK_SIZE="100"
```

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "disk_size": schema.Int64Attribute{
                Optional:            true,
                MarkdownDescription: "The size of the disk in MB.",
                PlanModifiers: []planmodifier.Int64{
                    int64pm.SetDefaultEnvVar("CAV_VAR_DEFAULT_DISK_SIZE"),
                },
            },
```

### `SetDefaultFunc`

This helper is used to set a default value for a boolean using a function.

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "enabled": schema.BoolAttribute{
                Optional:            true,
                MarkdownDescription: "The size of the disk in MB.",
                PlanModifiers: []planmodifier.Bool{
                    int64pm.SetDefaultFunc(int64pm.DefaultFunc(func(ctx context.Context, req planmodifier.Int64Request, resp *int64pm.DefaultFuncResponse) {
                        resp.Value = req.PlanValue * 1024
                    })),
                },
            },
```
