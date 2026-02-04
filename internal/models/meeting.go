package models

// MeetingConfig 会议配置
type MeetingConfig struct {
	MaxRounds       int  `json:"maxRounds"`       // 最大讨论轮数，默认 2
	EnableCrossTalk bool `json:"enableCrossTalk"` // 启用交叉讨论
	AutoSummarize   bool `json:"autoSummarize"`   // 自动总结
}

// DefaultMeetingConfig 默认会议配置
func DefaultMeetingConfig() MeetingConfig {
	return MeetingConfig{
		MaxRounds:       2,
		EnableCrossTalk: true,
		AutoSummarize:   true,
	}
}

// AgentProfile Agent 简介（用于互相了解）
type AgentProfile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Expertise   string `json:"expertise"`   // 擅长领域
	Perspective string `json:"perspective"` // 分析视角
}

// ToProfile 将 AgentConfig 转换为 AgentProfile
func (c *AgentConfig) ToProfile() AgentProfile {
	return AgentProfile{
		ID:          c.ID,
		Name:        c.Name,
		Role:        c.Role,
		Expertise:   c.Role,
		Perspective: c.Instruction,
	}
}
