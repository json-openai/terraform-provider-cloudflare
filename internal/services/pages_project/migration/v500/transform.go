package v500

import (
	"context"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Transform converts v0 (v4 provider) state to v5 state.
func Transform(ctx context.Context, source SourcePagesProjectModelV0) (*TargetPagesProjectModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create new v5 state
	target := &TargetPagesProjectModel{
		ID:               source.ID,
		Name:             source.Name,
		AccountID:        source.AccountID,
		ProductionBranch: source.ProductionBranch,
		CreatedOn:        source.CreatedOn,
		Subdomain:        source.Subdomain,
	}

	// Convert domains list
	if !source.Domains.IsNull() && !source.Domains.IsUnknown() {
		target.Domains = customfield.NewListMust[types.String](ctx, source.Domains.Elements())
	}

	// Convert build_config (v4 stores as list, take first element)
	if len(source.BuildConfig) > 0 {
		bc := source.BuildConfig[0]
		target.BuildConfig = customfield.NewObjectMust(ctx, &TargetPagesProjectBuildConfigModel{
			BuildCaching:      bc.BuildCaching,
			BuildCommand:      bc.BuildCommand,
			DestinationDir:    bc.DestinationDir,
			RootDir:           bc.RootDir,
			WebAnalyticsTag:   bc.WebAnalyticsTag,
			WebAnalyticsToken: bc.WebAnalyticsToken,
		})
	}

	// Convert source (v4 stores as list, take first element, handle field rename)
	if len(source.Source) > 0 {
		src := source.Source[0]
		target.Source = &TargetPagesProjectSourceModel{
			Type: src.Type,
		}
		if len(src.Config) > 0 {
			cfg := src.Config[0]
			target.Source.Config = &TargetPagesProjectSourceConfigModel{
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
				target.Source.Config.PathExcludes = customfield.NewListMust[types.String](ctx, cfg.PathExcludes.Elements())
			}
			if !cfg.PathIncludes.IsNull() && !cfg.PathIncludes.IsUnknown() {
				target.Source.Config.PathIncludes = customfield.NewListMust[types.String](ctx, cfg.PathIncludes.Elements())
			}
			if !cfg.PreviewBranchExcludes.IsNull() && !cfg.PreviewBranchExcludes.IsUnknown() {
				target.Source.Config.PreviewBranchExcludes = customfield.NewListMust[types.String](ctx, cfg.PreviewBranchExcludes.Elements())
			}
			if !cfg.PreviewBranchIncludes.IsNull() && !cfg.PreviewBranchIncludes.IsUnknown() {
				target.Source.Config.PreviewBranchIncludes = customfield.NewListMust[types.String](ctx, cfg.PreviewBranchIncludes.Elements())
			}
		}
	}

	// Convert deployment_configs (v4 stores as list, take first element)
	if len(source.DeploymentConfigs) > 0 {
		dc := source.DeploymentConfigs[0]
		var deploymentConfigs TargetPagesProjectDeploymentConfigsModel

		if len(dc.Preview) > 0 {
			preview := transformPreviewConfigV0ToV5(ctx, &dc.Preview[0])
			deploymentConfigs.Preview, _ = customfield.NewObject(ctx, preview)
		}

		if len(dc.Production) > 0 {
			production := transformProductionConfigV0ToV5(ctx, &dc.Production[0])
			deploymentConfigs.Production, _ = customfield.NewObject(ctx, production)
		}

		target.DeploymentConfigs, _ = customfield.NewObject(ctx, &deploymentConfigs)
	}

	return target, diags
}

// transformPreviewConfigV0ToV5 converts a v4 preview deployment config to v5 format.
func transformPreviewConfigV0ToV5(ctx context.Context, v0 *SourcePagesProjectDeploymentConfigModelV0) *TargetPagesProjectDeploymentConfigsPreviewModel {
	if v0 == nil {
		return nil
	}

	v5 := &TargetPagesProjectDeploymentConfigsPreviewModel{
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
		v5.Placement = &TargetPagesProjectDeploymentConfigsPreviewPlacementModel{
			Mode: v0.Placement[0].Mode,
		}
	}

	// Merge environment_variables and secrets into env_vars
	v5.EnvVars = MergeEnvVarsAndSecretsPreview(ctx, v0.EnvironmentVariables, v0.Secrets)

	// Convert kv_namespaces
	v5.KVNamespaces = ConvertKVNamespacesV0ToV5Preview(ctx, v0.KVNamespaces)

	// Convert d1_databases
	v5.D1Databases = ConvertD1DatabasesV0ToV5Preview(ctx, v0.D1Databases)

	// Convert r2_buckets
	v5.R2Buckets = ConvertR2BucketsV0ToV5Preview(ctx, v0.R2Buckets)

	// Convert durable_object_namespaces
	v5.DurableObjectNamespaces = ConvertDurableObjectNamespacesV0ToV5Preview(ctx, v0.DurableObjectNamespaces)

	// Convert service_binding list -> services map
	v5.Services = ConvertServiceBindingsV0ToV5Preview(ctx, v0.ServiceBindings)

	return v5
}

// transformProductionConfigV0ToV5 converts a v4 production deployment config to v5 format.
func transformProductionConfigV0ToV5(ctx context.Context, v0 *SourcePagesProjectDeploymentConfigModelV0) *TargetPagesProjectDeploymentConfigsProductionModel {
	if v0 == nil {
		return nil
	}

	v5 := &TargetPagesProjectDeploymentConfigsProductionModel{
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
		v5.Placement = &TargetPagesProjectDeploymentConfigsProductionPlacementModel{
			Mode: v0.Placement[0].Mode,
		}
	}

	// Merge environment_variables and secrets into env_vars
	v5.EnvVars = MergeEnvVarsAndSecretsProduction(ctx, v0.EnvironmentVariables, v0.Secrets)

	// Convert kv_namespaces
	v5.KVNamespaces = ConvertKVNamespacesV0ToV5Production(ctx, v0.KVNamespaces)

	// Convert d1_databases
	v5.D1Databases = ConvertD1DatabasesV0ToV5Production(ctx, v0.D1Databases)

	// Convert r2_buckets
	v5.R2Buckets = ConvertR2BucketsV0ToV5Production(ctx, v0.R2Buckets)

	// Convert durable_object_namespaces
	v5.DurableObjectNamespaces = ConvertDurableObjectNamespacesV0ToV5Production(ctx, v0.DurableObjectNamespaces)

	// Convert service_binding list -> services map
	v5.Services = ConvertServiceBindingsV0ToV5Production(ctx, v0.ServiceBindings)

	return v5
}

// =============================================================================
// Preview Config Conversion Helpers
// =============================================================================

// MergeEnvVarsAndSecretsPreview merges v4 environment_variables and secrets into v5 env_vars format for preview.
func MergeEnvVarsAndSecretsPreview(ctx context.Context, envVars, secrets types.Map) *map[string]TargetPagesProjectDeploymentConfigsPreviewEnvVarsModel {
	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewEnvVarsModel)

	// Add environment_variables as plain_text
	if !envVars.IsNull() && !envVars.IsUnknown() {
		for key, val := range envVars.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = TargetPagesProjectDeploymentConfigsPreviewEnvVarsModel{
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
				result[key] = TargetPagesProjectDeploymentConfigsPreviewEnvVarsModel{
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

// ConvertKVNamespacesV0ToV5Preview converts v4 kv_namespaces (map[string]string) to v5 preview format.
func ConvertKVNamespacesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsPreviewKVNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewKVNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsPreviewKVNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertD1DatabasesV0ToV5Preview converts v4 d1_databases (map[string]string) to v5 preview format.
func ConvertD1DatabasesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsPreviewD1DatabasesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewD1DatabasesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsPreviewD1DatabasesModel{
				ID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertR2BucketsV0ToV5Preview converts v4 r2_buckets (map[string]string) to v5 preview format.
func ConvertR2BucketsV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsPreviewR2BucketsModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewR2BucketsModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsPreviewR2BucketsModel{
				Name: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertDurableObjectNamespacesV0ToV5Preview converts v4 durable_object_namespaces (map[string]string) to v5 preview format.
func ConvertDurableObjectNamespacesV0ToV5Preview(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsPreviewDurableObjectNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertServiceBindingsV0ToV5Preview converts v4 service_binding list to v5 preview services map.
func ConvertServiceBindingsV0ToV5Preview(ctx context.Context, v0 types.List) *map[string]TargetPagesProjectDeploymentConfigsPreviewServicesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsPreviewServicesModel)

	var bindings []SourcePagesProjectServiceBindingModelV0
	if diags := v0.ElementsAs(ctx, &bindings, false); diags.HasError() {
		return nil
	}

	for _, binding := range bindings {
		if !binding.Name.IsNull() && !binding.Name.IsUnknown() {
			result[binding.Name.ValueString()] = TargetPagesProjectDeploymentConfigsPreviewServicesModel{
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

// MergeEnvVarsAndSecretsProduction merges v4 environment_variables and secrets into v5 env_vars format for production.
func MergeEnvVarsAndSecretsProduction(ctx context.Context, envVars, secrets types.Map) *map[string]TargetPagesProjectDeploymentConfigsProductionEnvVarsModel {
	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionEnvVarsModel)

	// Add environment_variables as plain_text
	if !envVars.IsNull() && !envVars.IsUnknown() {
		for key, val := range envVars.Elements() {
			if strVal, ok := val.(types.String); ok {
				result[key] = TargetPagesProjectDeploymentConfigsProductionEnvVarsModel{
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
				result[key] = TargetPagesProjectDeploymentConfigsProductionEnvVarsModel{
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

// ConvertKVNamespacesV0ToV5Production converts v4 kv_namespaces (map[string]string) to v5 production format.
func ConvertKVNamespacesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsProductionKVNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionKVNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsProductionKVNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertD1DatabasesV0ToV5Production converts v4 d1_databases (map[string]string) to v5 production format.
func ConvertD1DatabasesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsProductionD1DatabasesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionD1DatabasesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsProductionD1DatabasesModel{
				ID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertR2BucketsV0ToV5Production converts v4 r2_buckets (map[string]string) to v5 production format.
func ConvertR2BucketsV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsProductionR2BucketsModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionR2BucketsModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsProductionR2BucketsModel{
				Name: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertDurableObjectNamespacesV0ToV5Production converts v4 durable_object_namespaces (map[string]string) to v5 production format.
func ConvertDurableObjectNamespacesV0ToV5Production(ctx context.Context, v0 types.Map) *map[string]TargetPagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel)
	for key, val := range v0.Elements() {
		if strVal, ok := val.(types.String); ok {
			result[key] = TargetPagesProjectDeploymentConfigsProductionDurableObjectNamespacesModel{
				NamespaceID: strVal,
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return &result
}

// ConvertServiceBindingsV0ToV5Production converts v4 service_binding list to v5 production services map.
func ConvertServiceBindingsV0ToV5Production(ctx context.Context, v0 types.List) *map[string]TargetPagesProjectDeploymentConfigsProductionServicesModel {
	if v0.IsNull() || v0.IsUnknown() {
		return nil
	}

	result := make(map[string]TargetPagesProjectDeploymentConfigsProductionServicesModel)

	var bindings []SourcePagesProjectServiceBindingModelV0
	if diags := v0.ElementsAs(ctx, &bindings, false); diags.HasError() {
		return nil
	}

	for _, binding := range bindings {
		if !binding.Name.IsNull() && !binding.Name.IsUnknown() {
			result[binding.Name.ValueString()] = TargetPagesProjectDeploymentConfigsProductionServicesModel{
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
