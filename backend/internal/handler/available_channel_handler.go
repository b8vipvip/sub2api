package handler

import (
	"context"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AvailableChannelHandler 处理用户侧「模型广场」查询。
//
// 用户侧接口只返回模型广场需要的白名单字段，不暴露账号凭证、上游密钥、账号 ID 等敏感字段。
type AvailableChannelHandler struct {
	channelService *service.ChannelService
	apiKeyService  *service.APIKeyService
	settingService *service.SettingService
	accountRepo    service.AccountRepository
}

// NewAvailableChannelHandler 创建用户侧模型广场 handler。
func NewAvailableChannelHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	settingService *service.SettingService,
	accountRepo service.AccountRepository,
) *AvailableChannelHandler {
	return &AvailableChannelHandler{
		channelService: channelService,
		apiKeyService:  apiKeyService,
		settingService: settingService,
		accountRepo:    accountRepo,
	}
}

// featureEnabled 返回 available-channels 开关是否启用。默认关闭（opt-in）。
func (h *AvailableChannelHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return false
	}
	return h.settingService.GetAvailableChannelsRuntime(c.Request.Context()).Enabled
}

// userAvailableGroup 用户可见的分组概要（白名单字段）。
type userAvailableGroup struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Platform         string  `json:"platform"`
	SubscriptionType string  `json:"subscription_type"`
	RateMultiplier   float64 `json:"rate_multiplier"`
	IsExclusive      bool    `json:"is_exclusive"`
}

// userSupportedModelPricing 用户可见的定价字段白名单。
type userSupportedModelPricing struct {
	BillingMode      string                   `json:"billing_mode"`
	InputPrice       *float64                 `json:"input_price"`
	OutputPrice      *float64                 `json:"output_price"`
	CacheWritePrice  *float64                 `json:"cache_write_price"`
	CacheReadPrice   *float64                 `json:"cache_read_price"`
	ImageOutputPrice *float64                 `json:"image_output_price"`
	PerRequestPrice  *float64                 `json:"per_request_price"`
	Intervals        []userPricingIntervalDTO `json:"intervals"`
}

// userPricingIntervalDTO 定价区间白名单（去掉内部 ID、SortOrder 等前端不渲染的字段）。
type userPricingIntervalDTO struct {
	MinTokens       int      `json:"min_tokens"`
	MaxTokens       *int     `json:"max_tokens"`
	TierLabel       string   `json:"tier_label,omitempty"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

// userSupportedModel 用户可见的支持模型条目。
type userSupportedModel struct {
	Name     string                     `json:"name"`
	Platform string                     `json:"platform"`
	Pricing  *userSupportedModelPricing `json:"pricing"`
}

// userChannelPlatformSection 单渠道内某个平台的子视图。
type userChannelPlatformSection struct {
	Platform        string               `json:"platform"`
	Groups          []userAvailableGroup `json:"groups"`
	SupportedModels []userSupportedModel `json:"supported_models"`
}

// userAvailableChannel 用户可见的渠道条目（白名单字段）。
type userAvailableChannel struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Platforms   []userChannelPlatformSection `json:"platforms"`
}

// List 列出当前用户可见的模型广场。
// GET /api/v1/channels/available
func (h *AvailableChannelHandler) List(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	if !h.featureEnabled(c) {
		response.Success(c, []userAvailableChannel{})
		return
	}

	userGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	visibleGroups := groupsToUserAvailable(userGroups)
	accountModelsByPlatform, err := h.accountConfiguredModelsByPlatform(c.Request.Context(), visibleGroups)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if len(accountModelsByPlatform) > 0 {
		sections := buildAccountModelSections(visibleGroups, accountModelsByPlatform)
		if len(sections) > 0 {
			response.Success(c, []userAvailableChannel{{
				Name:        "模型广场",
				Description: "从管理员后台调度开关开启账号的模型白名单自动汇总。",
				Platforms:   sections,
			}})
			return
		}
	}

	// 没有账号白名单时，保留旧逻辑：从渠道配置里展示模型。
	allowedGroupIDs := make(map[int64]struct{}, len(userGroups))
	for i := range userGroups {
		allowedGroupIDs[userGroups[i].ID] = struct{}{}
	}

	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]userAvailableChannel, 0, len(channels))
	for _, ch := range channels {
		if ch.Status != service.StatusActive {
			continue
		}
		visibleGroups := filterUserVisibleGroups(ch.Groups, allowedGroupIDs)
		if len(visibleGroups) == 0 {
			continue
		}
		fallbackModelsByPlatform, err := h.accountConfiguredModelsByPlatform(c.Request.Context(), visibleGroups)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		sections := buildPlatformSections(ch, visibleGroups, fallbackModelsByPlatform)
		if len(sections) == 0 {
			continue
		}
		out = append(out, userAvailableChannel{
			Name:        ch.Name,
			Description: ch.Description,
			Platforms:   sections,
		})
	}

	response.Success(c, out)
}

func groupsToUserAvailable(groups []service.AvailableGroupRef) []userAvailableGroup {
	out := make([]userAvailableGroup, 0, len(groups))
	seen := make(map[int64]struct{}, len(groups))
	for _, g := range groups {
		if g.ID <= 0 || g.Platform == "" {
			continue
		}
		if _, ok := seen[g.ID]; ok {
			continue
		}
		seen[g.ID] = struct{}{}
		out = append(out, userAvailableGroup{
			ID:               g.ID,
			Name:             g.Name,
			Platform:         g.Platform,
			SubscriptionType: g.SubscriptionType,
			RateMultiplier:   g.RateMultiplier,
			IsExclusive:      g.IsExclusive,
		})
	}
	return out
}

func buildAccountModelSections(
	visibleGroups []userAvailableGroup,
	modelsByPlatform map[string][]userSupportedModel,
) []userChannelPlatformSection {
	groupsByPlatform := make(map[string][]userAvailableGroup, 4)
	for _, g := range visibleGroups {
		if g.Platform == "" {
			continue
		}
		groupsByPlatform[g.Platform] = append(groupsByPlatform[g.Platform], g)
	}

	platforms := make([]string, 0, len(modelsByPlatform))
	for platform, models := range modelsByPlatform {
		if platform != "" && len(models) > 0 {
			platforms = append(platforms, platform)
		}
	}
	sort.Strings(platforms)

	sections := make([]userChannelPlatformSection, 0, len(platforms))
	for _, platform := range platforms {
		sections = append(sections, userChannelPlatformSection{
			Platform:        platform,
			Groups:          groupsByPlatform[platform],
			SupportedModels: modelsByPlatform[platform],
		})
	}
	return sections
}

// buildPlatformSections 把一个渠道按 visibleGroups 的平台集合拆成有序的 section 列表。
func buildPlatformSections(
	ch service.AvailableChannel,
	visibleGroups []userAvailableGroup,
	fallbackModelsByPlatform map[string][]userSupportedModel,
) []userChannelPlatformSection {
	groupsByPlatform := make(map[string][]userAvailableGroup, 4)
	for _, g := range visibleGroups {
		if g.Platform == "" {
			continue
		}
		groupsByPlatform[g.Platform] = append(groupsByPlatform[g.Platform], g)
	}
	if len(groupsByPlatform) == 0 {
		return nil
	}

	platforms := make([]string, 0, len(groupsByPlatform))
	for p := range groupsByPlatform {
		platforms = append(platforms, p)
	}
	sort.Strings(platforms)

	sections := make([]userChannelPlatformSection, 0, len(platforms))
	for _, platform := range platforms {
		platformSet := map[string]struct{}{platform: {}}
		models := toUserSupportedModels(ch.SupportedModels, platformSet)
		if len(models) == 0 && fallbackModelsByPlatform != nil {
			models = fallbackModelsByPlatform[platform]
		}
		sections = append(sections, userChannelPlatformSection{
			Platform:        platform,
			Groups:          groupsByPlatform[platform],
			SupportedModels: models,
		})
	}
	return sections
}

// filterUserVisibleGroups 仅保留用户可访问的分组。
func filterUserVisibleGroups(
	groups []service.AvailableGroupRef,
	allowed map[int64]struct{},
) []userAvailableGroup {
	visible := make([]userAvailableGroup, 0, len(groups))
	for _, g := range groups {
		if _, ok := allowed[g.ID]; !ok {
			continue
		}
		visible = append(visible, userAvailableGroup{
			ID:               g.ID,
			Name:             g.Name,
			Platform:         g.Platform,
			SubscriptionType: g.SubscriptionType,
			RateMultiplier:   g.RateMultiplier,
			IsExclusive:      g.IsExclusive,
		})
	}
	return visible
}

// accountConfiguredModelsByPlatform 从当前用户可见分组中的调度账号配置里汇总模型列表。
func (h *AvailableChannelHandler) accountConfiguredModelsByPlatform(
	ctx context.Context,
	groups []userAvailableGroup,
) (map[string][]userSupportedModel, error) {
	if h.accountRepo == nil || h.channelService == nil || len(groups) == 0 {
		return nil, nil
	}

	seenGroups := make(map[int64]struct{}, len(groups))
	seenModels := make(map[string]map[string]struct{})
	modelsByPlatform := make(map[string][]userSupportedModel)

	for _, g := range groups {
		if g.ID <= 0 || g.Platform == "" {
			continue
		}
		if _, ok := seenGroups[g.ID]; ok {
			continue
		}
		seenGroups[g.ID] = struct{}{}

		accounts, err := h.accountRepo.ListByGroup(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		for i := range accounts {
			acc := &accounts[i]
			if acc == nil || !acc.IsActive() || !acc.Schedulable || acc.Platform != g.Platform {
				continue
			}
			for _, model := range accountConfiguredModelNames(acc) {
				if model == "" {
					continue
				}
				if seenModels[g.Platform] == nil {
					seenModels[g.Platform] = make(map[string]struct{})
				}
				modelKey := strings.ToLower(model)
				if _, exists := seenModels[g.Platform][modelKey]; exists {
					continue
				}
				seenModels[g.Platform][modelKey] = struct{}{}
				modelsByPlatform[g.Platform] = append(modelsByPlatform[g.Platform], userSupportedModel{
					Name:     model,
					Platform: g.Platform,
					Pricing:  toUserPricing(h.channelService.GetChannelModelPricing(ctx, g.ID, model)),
				})
			}
		}
	}

	for platform := range modelsByPlatform {
		sort.Slice(modelsByPlatform[platform], func(i, j int) bool {
			return modelsByPlatform[platform][i].Name < modelsByPlatform[platform][j].Name
		})
	}
	return modelsByPlatform, nil
}

func accountConfiguredModelNames(acc *service.Account) []string {
	if acc == nil {
		return nil
	}
	seen := make(map[string]struct{})
	out := make([]string, 0)

	add := func(model string) {
		model = strings.TrimSpace(model)
		if model == "" {
			return
		}
		key := strings.ToLower(model)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		out = append(out, model)
	}

	// 账号模型限制当前主要存储为 credentials.model_mapping：
	// - 模型白名单模式：{"gpt-5.5":"gpt-5.5"}
	// - 模型映射模式：{"custom-name":"upstream-name"}
	// 对用户端来说，应该展示客户端可填写的请求模型 ID，即 map key。
	for model := range acc.GetModelMapping() {
		add(model)
	}

	// 兼容旧数据：credentials.model_whitelist 可能是 []string 或 []any。
	if acc.Credentials != nil {
		switch raw := acc.Credentials["model_whitelist"].(type) {
		case []string:
			for _, model := range raw {
				add(model)
			}
		case []any:
			for _, item := range raw {
				if model, ok := item.(string); ok {
					add(model)
				}
			}
		}
	}

	return out
}

// toUserSupportedModels 将 service 层支持模型转换为用户 DTO（字段白名单）。
func toUserSupportedModels(
	src []service.SupportedModel,
	allowedPlatforms map[string]struct{},
) []userSupportedModel {
	out := make([]userSupportedModel, 0, len(src))
	for i := range src {
		m := src[i]
		if allowedPlatforms != nil {
			if _, ok := allowedPlatforms[m.Platform]; !ok {
				continue
			}
		}
		out = append(out, userSupportedModel{
			Name:     m.Name,
			Platform: m.Platform,
			Pricing:  toUserPricing(m.Pricing),
		})
	}
	return out
}

// toUserPricing 将 service 层定价转换为用户 DTO；入参为 nil 时返回 nil。
func toUserPricing(p *service.ChannelModelPricing) *userSupportedModelPricing {
	if p == nil {
		return nil
	}
	intervals := make([]userPricingIntervalDTO, 0, len(p.Intervals))
	for _, iv := range p.Intervals {
		intervals = append(intervals, userPricingIntervalDTO{
			MinTokens:       iv.MinTokens,
			MaxTokens:       iv.MaxTokens,
			TierLabel:       iv.TierLabel,
			InputPrice:      iv.InputPrice,
			OutputPrice:     iv.OutputPrice,
			CacheWritePrice: iv.CacheWritePrice,
			CacheReadPrice:  iv.CacheReadPrice,
			PerRequestPrice: iv.PerRequestPrice,
		})
	}
	billingMode := string(p.BillingMode)
	if billingMode == "" {
		billingMode = string(service.BillingModeToken)
	}
	return &userSupportedModelPricing{
		BillingMode:      billingMode,
		InputPrice:       p.InputPrice,
		OutputPrice:      p.OutputPrice,
		CacheWritePrice:  p.CacheWritePrice,
		CacheReadPrice:   p.CacheReadPrice,
		ImageOutputPrice: p.ImageOutputPrice,
		PerRequestPrice:  p.PerRequestPrice,
		Intervals:        intervals,
	}
}
