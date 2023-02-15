# CloudAvenue Plan Modifier String Helper

This helper is used to modify a string value in a plan.

## Helpers Available

### `SetDefault`

This helper is used to set a default value for a string.

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Optional:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.SetDefault("default-name"),
                },
            },
```

### `SetDefaultEnvVar`

This helper is used to set a default value for a string from an environment variable.

```sh
export CAV_VAR_DEFAULT_NAME="default-name"
```

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Optional:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.SetDefaultEnvVar("CAV_VAR_DEFAULT_NAME"),
                },
            },
```

### `SetDefaultFunc`

This helper is used to set a default value for a string using a function.

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Optional:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.SetDefaultFunc(stringpm.DefaultFunc(func(ctx context.Context, req planmodifier.StringRequest, resp *stringpm.DefaultFuncResponse) {
                        if strings.Contains(req.PlanValue, "foo") {
                            resp.Value = "bar"
                            return
                        }
                    })),
                },
            },
```

### `SetDefaultEmptyString`

This helper is used to set a default empty value for a string.

```go
// Schema defines the schema for the resource.
func (r *xResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        (...)
            "name": schema.StringAttribute{
                Optional:            true,
                MarkdownDescription: "A name for ...",
                PlanModifiers: []planmodifier.String{
                    stringpm.SetDefaultEmptyString(),
                },
            },
```
