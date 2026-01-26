// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package pages_project

import (
	"context"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithUpgradeState = (*PagesProjectResource)(nil)

// UpgradeState returns state upgraders for handling schema version migrations.
// Version 0: v4 provider schema (pre-5.x) - blocks stored as lists (SDKv2 style)
// Version 1: v5 provider schema - single nested attributes
func (r *PagesProjectResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   resourceSchemaV0(ctx),
			StateUpgrader: upgradeStateFromV0,
		},
	}
}

// =============================================================================
// V0 Schema Definition (v4 provider)
// =============================================================================

// resourceSchemaV0 returns the v4 provider schema for pages_project.
// Key differences from v5:
//   - Blocks are represented as lists (SDKv2 style) not single nested attributes
//   - environment_variables and secrets are separate maps (merged to env_vars in v5)
//   - service_binding is a list (converted to services map in v5)
//   - kv_namespaces, d1_databases, r2_buckets, durable_object_namespaces are simple string maps
//     (wrapped in objects in v5)
//   - production_deployment_enabled renamed to production_deployments_enabled
func resourceSchemaV0(ctx context.Context) *schema.Schema {
	return &schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"account_id": schema.StringAttribute{
				Required: true,
			},
			"production_branch": schema.StringAttribute{
				Required: true,
			},
			// v4 SDKv2 stores blocks as lists
			"build_config": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"build_caching": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"build_command": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"destination_dir": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"root_dir": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"web_analytics_tag": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"web_analytics_token": schema.StringAttribute{
							Optional:  true,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
			"source": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional: true,
						},
						"config": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"deployments_enabled": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
									"owner": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"path_excludes": schema.ListAttribute{
										Optional:    true,
										Computed:    true,
										ElementType: types.StringType,
									},
									"path_includes": schema.ListAttribute{
										Optional:    true,
										Computed:    true,
										ElementType: types.StringType,
									},
									"pr_comments_enabled": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
									"preview_branch_excludes": schema.ListAttribute{
										Optional:    true,
										Computed:    true,
										ElementType: types.StringType,
									},
									"preview_branch_includes": schema.ListAttribute{
										Optional:    true,
										Computed:    true,
										ElementType: types.StringType,
									},
									"preview_deployment_setting": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"production_branch": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									// v4 field name - renamed to production_deployments_enabled in v5
									"production_deployment_enabled": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
									"repo_name": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"deployment_configs": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"preview":    deploymentConfigSchemaV0(),
						"production": deploymentConfigSchemaV0(),
					},
				},
			},
			// Computed fields
			"created_on": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
			"subdomain": schema.StringAttribute{
				Computed: true,
			},
			"domains": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// deploymentConfigSchemaV0 returns the v4 schema for deployment config (preview/production).
// v4 SDKv2 stores blocks as lists.
func deploymentConfigSchemaV0() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"always_use_latest_compatibility_date": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"compatibility_date": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"compatibility_flags": schema.ListAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"usage_model": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"fail_open": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				// v4: separate maps for env vars and secrets
				"environment_variables": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"secrets": schema.MapAttribute{
					Optional:    true,
					Sensitive:   true,
					ElementType: types.StringType,
				},
				// v4: simple string maps for bindings
				"kv_namespaces": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"d1_databases": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"r2_buckets": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"durable_object_namespaces": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				// v4: service_binding is a list of objects
				"service_binding": schema.ListNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required: true,
							},
							"service": schema.StringAttribute{
								Required: true,
							},
							"environment": schema.StringAttribute{
								Optional: true,
							},
							"entrypoint": schema.StringAttribute{
								Optional: true,
							},
						},
					},
				},
				"placement": schema.ListNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

// =============================================================================
// V0 Model Definitions (v4 provider)
// =============================================================================

// PagesProjectModelV0 represents the v4 state structure.
// Note: v4 SDKv2 stores blocks as lists, so all block types are slices.
type PagesProjectModelV0 struct {
	ID                types.String                           `tfsdk:"id"`
	Name              types.String                           `tfsdk:"name"`
	AccountID         types.String                           `tfsdk:"account_id"`
	ProductionBranch  types.String                           `tfsdk:"production_branch"`
	BuildConfig       []PagesProjectBuildConfigModelV0       `tfsdk:"build_config"`
	Source            []PagesProjectSourceModelV0            `tfsdk:"source"`
	DeploymentConfigs []PagesProjectDeploymentConfigsModelV0 `tfsdk:"deployment_configs"`
	CreatedOn         timetypes.RFC3339                      `tfsdk:"created_on"`
	Subdomain         types.String                           `tfsdk:"subdomain"`
	Domains           types.List                             `tfsdk:"domains"`
}

type PagesProjectBuildConfigModelV0 struct {
	BuildCaching      types.Bool   `tfsdk:"build_caching"`
	BuildCommand      types.String `tfsdk:"build_command"`
	DestinationDir    types.String `tfsdk:"destination_dir"`
	RootDir           types.String `tfsdk:"root_dir"`
	WebAnalyticsTag   types.String `tfsdk:"web_analytics_tag"`
	WebAnalyticsToken types.String `tfsdk:"web_analytics_token"`
}

type PagesProjectSourceModelV0 struct {
	Type   types.String                       `tfsdk:"type"`
	Config []PagesProjectSourceConfigModelV0 `tfsdk:"config"`
}

type PagesProjectSourceConfigModelV0 struct {
	DeploymentsEnabled          types.Bool   `tfsdk:"deployments_enabled"`
	Owner                       types.String `tfsdk:"owner"`
	PathExcludes                types.List   `tfsdk:"path_excludes"`
	PathIncludes                types.List   `tfsdk:"path_includes"`
	PrCommentsEnabled           types.Bool   `tfsdk:"pr_comments_enabled"`
	PreviewBranchExcludes       types.List   `tfsdk:"preview_branch_excludes"`
	PreviewBranchIncludes       types.List   `tfsdk:"preview_branch_includes"`
	PreviewDeploymentSetting    types.String `tfsdk:"preview_deployment_setting"`
	ProductionBranch            types.String `tfsdk:"production_branch"`
	ProductionDeploymentEnabled types.Bool   `tfsdk:"production_deployment_enabled"` // renamed in v5
	RepoName                    types.String `tfsdk:"repo_name"`
}

type PagesProjectDeploymentConfigsModelV0 struct {
	Preview    []PagesProjectDeploymentConfigModelV0 `tfsdk:"preview"`
	Production []PagesProjectDeploymentConfigModelV0 `tfsdk:"production"`
}

type PagesProjectDeploymentConfigModelV0 struct {
	AlwaysUseLatestCompatibilityDate types.Bool                     `tfsdk:"always_use_latest_compatibility_date"`
	CompatibilityDate                types.String                   `tfsdk:"compatibility_date"`
	CompatibilityFlags               types.List                     `tfsdk:"compatibility_flags"`
	UsageModel                       types.String                   `tfsdk:"usage_model"`
	FailOpen                         types.Bool                     `tfsdk:"fail_open"`
	EnvironmentVariables             types.Map                      `tfsdk:"environment_variables"`
	Secrets                          types.Map                      `tfsdk:"secrets"`
	KVNamespaces                     types.Map                      `tfsdk:"kv_namespaces"`
	D1Databases                      types.Map                      `tfsdk:"d1_databases"`
	R2Buckets                        types.Map                      `tfsdk:"r2_buckets"`
	DurableObjectNamespaces          types.Map                      `tfsdk:"durable_object_namespaces"`
	ServiceBindings                  types.List                     `tfsdk:"service_binding"`
	Placement                        []PagesProjectPlacementModelV0 `tfsdk:"placement"`
}

type PagesProjectPlacementModelV0 struct {
	Mode types.String `tfsdk:"mode"`
}

type PagesProjectServiceBindingModelV0 struct {
	Name        types.String `tfsdk:"name"`
	Service     types.String `tfsdk:"service"`
	Environment types.String `tfsdk:"environment"`
	Entrypoint  types.String `tfsdk:"entrypoint"`
}

// =============================================================================
// State Upgrade Function V0 -> V1 (v4 -> v5)
// =============================================================================

// upgradeStateFromV0 upgrades the state from v4 (version 0) to v5 (version 2).
// v4 SDKv2 stores blocks as lists, so we access the first element of each slice.
func upgradeStateFromV0(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var priorState PagesProjectModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new v5 state
	newState := PagesProjectModel{
		ID:               priorState.ID,
		Name:             priorState.Name,
		AccountID:        priorState.AccountID,
		ProductionBranch: priorState.ProductionBranch,
		CreatedOn:        priorState.CreatedOn,
		Subdomain:        priorState.Subdomain,
	}

	// Convert domains list
	if !priorState.Domains.IsNull() && !priorState.Domains.IsUnknown() {
		newState.Domains = customfield.NewListMust[types.String](ctx, priorState.Domains.Elements())
	}

	// Convert build_config (v4 stores as list, take first element)
	if len(priorState.BuildConfig) > 0 {
		bc := priorState.BuildConfig[0]
		newState.BuildConfig = customfield.NewObjectMust(ctx, &PagesProjectBuildConfigModel{
			BuildCaching:      bc.BuildCaching,
			BuildCommand:      bc.BuildCommand,
			DestinationDir:    bc.DestinationDir,
			RootDir:           bc.RootDir,
			WebAnalyticsTag:   bc.WebAnalyticsTag,
			WebAnalyticsToken: bc.WebAnalyticsToken,
		})
	}

	// Convert source (v4 stores as list, take first element, handle field rename)
	if len(priorState.Source) > 0 {
		src := priorState.Source[0]
		newState.Source = &PagesProjectSourceModel{
			Type: src.Type,
		}
		if len(src.Config) > 0 {
			cfg := src.Config[0]
			newState.Source.Config = &PagesProjectSourceConfigModel{
				DeploymentsEnabled:           cfg.DeploymentsEnabled,
				Owner:                        cfg.Owner,
				PrCommentsEnabled:            cfg.PrCommentsEnabled,
				PreviewDeploymentSetting:     cfg.PreviewDeploymentSetting,
				ProductionBranch:             cfg.ProductionBranch,
				ProductionDeploymentsEnabled: cfg.ProductionDeploymentEnabled, // renamed field
				RepoName:                     cfg.RepoName,
			}
			// Convert list fields
			if !cfg.PathExcludes.IsNull() && !cfg.PathExcludes.IsUnknown() {
				newState.Source.Config.PathExcludes = customfield.NewListMust[types.String](ctx, cfg.PathExcludes.Elements())
			}
			if !cfg.PathIncludes.IsNull() && !cfg.PathIncludes.IsUnknown() {
				newState.Source.Config.PathIncludes = customfield.NewListMust[types.String](ctx, cfg.PathIncludes.Elements())
			}
			if !cfg.PreviewBranchExcludes.IsNull() && !cfg.PreviewBranchExcludes.IsUnknown() {
				newState.Source.Config.PreviewBranchExcludes = customfield.NewListMust[types.String](ctx, cfg.PreviewBranchExcludes.Elements())
			}
			if !cfg.PreviewBranchIncludes.IsNull() && !cfg.PreviewBranchIncludes.IsUnknown() {
				newState.Source.Config.PreviewBranchIncludes = customfield.NewListMust[types.String](ctx, cfg.PreviewBranchIncludes.Elements())
			}
		}
	}

	// Convert deployment_configs (v4 stores as list, take first element)
	if len(priorState.DeploymentConfigs) > 0 {
		dc := priorState.DeploymentConfigs[0]
		var deploymentConfigs PagesProjectDeploymentConfigsModel

		if len(dc.Preview) > 0 {
			preview := upgradePreviewConfigV0ToV5(ctx, &dc.Preview[0])
			deploymentConfigs.Preview, _ = customfield.NewObject(ctx, preview)
		}

		if len(dc.Production) > 0 {
			production := upgradeProductionConfigV0ToV5(ctx, &dc.Production[0])
			deploymentConfigs.Production, _ = customfield.NewObject(ctx, production)
		}

		newState.DeploymentConfigs, _ = customfield.NewObject(ctx, &deploymentConfigs)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

// upgradePreviewConfigV0ToV5 converts a v4 preview deployment config to v5 format.
func upgradePreviewConfigV0ToV5(ctx context.Context, v0 *PagesProjectDeploymentConfigModelV0) *PagesProjectDeploymentConfigsPreviewModel {
	if v0 == nil {
		return nil
	}

	v5 := &PagesProjectDeploymentConfigsPreviewModel{
		AlwaysUseLatestCompatibilityDate: v0.AlwaysUseLatestCompatibilityDate,
		CompatibilityDate:                v0.CompatibilityDate,
		// Note: UsageModel is intentionally not copied from v4 state.
		// v4 had "bundled" as default, but v5 API returns "standard" for all projects.
		// Let Read() populate this from the API to avoid plan diffs.
		FailOpen: v0.FailOpen,
	}

	// Convert compatibility_flags list to pointer slice
	if !v0.CompatibilityFlags.IsNull() && !v0.CompatibilityFlags.IsUnknown() {
		flags := make([]types.String, 0, len(v0.CompatibilityFlags.Elements()))
		for _, elem := range v0.CompatibilityFlags.Elements() {
			if str, ok := elem.(types.String); ok {
				flags = append(flags, str)
			}
		}
		v5.CompatibilityFlags = &flags
	}

	// Convert placement (v4 stores as list, take first element)
	if len(v0.Placement) > 0 && !v0.Placement[0].Mode.IsNull() {
		v5.Placement = &PagesProjectDeploymentConfigsPreviewPlacementModel{
			Mode: v0.Placement[0].Mode,
		}
	}

	// Merge environment_variables and secrets into env_vars
	v5.EnvVars = mergeEnvVarsAndSecretsPreview(ctx, v0.EnvironmentVariables, v0.Secrets)

	// Convert kv_namespaces
	v5.KVNamespaces = convertKVNamespacesV0ToV5Preview(ctx, v0.KVNamespaces)

	// Convert d1_databases
	v5.D1Databases = convertD1DatabasesV0ToV5Preview(ctx, v0.D1Databases)

	// Convert r2_buckets
	v5.R2Buckets = convertR2BucketsV0ToV5Preview(ctx, v0.R2Buckets)

	// Convert durable_object_namespaces
	v5.DurableObjectNamespaces = convertDurableObjectNamespacesV0ToV5Preview(ctx, v0.DurableObjectNamespaces)

	// Convert service_binding list -> services map
	v5.Services = convertServiceBindingsV0ToV5Preview(ctx, v0.ServiceBindings)

	return v5
}

// upgradeProductionConfigV0ToV5 converts a v4 production deployment config to v5 format.
func upgradeProductionConfigV0ToV5(ctx context.Context, v0 *PagesProjectDeploymentConfigModelV0) *PagesProjectDeploymentConfigsProductionModel {
	if v0 == nil {
		return nil
	}

	v5 := &PagesProjectDeploymentConfigsProductionModel{
		AlwaysUseLatestCompatibilityDate: v0.AlwaysUseLatestCompatibilityDate,
		CompatibilityDate:                v0.CompatibilityDate,
		// Note: UsageModel is intentionally not copied from v4 state.
		// v4 had "bundled" as default, but v5 API returns "standard" for all projects.
		// Let Read() populate this from the API to avoid plan diffs.
		FailOpen: v0.FailOpen,
	}

	// Convert compatibility_flags list to pointer slice
	if !v0.CompatibilityFlags.IsNull() && !v0.CompatibilityFlags.IsUnknown() {
		flags := make([]types.String, 0, len(v0.CompatibilityFlags.Elements()))
		for _, elem := range v0.CompatibilityFlags.Elements() {
			if str, ok := elem.(types.String); ok {
				flags = append(flags, str)
			}
		}
		v5.CompatibilityFlags = &flags
	}

	// Convert placement (v4 stores as list, take first element)
	if len(v0.Placement) > 0 && !v0.Placement[0].Mode.IsNull() {
		v5.Placement = &PagesProjectDeploymentConfigsProductionPlacementModel{
			Mode: v0.Placement[0].Mode,
		}
	}

	// Merge environment_variables and secrets into env_vars
	v5.EnvVars = mergeEnvVarsAndSecretsProduction(ctx, v0.EnvironmentVariables, v0.Secrets)

	// Convert kv_namespaces
	v5.KVNamespaces = convertKVNamespacesV0ToV5Production(ctx, v0.KVNamespaces)

	// Convert d1_databases
	v5.D1Databases = convertD1DatabasesV0ToV5Production(ctx, v0.D1Databases)

	// Convert r2_buckets
	v5.R2Buckets = convertR2BucketsV0ToV5Production(ctx, v0.R2Buckets)

	// Convert durable_object_namespaces
	v5.DurableObjectNamespaces = convertDurableObjectNamespacesV0ToV5Production(ctx, v0.DurableObjectNamespaces)

	// Convert service_binding list -> services map
	v5.Services = convertServiceBindingsV0ToV5Production(ctx, v0.ServiceBindings)

	return v5
}

// =============================================================================
// Preview Config Conversion Helpers
// =============================================================================

// mergeEnvVarsAndSecretsPreview merges v4 environment_variables and secrets into v5 env_vars format for preview.
func mergeEnvVarsAndSecretsPreview(ctx context.Context, envVars, secrets types.Map) *map[string]PagesProjectDeploymentConfigsPreviewEnvVarsModel {
	result := make(map[string]PagesProjectDeploymentConfigsPreviewEnvVarsModel)

	// Add environment_variables as plain_text
	if !envVars.IsNull() && !envVars.IsUnknown() {
		for key, val := range envVars.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = PagesProjectDeploymentConfigsPreviewEnvVarsModel{
					Type:  types.StringValue("plain_text"),
					Value: strVal,
				}
			}
		}
	}

	// Add secrets as secret_text
	if !secrets.IsNull() && !secrets.IsUnknown() {
		for key, val := range secrets.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = PagesProjectDeploymentConfigsPreviewEnvVarsModel{
					Type:  types.StringValue("secret_text"),
					Value: strVal,
				}
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertKVNamespacesV0ToV5Preview converts v4 kv_namespaces (map[string]string) to v5 preview format.
func convertKVNamespacesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewKVNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsPreviewKVNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsPreviewKVNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertD1DatabasesV0ToV5Preview converts v4 d1_databases (map[string]string) to v5 preview format.
func convertD1DatabasesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewD1DatabasesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsPreviewD1DatabasesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsPreviewD1DatabasesModel{
				ID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertR2BucketsV0ToV5Preview converts v4 r2_buckets (map[string]string) to v5 preview format.
func convertR2BucketsV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewR2BucketsModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsPreviewR2BucketsModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsPreviewR2BucketsModel{
				Name: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertDurableObjectNamespacesV0ToV5Preview converts v4 durable_object_namespaces (map[string]string) to v5 preview format.
func convertDurableObjectNamespacesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertServiceBindingsV0ToV5Preview converts v4 service_binding list to v5 preview services map.
func convertServiceBindingsV0ToV5Preview(ctx context.Context, v0 types.List) *map[string]PagesProjectDeploymentConfigsPreviewServicesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsPreviewServicesModel)

	var bindings []PagesProjectServiceBindingModelV0
	if diags := v0.ElementsAs(ctx, &bindings, false); diags.HasError() {
		return nil
	}

	for _, binding := range bindings {
		if !binding.Name.IsNull() && !binding.Name.IsUnknown() {
			result[binding.Name.ValueString()] = PagesProjectDeploymentConfigsPreviewServicesModel{
				Service:     binding.Service,
				Environment: binding.Environment,
				Entrypoint:  binding.Entrypoint,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// =============================================================================
// Production Config Conversion Helpers
// =============================================================================

// mergeEnvVarsAndSecretsProduction merges v4 environment_variables and secrets into v5 env_vars format for production.
func mergeEnvVarsAndSecretsProduction(ctx context.Context, envVars, secrets types.Map) *map[string]PagesProjectDeploymentConfigsProductionEnvVarsModel {
	result := make(map[string]PagesProjectDeploymentConfigsProductionEnvVarsModel)

	// Add environment_variables as plain_text
	if !envVars.IsNull() && !envVars.IsUnknown() {
		for key, val := range envVars.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = PagesProjectDeploymentConfigsProductionEnvVarsModel{
					Type:  types.StringValue("plain_text"),
					Value: strVal,
				}
			}
		}
	}

	// Add secrets as secret_text
	if !secrets.IsNull() && !secrets.IsUnknown() {
		for key, val := range secrets.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = PagesProjectDeploymentConfigsProductionEnvVarsModel{
					Type:  types.StringValue("secret_text"),
					Value: strVal,
				}
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertKVNamespacesV0ToV5Production converts v4 kv_namespaces (map[string]string) to v5 production format.
func convertKVNamespacesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsProductionKVNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsProductionKVNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsProductionKVNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertD1DatabasesV0ToV5Production converts v4 d1_databases (map[string]string) to v5 production format.
func convertD1DatabasesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsProductionD1DatabasesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsProductionD1DatabasesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsProductionD1DatabasesModel{
				ID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertR2BucketsV0ToV5Production converts v4 r2_buckets (map[string]string) to v5 production format.
func convertR2BucketsV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsProductionR2BucketsModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsProductionR2BucketsModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsProductionR2BucketsModel{
				Name: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertDurableObjectNamespacesV0ToV5Production converts v4 durable_object_namespaces (map[string]string) to v5 production format.
func convertDurableObjectNamespacesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = PagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// convertServiceBindingsV0ToV5Production converts v4 service_binding list to v5 production services map.
func convertServiceBindingsV0ToV5Production(ctx context.Context, v0 types.List) *map[string]PagesProjectDeploymentConfigsProductionServicesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]PagesProjectDeploymentConfigsProductionServicesModel)

	var bindings []PagesProjectServiceBindingModelV0
	if diags := v0.ElementsAs(ctx, &bindings, false); diags.HasError() {
		return nil
	}

	for _, binding := range bindings {
		if !binding.Name.IsNull() && !binding.Name.IsUnknown() {
			result[binding.Name.ValueString()] = PagesProjectDeploymentConfigsProductionServicesModel{
				Service:     binding.Service,
				Environment: binding.Environment,
				Entrypoint:  binding.Entrypoint,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// =============================================================================
// Exported Functions for Testing
// =============================================================================

// MergeEnvVarsAndSecretsPreviewExported is an exported wrapper for testing.
func MergeEnvVarsAndSecretsPreviewExported(ctx context.Context, envVars, secrets types.Map) *map[string]PagesProjectDeploymentConfigsPreviewEnvVarsModel {
	return mergeEnvVarsAndSecretsPreview(ctx, envVars, secrets)
}

// ConvertKVNamespacesV0ToV5PreviewExported is an exported wrapper for testing.
func ConvertKVNamespacesV0ToV5PreviewExported(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewKVNamespacesModel {
	return convertKVNamespacesV0ToV5Preview(ctx, v0)
}

// ConvertD1DatabasesV0ToV5PreviewExported is an exported wrapper for testing.
func ConvertD1DatabasesV0ToV5PreviewExported(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewD1DatabasesModel {
	return convertD1DatabasesV0ToV5Preview(ctx, v0)
}

// ConvertR2BucketsV0ToV5PreviewExported is an exported wrapper for testing.
func ConvertR2BucketsV0ToV5PreviewExported(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewR2BucketsModel {
	return convertR2BucketsV0ToV5Preview(ctx, v0)
}

// ConvertDurableObjectNamespacesV0ToV5PreviewExported is an exported wrapper for testing.
func ConvertDurableObjectNamespacesV0ToV5PreviewExported(ctx context.Context, v0 types.Map) *map[string]PagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel {
	return convertDurableObjectNamespacesV0ToV5Preview(ctx, v0)
}
