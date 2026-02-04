package models

// AgentConfig Agent配置（可通过配置文件添加）
type AgentConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Avatar      string   `json:"avatar"`
	Color       string   `json:"color"`
	Instruction string   `json:"instruction"` // Agent系统指令
	Tools       []string `json:"tools"`       // 可用工具列表
	MCPServers  []string `json:"mcpServers"`  // 关联的 MCP 服务器 ID 列表
	Priority    int      `json:"priority"`    // 显示优先级
	IsBuiltin   bool     `json:"isBuiltin"`   // 是否内置Agent
	Enabled     bool     `json:"enabled"`     // 是否全局启用
	ProviderID  string   `json:"providerId"`  // 关联的Provider ID（空则使用默认）
}

// Agent AI分析师（兼容旧结构）
type Agent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
	Color  string `json:"color"`
}

// ToAgent 将AgentConfig转换为Agent
func (c *AgentConfig) ToAgent() Agent {
	return Agent{
		ID:     c.ID,
		Name:   c.Name,
		Role:   c.Role,
		Avatar: c.Avatar,
		Color:  c.Color,
	}
}
