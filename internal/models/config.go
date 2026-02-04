package models

// AIProvider AI服务提供商类型
type AIProvider string

const (
	AIProviderOpenAI AIProvider = "openai"
	AIProviderGemini AIProvider = "gemini"
)

// AIConfig AI服务配置
type AIConfig struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Provider         AIProvider `json:"provider"`
	BaseURL          string     `json:"baseUrl"`
	APIKey           string     `json:"apiKey"`
	ModelName        string     `json:"modelName"`
	MaxTokens        int        `json:"maxTokens"`
	Temperature      float64    `json:"temperature"`
	Timeout          int        `json:"timeout"`
	HttpProxy        string     `json:"httpProxy"`
	HttpProxyEnabled bool       `json:"httpProxyEnabled"`
	IsDefault        bool       `json:"isDefault"`
}

// MCPTransportType MCP传输类型
type MCPTransportType string

const (
	MCPTransportHTTP    MCPTransportType = "http"    // StreamableHTTP 传输
	MCPTransportSSE     MCPTransportType = "sse"     // SSE 传输（已废弃）
	MCPTransportCommand MCPTransportType = "command" // 命令行传输
)

// MCPServerConfig MCP服务器配置
type MCPServerConfig struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	TransportType MCPTransportType `json:"transportType"`
	Endpoint      string           `json:"endpoint"`      // HTTP/SSE 端点 URL
	Command       string           `json:"command"`       // 命令行传输的命令
	Args          []string         `json:"args"`          // 命令行参数
	ToolFilter    []string         `json:"toolFilter"`    // 工具过滤列表（空则全部）
	Enabled       bool             `json:"enabled"`       // 是否启用
}

// AppConfig 应用配置
type AppConfig struct {
	RefreshInterval int64             `json:"refreshInterval"`
	DarkTheme       bool              `json:"darkTheme"`
	AIConfigs       []AIConfig        `json:"aiConfigs"`
	DefaultAIID     string            `json:"defaultAiId"`
	MCPServers      []MCPServerConfig `json:"mcpServers"` // MCP服务器配置列表
}
