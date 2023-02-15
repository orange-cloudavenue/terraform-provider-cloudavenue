# CloudAvenue Plan Modifier Bool Helper

This helper is used to modify a boolean value in a plan.

## Helpers Available

### `SetDefault`

This helper is used to set a default value for a boolean.

```go
// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "enabled": schema.BoolAttribute{
                Required:            true,
                MarkdownDescription: "Enable or disable ...",
                PlanModifiers: []planmodifier.Bool{
                    int64pm.SetDefault(true),
                },
            },
```

### `SetDefaultEnvVar`

This helper is used to set a default value for a boolean from an environment variable.

```sh
export CAV_VAR_DEFAULT_NAME="true"
```

```go
// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "enabled": schema.BoolAttribute{
                Required:            true,
                MarkdownDescription: "Enable or disable ...",
                PlanModifiers: []planmodifier.Bool{
                    int64pm.SetDefaultEnvVar("CAV_VAR_DEFAULT_NAME"),
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
                Required:            true,
                MarkdownDescription: "Enable or disable ...",
                PlanModifiers: []planmodifier.Bool{
                    int64pm.SetDefaultFunc(int64pm.DefaultFunc(func(ctx context.Context, req planmodifier.BoolRequest, resp *int64pm.DefaultFuncResponse) {
                        if os.Getenv("CAV_VAR_1") == "foo" && os.Getenv("CAV_VAR_2") == "bar" {
                            resp.Value = true
                            return
                        }
                    })),
                },
            },
```
