package v500

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ============================================================================
// Source V4 Plain Go structs for JSON unmarshaling (v4 SDKv2 format)
// ============================================================================

// SourceV4PageRuleJSON is the plain Go struct for v4 state JSON unmarshaling.
type SourceV4PageRuleJSON struct {
	ID       string          `json:"id"`
	ZoneID   string          `json:"zone_id"`
	Target   string          `json:"target"`
	Priority *float64        `json:"priority"` // JSON numbers are float64
	Status   string          `json:"status"`
	Actions  []SourceV4ActionsJSON `json:"actions"`
}

// SourceV4ActionsJSON is the plain Go struct for v4 actions array element.
type SourceV4ActionsJSON struct {
	AlwaysUseHTTPS          *bool    `json:"always_use_https"`
	DisableApps             *bool    `json:"disable_apps"`
	DisablePerformance      *bool    `json:"disable_performance"`
	DisableRailgun          *bool    `json:"disable_railgun"`
	DisableSecurity         *bool    `json:"disable_security"`
	DisableZaraz            *bool    `json:"disable_zaraz"`
	AutomaticHTTPSRewrites  *string  `json:"automatic_https_rewrites"`
	BrowserCacheTTL         *string  `json:"browser_cache_ttl"` // String in v4!
	BrowserCheck            *string  `json:"browser_check"`
	BypassCacheOnCookie     *string  `json:"bypass_cache_on_cookie"`
	CacheByDeviceType       *string  `json:"cache_by_device_type"`
	CacheDeceptionArmor     *string  `json:"cache_deception_armor"`
	CacheLevel              *string  `json:"cache_level"`
	CacheOnCookie           *string  `json:"cache_on_cookie"`
	EmailObfuscation        *string  `json:"email_obfuscation"`
	ExplicitCacheControl    *string  `json:"explicit_cache_control"`
	HostHeaderOverride      *string  `json:"host_header_override"`
	IPGeolocation           *string  `json:"ip_geolocation"`
	Mirage                  *string  `json:"mirage"`
	OpportunisticEncryption *string  `json:"opportunistic_encryption"`
	OriginErrorPagePassThru *string  `json:"origin_error_page_pass_thru"`
	Polish                  *string  `json:"polish"`
	ResolveOverride         *string  `json:"resolve_override"`
	RespectStrongEtag       *string  `json:"respect_strong_etag"`
	ResponseBuffering       *string  `json:"response_buffering"`
	RocketLoader            *string  `json:"rocket_loader"`
	SecurityLevel           *string  `json:"security_level"`
	ServerSideExclude       *string  `json:"server_side_exclude"`
	SortQueryStringForCache *string  `json:"sort_query_string_for_cache"`
	SSL                     *string  `json:"ssl"`
	TrueClientIPHeader      *string  `json:"true_client_ip_header"`
	WAF                     *string  `json:"waf"`
	EdgeCacheTTL            *float64 `json:"edge_cache_ttl"`

	ForwardingURL    []SourceV4ForwardingURLJSON    `json:"forwarding_url"`
	Minify           []SourceV4MinifyJSON           `json:"minify"`
	CacheKeyFields   []SourceV4CacheKeyFieldsJSON   `json:"cache_key_fields"`
	CacheTTLByStatus []SourceV4CacheTTLByStatusJSON `json:"cache_ttl_by_status"`
}

type SourceV4ForwardingURLJSON struct {
	URL        string   `json:"url"`
	StatusCode *float64 `json:"status_code"`
}

type SourceV4MinifyJSON struct {
	JS   string `json:"js"`
	CSS  string `json:"css"`
	HTML string `json:"html"`
}

type SourceV4CacheKeyFieldsJSON struct {
	Cookie      []SourceV4CacheKeyFieldsCookieJSON      `json:"cookie"`
	Header      []SourceV4CacheKeyFieldsHeaderJSON      `json:"header"`
	Host        []SourceV4CacheKeyFieldsHostJSON        `json:"host"`
	QueryString []SourceV4CacheKeyFieldsQueryStringJSON `json:"query_string"`
	User        []SourceV4CacheKeyFieldsUserJSON        `json:"user"`
}

type SourceV4CacheKeyFieldsCookieJSON struct {
	CheckPresence []string `json:"check_presence"`
	Include       []string `json:"include"`
}

type SourceV4CacheKeyFieldsHeaderJSON struct {
	CheckPresence []string `json:"check_presence"`
	Include       []string `json:"include"`
	Exclude       []string `json:"exclude"`
}

type SourceV4CacheKeyFieldsHostJSON struct {
	Resolved *bool `json:"resolved"`
}

type SourceV4CacheKeyFieldsQueryStringJSON struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Ignore  *bool    `json:"ignore"`
}

type SourceV4CacheKeyFieldsUserJSON struct {
	DeviceType *bool `json:"device_type"`
	Geo        *bool `json:"geo"`
	Lang       *bool `json:"lang"`
}

type SourceV4CacheTTLByStatusJSON struct {
	Codes string   `json:"codes"`
	TTL   *float64 `json:"ttl"`
}

// ============================================================================
// Source V5 Plain Go structs for JSON unmarshaling (v5 Plugin Framework format)
// ============================================================================

type SourceV5PageRuleJSON struct {
	ID       string          `json:"id"`
	ZoneID   string          `json:"zone_id"`
	Target   string          `json:"target"`
	Priority *float64        `json:"priority"`
	Status   string          `json:"status"`
	Actions  *SourceV5ActionsJSON  `json:"actions"` // Object in v5, not array
}

type SourceV5ActionsJSON struct {
	AlwaysUseHTTPS          *bool             `json:"always_use_https"`
	DisableApps             *bool             `json:"disable_apps"`
	DisablePerformance      *bool             `json:"disable_performance"`
	DisableSecurity         *bool             `json:"disable_security"`
	DisableZaraz            *bool             `json:"disable_zaraz"`
	AutomaticHTTPSRewrites  *string           `json:"automatic_https_rewrites"`
	BrowserCacheTTL         *float64          `json:"browser_cache_ttl"` // Int64 in v5
	BrowserCheck            *string           `json:"browser_check"`
	BypassCacheOnCookie     *string           `json:"bypass_cache_on_cookie"`
	CacheByDeviceType       *string           `json:"cache_by_device_type"`
	CacheDeceptionArmor     *string           `json:"cache_deception_armor"`
	CacheLevel              *string           `json:"cache_level"`
	CacheOnCookie           *string           `json:"cache_on_cookie"`
	EdgeCacheTTL            *float64          `json:"edge_cache_ttl"`
	EmailObfuscation        *string           `json:"email_obfuscation"`
	ExplicitCacheControl    *string           `json:"explicit_cache_control"`
	HostHeaderOverride      *string           `json:"host_header_override"`
	IPGeolocation           *string           `json:"ip_geolocation"`
	Mirage                  *string           `json:"mirage"`
	OpportunisticEncryption *string           `json:"opportunistic_encryption"`
	OriginErrorPagePassThru *string           `json:"origin_error_page_pass_thru"`
	Polish                  *string           `json:"polish"`
	ResolveOverride         *string           `json:"resolve_override"`
	RespectStrongEtag       *string           `json:"respect_strong_etag"`
	ResponseBuffering       *string           `json:"response_buffering"`
	RocketLoader            *string           `json:"rocket_loader"`
	SecurityLevel           *string           `json:"security_level"`
	SortQueryStringForCache *string           `json:"sort_query_string_for_cache"`
	SSL                     *string           `json:"ssl"`
	TrueClientIPHeader      *string           `json:"true_client_ip_header"`
	WAF                     *string           `json:"waf"`
	ForwardingURL           *SourceV5ForwardingURLJSON  `json:"forwarding_url"`
	CacheKeyFields          *SourceV5CacheKeyFieldsJSON `json:"cache_key_fields"`
	CacheTTLByStatus        map[string]string     `json:"cache_ttl_by_status"`
}

type SourceV5ForwardingURLJSON struct {
	URL        string   `json:"url"`
	StatusCode *float64 `json:"status_code"`
}

type SourceV5CacheKeyFieldsJSON struct {
	Cookie      *SourceV5CacheKeyFieldsCookieJSON      `json:"cookie"`
	Header      *SourceV5CacheKeyFieldsHeaderJSON      `json:"header"`
	Host        *SourceV5CacheKeyFieldsHostJSON        `json:"host"`
	QueryString *SourceV5CacheKeyFieldsQueryStringJSON `json:"query_string"`
	User        *SourceV5CacheKeyFieldsUserJSON        `json:"user"`
}

type SourceV5CacheKeyFieldsCookieJSON struct {
	CheckPresence []string `json:"check_presence"`
	Include       []string `json:"include"`
}

type SourceV5CacheKeyFieldsHeaderJSON struct {
	CheckPresence []string `json:"check_presence"`
	Include       []string `json:"include"`
	Exclude       []string `json:"exclude"`
}

type SourceV5CacheKeyFieldsHostJSON struct {
	Resolved *bool `json:"resolved"`
}

type SourceV5CacheKeyFieldsQueryStringJSON struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type SourceV5CacheKeyFieldsUserJSON struct {
	DeviceType *bool `json:"device_type"`
	Geo        *bool `json:"geo"`
	Lang       *bool `json:"lang"`
}

// ============================================================================
// Parsing functions
// ============================================================================

// parseV4JSONToModel parses raw JSON (v4 format) into SourceV4PageRuleModel.
func parseV4JSONToModel(ctx context.Context, rawJSON []byte) (SourceV4PageRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var jsonData SourceV4PageRuleJSON

	if err := json.Unmarshal(rawJSON, &jsonData); err != nil {
		diags.AddError("Failed to unmarshal v4 JSON", err.Error())
		return SourceV4PageRuleModel{}, diags
	}

	return convertV4JSONToModel(ctx, jsonData), diags
}

// parseV5JSONToModel parses raw JSON (v5 format) into TargetV5PageRuleModel.
func parseV5JSONToModel(ctx context.Context, rawJSON []byte) (*TargetV5PageRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var jsonData SourceV5PageRuleJSON

	if err := json.Unmarshal(rawJSON, &jsonData); err != nil {
		diags.AddError("Failed to unmarshal v5 JSON", err.Error())
		return nil, diags
	}

	return convertV5JSONToModel(ctx, jsonData), diags
}

// ============================================================================
// Conversion: v4 JSON -> v4 Terraform Model
// ============================================================================

func convertV4JSONToModel(ctx context.Context, j SourceV4PageRuleJSON) SourceV4PageRuleModel {
	model := SourceV4PageRuleModel{
		ID:       toStringValue(j.ID),
		ZoneID:   toStringValue(j.ZoneID),
		Target:   toStringValue(j.Target),
		Priority: toInt64FromFloat(j.Priority),
		Status:   toStringValue(j.Status),
	}

	if len(j.Actions) > 0 {
		model.Actions = []SourceV4ActionsModel{convertSourceV4ActionsJSON(ctx, j.Actions[0])}
	}

	return model
}

func convertSourceV4ActionsJSON(ctx context.Context, j SourceV4ActionsJSON) SourceV4ActionsModel {
	model := SourceV4ActionsModel{
		AlwaysUseHTTPS:          toBoolValue(j.AlwaysUseHTTPS),
		DisableApps:             toBoolValue(j.DisableApps),
		DisablePerformance:      toBoolValue(j.DisablePerformance),
		DisableRailgun:          toBoolValue(j.DisableRailgun),
		DisableSecurity:         toBoolValue(j.DisableSecurity),
		DisableZaraz:            toBoolValue(j.DisableZaraz),
		AutomaticHTTPSRewrites:  toStringPtrValue(j.AutomaticHTTPSRewrites),
		BrowserCacheTTL:         toStringPtrValue(j.BrowserCacheTTL),
		BrowserCheck:            toStringPtrValue(j.BrowserCheck),
		BypassCacheOnCookie:     toStringPtrValue(j.BypassCacheOnCookie),
		CacheByDeviceType:       toStringPtrValue(j.CacheByDeviceType),
		CacheDeceptionArmor:     toStringPtrValue(j.CacheDeceptionArmor),
		CacheLevel:              toStringPtrValue(j.CacheLevel),
		CacheOnCookie:           toStringPtrValue(j.CacheOnCookie),
		EmailObfuscation:        toStringPtrValue(j.EmailObfuscation),
		ExplicitCacheControl:    toStringPtrValue(j.ExplicitCacheControl),
		HostHeaderOverride:      toStringPtrValue(j.HostHeaderOverride),
		IPGeolocation:           toStringPtrValue(j.IPGeolocation),
		Mirage:                  toStringPtrValue(j.Mirage),
		OpportunisticEncryption: toStringPtrValue(j.OpportunisticEncryption),
		OriginErrorPagePassThru: toStringPtrValue(j.OriginErrorPagePassThru),
		Polish:                  toStringPtrValue(j.Polish),
		ResolveOverride:         toStringPtrValue(j.ResolveOverride),
		RespectStrongEtag:       toStringPtrValue(j.RespectStrongEtag),
		ResponseBuffering:       toStringPtrValue(j.ResponseBuffering),
		RocketLoader:            toStringPtrValue(j.RocketLoader),
		SecurityLevel:           toStringPtrValue(j.SecurityLevel),
		ServerSideExclude:       toStringPtrValue(j.ServerSideExclude),
		SortQueryStringForCache: toStringPtrValue(j.SortQueryStringForCache),
		SSL:                     toStringPtrValue(j.SSL),
		TrueClientIPHeader:      toStringPtrValue(j.TrueClientIPHeader),
		WAF:                     toStringPtrValue(j.WAF),
		EdgeCacheTTL:            toInt64FromFloat(j.EdgeCacheTTL),
	}

	if len(j.ForwardingURL) > 0 {
		model.ForwardingURL = []SourceV4ForwardingURLModel{{
			URL:        toStringValue(j.ForwardingURL[0].URL),
			StatusCode: toInt64FromFloat(j.ForwardingURL[0].StatusCode),
		}}
	}

	if len(j.Minify) > 0 {
		model.Minify = []SourceV4MinifyModel{{
			JS:   toStringValue(j.Minify[0].JS),
			CSS:  toStringValue(j.Minify[0].CSS),
			HTML: toStringValue(j.Minify[0].HTML),
		}}
	}

	if len(j.CacheKeyFields) > 0 {
		model.CacheKeyFields = []SourceV4CacheKeyFieldsModel{convertSourceV4CacheKeyFieldsJSON(ctx, j.CacheKeyFields[0])}
	}

	for _, entry := range j.CacheTTLByStatus {
		model.CacheTTLByStatus = append(model.CacheTTLByStatus, SourceV4CacheTTLByStatusModel{
			Codes: toStringValue(entry.Codes),
			TTL:   toInt64FromFloat(entry.TTL),
		})
	}

	return model
}

func convertSourceV4CacheKeyFieldsJSON(ctx context.Context, j SourceV4CacheKeyFieldsJSON) SourceV4CacheKeyFieldsModel {
	model := SourceV4CacheKeyFieldsModel{}

	if len(j.Cookie) > 0 {
		model.Cookie = []SourceV4CacheKeyFieldsCookieModel{{
			CheckPresence: toSetValue(ctx, j.Cookie[0].CheckPresence),
			Include:       toSetValue(ctx, j.Cookie[0].Include),
		}}
	}

	if len(j.Header) > 0 {
		model.Header = []SourceV4CacheKeyFieldsHeaderModel{{
			CheckPresence: toSetValue(ctx, j.Header[0].CheckPresence),
			Include:       toSetValue(ctx, j.Header[0].Include),
			Exclude:       toSetValue(ctx, j.Header[0].Exclude),
		}}
	}

	if len(j.Host) > 0 {
		model.Host = []SourceV4CacheKeyFieldsHostModel{{
			Resolved: toBoolValue(j.Host[0].Resolved),
		}}
	}

	if len(j.QueryString) > 0 {
		model.QueryString = []SourceV4CacheKeyFieldsQueryStringModel{{
			Include: toSetValue(ctx, j.QueryString[0].Include),
			Exclude: toSetValue(ctx, j.QueryString[0].Exclude),
			Ignore:  toBoolValue(j.QueryString[0].Ignore),
		}}
	}

	if len(j.User) > 0 {
		model.User = []SourceV4CacheKeyFieldsUserModel{{
			DeviceType: toBoolValue(j.User[0].DeviceType),
			Geo:        toBoolValue(j.User[0].Geo),
			Lang:       toBoolValue(j.User[0].Lang),
		}}
	}

	return model
}

// ============================================================================
// Conversion: v5 JSON -> v5 Terraform Model
// ============================================================================

func convertV5JSONToModel(ctx context.Context, j SourceV5PageRuleJSON) *TargetV5PageRuleModel {
	model := &TargetV5PageRuleModel{
		ID:       toStringValue(j.ID),
		ZoneID:   toStringValue(j.ZoneID),
		Target:   toStringValue(j.Target),
		Priority: toInt64FromFloat(j.Priority),
		Status:   toStringValue(j.Status),
	}

	if j.Actions != nil {
		model.Actions = convertSourceV5ActionsJSON(ctx, j.Actions)
	}

	return model
}

func convertSourceV5ActionsJSON(ctx context.Context, j *SourceV5ActionsJSON) *TargetV5ActionsModel {
	model := &TargetV5ActionsModel{
		AlwaysUseHTTPS:          toBoolValue(j.AlwaysUseHTTPS),
		DisableApps:             toBoolValue(j.DisableApps),
		DisablePerformance:      toBoolValue(j.DisablePerformance),
		DisableSecurity:         toBoolValue(j.DisableSecurity),
		DisableZaraz:            toBoolValue(j.DisableZaraz),
		AutomaticHTTPSRewrites:  toStringPtrValue(j.AutomaticHTTPSRewrites),
		BrowserCacheTTL:         toInt64FromFloat(j.BrowserCacheTTL),
		BrowserCheck:            toStringPtrValue(j.BrowserCheck),
		BypassCacheOnCookie:     toStringPtrValue(j.BypassCacheOnCookie),
		CacheByDeviceType:       toStringPtrValue(j.CacheByDeviceType),
		CacheDeceptionArmor:     toStringPtrValue(j.CacheDeceptionArmor),
		CacheLevel:              toStringPtrValue(j.CacheLevel),
		CacheOnCookie:           toStringPtrValue(j.CacheOnCookie),
		EdgeCacheTTL:            toInt64FromFloat(j.EdgeCacheTTL),
		EmailObfuscation:        toStringPtrValue(j.EmailObfuscation),
		ExplicitCacheControl:    toStringPtrValue(j.ExplicitCacheControl),
		HostHeaderOverride:      toStringPtrValue(j.HostHeaderOverride),
		IPGeolocation:           toStringPtrValue(j.IPGeolocation),
		Mirage:                  toStringPtrValue(j.Mirage),
		OpportunisticEncryption: toStringPtrValue(j.OpportunisticEncryption),
		OriginErrorPagePassThru: toStringPtrValue(j.OriginErrorPagePassThru),
		Polish:                  toStringPtrValue(j.Polish),
		ResolveOverride:         toStringPtrValue(j.ResolveOverride),
		RespectStrongEtag:       toStringPtrValue(j.RespectStrongEtag),
		ResponseBuffering:       toStringPtrValue(j.ResponseBuffering),
		RocketLoader:            toStringPtrValue(j.RocketLoader),
		SecurityLevel:           toStringPtrValue(j.SecurityLevel),
		SortQueryStringForCache: toStringPtrValue(j.SortQueryStringForCache),
		SSL:                     toStringPtrValue(j.SSL),
		TrueClientIPHeader:      toStringPtrValue(j.TrueClientIPHeader),
		WAF:                     toStringPtrValue(j.WAF),
		CacheTTLByStatus:        types.MapNull(types.StringType),
	}

	if j.ForwardingURL != nil {
		model.ForwardingURL = &TargetV5ForwardingURLModel{
			URL:        toStringValue(j.ForwardingURL.URL),
			StatusCode: toInt64FromFloat(j.ForwardingURL.StatusCode),
		}
	}

	if j.CacheKeyFields != nil {
		model.CacheKeyFields = convertSourceV5CacheKeyFieldsJSON(ctx, j.CacheKeyFields)
	}

	if len(j.CacheTTLByStatus) > 0 {
		elements := make(map[string]attr.Value, len(j.CacheTTLByStatus))
		for k, v := range j.CacheTTLByStatus {
			elements[k] = types.StringValue(v)
		}
		model.CacheTTLByStatus = types.MapValueMust(types.StringType, elements)
	}

	return model
}

func convertSourceV5CacheKeyFieldsJSON(ctx context.Context, j *SourceV5CacheKeyFieldsJSON) *TargetV5CacheKeyFieldsModel {
	model := &TargetV5CacheKeyFieldsModel{}

	if j.Cookie != nil {
		model.Cookie = &TargetV5CacheKeyFieldsCookieModel{
			CheckPresence: toTypesStringSlice(j.Cookie.CheckPresence),
			Include:       toTypesStringSlice(j.Cookie.Include),
		}
	}

	if j.Header != nil {
		model.Header = &TargetV5CacheKeyFieldsHeaderModel{
			CheckPresence: toTypesStringSlice(j.Header.CheckPresence),
			Include:       toTypesStringSlice(j.Header.Include),
			Exclude:       toTypesStringSlice(j.Header.Exclude),
		}
	}

	if j.Host != nil {
		model.Host = &TargetV5CacheKeyFieldsHostModel{
			Resolved: toBoolValue(j.Host.Resolved),
		}
	}

	if j.QueryString != nil {
		model.QueryString = &TargetV5CacheKeyFieldsQueryStringModel{
			Include: toTypesStringSlice(j.QueryString.Include),
			Exclude: toTypesStringSlice(j.QueryString.Exclude),
		}
	}

	if j.User != nil {
		model.User = &TargetV5CacheKeyFieldsUserModel{
			DeviceType: toBoolValue(j.User.DeviceType),
			Geo:        toBoolValue(j.User.Geo),
			Lang:       toBoolValue(j.User.Lang),
		}
	}

	return model
}

// ============================================================================
// Type conversion helpers
// ============================================================================

func toStringValue(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func toStringPtrValue(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

func toInt64FromFloat(f *float64) types.Int64 {
	if f == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*f))
}

func toBoolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

func toSetValue(ctx context.Context, s []string) types.Set {
	if len(s) == 0 {
		return types.SetNull(types.StringType)
	}
	set, _ := types.SetValueFrom(ctx, types.StringType, s)
	return set
}

func toTypesStringSlice(s []string) []types.String {
	if s == nil {
		return nil
	}
	result := make([]types.String, len(s))
	for i, v := range s {
		result[i] = types.StringValue(v)
	}
	return result
}
