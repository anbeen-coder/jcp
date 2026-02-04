// Package mcp 提供 MCP (Model Context Protocol) 集成功能
package mcp

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/run-bigpig/jcp/internal/models"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

// ServerStatus MCP 服务器状态
type ServerStatus struct {
	ID        string `json:"id"`
	Connected bool   `json:"connected"`
	Error     string `json:"error"`
}

// ToolInfo MCP 工具信息
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServerID    string `json:"serverId"`
	ServerName  string `json:"serverName"`
}

// Manager MCP 服务管理器
type Manager struct {
	mu       sync.RWMutex
	toolsets map[string]tool.Toolset
	configs  map[string]*models.MCPServerConfig
}

// NewManager 创建 MCP 管理器
func NewManager() *Manager {
	return &Manager{
		toolsets: make(map[string]tool.Toolset),
		configs:  make(map[string]*models.MCPServerConfig),
	}
}

// LoadConfigs 加载 MCP 服务器配置
func (m *Manager) LoadConfigs(configs []models.MCPServerConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.toolsets = make(map[string]tool.Toolset)
	m.configs = make(map[string]*models.MCPServerConfig)

	for i := range configs {
		cfg := &configs[i]
		if !cfg.Enabled {
			continue
		}
		m.configs[cfg.ID] = cfg

		ts, err := mcptoolset.New(mcptoolset.Config{
			Transport:  createTransport(cfg),
			ToolFilter: tool.StringPredicate(cfg.ToolFilter),
		})
		if err == nil {
			m.toolsets[cfg.ID] = ts
		}
	}
	return nil
}

// createTransport 根据配置创建 MCP 传输层
func createTransport(cfg *models.MCPServerConfig) mcp.Transport {
	switch cfg.TransportType {
	case models.MCPTransportSSE:
		return &mcp.SSEClientTransport{Endpoint: cfg.Endpoint}
	case models.MCPTransportCommand:
		return &mcp.CommandTransport{Command: exec.Command(cfg.Command, cfg.Args...)}
	default: // http
		return &mcp.StreamableClientTransport{Endpoint: cfg.Endpoint}
	}
}

// GetToolset 获取指定 MCP 服务器的 toolset
func (m *Manager) GetToolset(serverID string) (tool.Toolset, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ts, ok := m.toolsets[serverID]
	return ts, ok
}

// GetAllToolsets 获取所有已启用的 toolsets
func (m *Manager) GetAllToolsets() []tool.Toolset {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]tool.Toolset, 0, len(m.toolsets))
	for _, ts := range m.toolsets {
		result = append(result, ts)
	}
	return result
}

// GetToolsetsByIDs 根据 ID 列表获取 toolsets
func (m *Manager) GetToolsetsByIDs(ids []string) []tool.Toolset {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []tool.Toolset
	for _, id := range ids {
		if ts, ok := m.toolsets[id]; ok {
			result = append(result, ts)
		}
	}
	return result
}

// TestConnection 测试指定 MCP 服务器的连接
func (m *Manager) TestConnection(serverID string) *ServerStatus {
	m.mu.RLock()
	cfg, ok := m.configs[serverID]
	m.mu.RUnlock()

	if !ok {
		return &ServerStatus{ID: serverID, Connected: false, Error: "服务器未配置"}
	}

	// 使用 MCP SDK 原生 Client.Connect 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	impl := &mcp.Implementation{Name: cfg.Name, Version: "1.0.0"}
	client := mcp.NewClient(impl, nil)
	_, err := client.Connect(ctx, createTransport(cfg), nil)

	if err != nil {
		return &ServerStatus{ID: serverID, Connected: false, Error: err.Error()}
	}
	return &ServerStatus{ID: serverID, Connected: true}
}

// GetAllStatus 获取所有服务器状态
func (m *Manager) GetAllStatus() []ServerStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]ServerStatus, 0, len(m.configs))
	for id := range m.configs {
		result = append(result, ServerStatus{ID: id})
	}
	return result
}

// GetServerTools 获取指定 MCP 服务器的工具列表
func (m *Manager) GetServerTools(serverID string) ([]ToolInfo, error) {
	m.mu.RLock()
	cfg, ok := m.configs[serverID]
	m.mu.RUnlock()

	if !ok {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	impl := &mcp.Implementation{Name: cfg.Name, Version: "1.0.0"}
	client := mcp.NewClient(impl, nil)
	session, err := client.Connect(ctx, createTransport(cfg), nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// 获取工具列表
	toolsResp, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}

	var tools []ToolInfo
	for _, t := range toolsResp.Tools {
		tools = append(tools, ToolInfo{
			Name:        t.Name,
			Description: t.Description,
			ServerID:    serverID,
			ServerName:  cfg.Name,
		})
	}
	return tools, nil
}

// GetAllServerTools 获取所有已启用 MCP 服务器的工具列表
func (m *Manager) GetAllServerTools() []ToolInfo {
	m.mu.RLock()
	serverIDs := make([]string, 0, len(m.configs))
	for id := range m.configs {
		serverIDs = append(serverIDs, id)
	}
	m.mu.RUnlock()

	var allTools []ToolInfo
	for _, id := range serverIDs {
		tools, err := m.GetServerTools(id)
		if err == nil && tools != nil {
			allTools = append(allTools, tools...)
		}
	}
	return allTools
}

// GetToolInfosByServerIDs 根据服务器 ID 列表获取工具信息
func (m *Manager) GetToolInfosByServerIDs(serverIDs []string) []ToolInfo {
	var allTools []ToolInfo
	for _, id := range serverIDs {
		tools, err := m.GetServerTools(id)
		if err == nil && tools != nil {
			allTools = append(allTools, tools...)
		}
	}
	return allTools
}
