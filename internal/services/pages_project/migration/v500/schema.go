package v500

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SourcePagesProjectSchemaV0 returns the v4 provider schema for pages_project.
// Key differences from v5:
//   - Blocks are represented as lists (SDKv2 style) not single nested attributes
//   - environment_variables and secrets are separate maps (merged to env_vars in v5)
//   - service_binding is a list (converted to services map in v5)
//   - kv_namespaces, d1_databases, r2_buckets, durable_object_namespaces are simple string maps
//     (wrapped in objects in v5)
//   - production_deployment_enabled renamed to production_deployments_enabled
func SourcePagesProjectSchemaV0(ctx context.Context) *schema.Schema {
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
