package adk

import (
	"fmt"
	"time"

	"github.com/run-bigpig/jcp/internal/adk/mcp"
	"github.com/run-bigpig/jcp/internal/adk/tools"
	"github.com/run-bigpig/jcp/internal/models"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

// ExpertAgentBuilder 专家 Agent 构建器
type ExpertAgentBuilder struct {
	llm          model.LLM
	toolRegistry *tools.Registry
	mcpManager   *mcp.Manager
}

// NewExpertAgentBuilder 创建专家 Agent 构建器
func NewExpertAgentBuilder(llm model.LLM) *ExpertAgentBuilder {
	return &ExpertAgentBuilder{llm: llm}
}

// NewExpertAgentBuilderWithTools 创建带工具的专家 Agent 构建器
func NewExpertAgentBuilderWithTools(llm model.LLM, registry *tools.Registry) *ExpertAgentBuilder {
	return &ExpertAgentBuilder{llm: llm, toolRegistry: registry}
}

// NewExpertAgentBuilderFull 创建完整配置的专家 Agent 构建器
func NewExpertAgentBuilderFull(llm model.LLM, registry *tools.Registry, mcpMgr *mcp.Manager) *ExpertAgentBuilder {
	return &ExpertAgentBuilder{llm: llm, toolRegistry: registry, mcpManager: mcpMgr}
}

// BuildAgent 根据配置构建 LLM Agent
func (b *ExpertAgentBuilder) BuildAgent(config *models.AgentConfig, stock *models.Stock, query string, position *models.StockPosition) (agent.Agent, error) {
	return b.BuildAgentWithContext(config, stock, query, "", position)
}

// BuildAgentWithContext 根据配置构建 LLM Agent（支持引用上下文）
func (b *ExpertAgentBuilder) BuildAgentWithContext(config *models.AgentConfig, stock *models.Stock, query string, replyContent string, position *models.StockPosition) (agent.Agent, error) {
	instruction := b.buildInstructionWithContext(config, stock, query, replyContent, position)

	// 获取 Agent 配置的工具
	var agentTools []tool.Tool
	if b.toolRegistry != nil && len(config.Tools) > 0 {
		agentTools = b.toolRegistry.GetTools(config.Tools)
	}

	// 获取 MCP toolsets
	var toolsets []tool.Toolset
	if b.mcpManager != nil && len(config.MCPServers) > 0 {
		toolsets = b.mcpManager.GetToolsetsByIDs(config.MCPServers)
	}

	return llmagent.New(llmagent.Config{
		Name:        config.ID,
		Model:       b.llm,
		Description: config.Role,
		Instruction: instruction,
		Tools:       agentTools,
		Toolsets:    toolsets,
	})
}

// buildInstruction 构建 Agent 指令
func (b *ExpertAgentBuilder) buildInstruction(config *models.AgentConfig, stock *models.Stock, query string, position *models.StockPosition) string {
	return b.buildInstructionWithContext(config, stock, query, "", position)
}

// buildInstructionWithContext 构建 Agent 指令（支持引用上下文）
func (b *ExpertAgentBuilder) buildInstructionWithContext(config *models.AgentConfig, stock *models.Stock, query string, replyContent string, position *models.StockPosition) string {
	baseInstruction := config.Instruction
	if baseInstruction == "" {
		baseInstruction = fmt.Sprintf("你是一位%s，名字是%s。", config.Role, config.Name)
	}

	// 构建可用工具说明
	toolsDescription := b.buildToolsDescription(config)

	// 获取当前时间和盘中状态
	now := time.Now()
	timeStr := now.Format("2006-01-02 15:04:05")
	weekday := now.Weekday()
	hour, minute := now.Hour(), now.Minute()
	currentMinutes := hour*60 + minute

	// 判断盘中状态（A股交易时间：9:30-11:30, 13:00-15:00，周一至周五）
	var marketStatus string
	if weekday == time.Saturday || weekday == time.Sunday {
		marketStatus = "休市（周末）"
	} else if currentMinutes >= 9*60+30 && currentMinutes <= 11*60+30 {
		marketStatus = "盘中（上午交易时段）"
	} else if currentMinutes >= 13*60 && currentMinutes <= 15*60 {
		marketStatus = "盘中（下午交易时段）"
	} else if currentMinutes < 9*60+30 {
		marketStatus = "盘前"
	} else if currentMinutes > 15*60 {
		marketStatus = "盘后"
	} else {
		marketStatus = "午间休市"
	}

	prompt := fmt.Sprintf(`%s
%s
当前时间: %s
市场状态: %s

股票: %s (%s)
当前价格: %.2f
涨跌幅: %.2f%%
`, baseInstruction, toolsDescription, timeStr, marketStatus, stock.Symbol, stock.Name, stock.Price, stock.ChangePercent)

	// 如果有持仓信息，加入上下文
	if position != nil && position.Shares > 0 {
		marketValue := float64(position.Shares) * stock.Price
		costAmount := float64(position.Shares) * position.CostPrice
		profitLoss := marketValue - costAmount
		profitPercent := 0.0
		if costAmount > 0 {
			profitPercent = (profitLoss / costAmount) * 100
		}
		prompt += fmt.Sprintf(`
用户持仓: %d股，成本价 %.2f
持仓市值: %.2f，盈亏: %.2f (%.2f%%)
`, position.Shares, position.CostPrice, marketValue, profitLoss, profitPercent)
	}

	// 如果有引用内容，加入上下文
	if replyContent != "" {
		prompt += fmt.Sprintf(`--- 引用的观点 ---
%s
---

小韭菜问题: %s

请结合以上引用的观点，发表你的看法。可以赞同、补充或反驳。回复控制在150字以内。`, replyContent, query)
	} else {
		prompt += fmt.Sprintf(`小韭菜问题: %s

请用简洁专业的语言回答，控制在150字以内。`, query)
	}

	return prompt
}

// buildToolsDescription 构建可用工具说明
func (b *ExpertAgentBuilder) buildToolsDescription(config *models.AgentConfig) string {
	var toolDescriptions []string

	// 获取内置工具信息
	if b.toolRegistry != nil && len(config.Tools) > 0 {
		toolInfos := b.toolRegistry.GetToolInfosByNames(config.Tools)
		for _, info := range toolInfos {
			toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", info.Name, info.Description))
		}
	}

	// 获取 MCP 工具信息
	if b.mcpManager != nil && len(config.MCPServers) > 0 {
		mcpTools := b.mcpManager.GetToolInfosByServerIDs(config.MCPServers)
		for _, info := range mcpTools {
			toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s (来自 %s)", info.Name, info.Description, info.ServerName))
		}
	}

	if len(toolDescriptions) == 0 {
		return ""
	}

	result := "\n可用工具:\n"
	for _, desc := range toolDescriptions {
		result += desc + "\n"
	}
	return result
}
