import { GetAgentConfigs, AddAgentConfig, UpdateAgentConfig, DeleteAgentConfig } from '../../wailsjs/go/main/App';

export interface AgentConfig {
  id: string;
  name: string;
  role: string;
  avatar: string;
  color: string;
  instruction: string;
  tools: string[];
  mcpServers: string[];  // 关联的 MCP 服务器 ID 列表
  priority: number;
  isBuiltin: boolean;
  enabled: boolean;
  providerId: string;  // 关联的Provider ID（空则使用默认）
}

// 获取所有Agent配置
export const getAgentConfigs = async (): Promise<AgentConfig[]> => {
  return await GetAgentConfigs();
};

// 添加Agent配置
export const addAgentConfig = async (config: AgentConfig): Promise<string> => {
  return await AddAgentConfig(config);
};

// 更新Agent配置
export const updateAgentConfig = async (config: AgentConfig): Promise<string> => {
  return await UpdateAgentConfig(config);
};

// 删除Agent配置
export const deleteAgentConfig = async (id: string): Promise<string> => {
  return await DeleteAgentConfig(id);
};
