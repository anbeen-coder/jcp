// 配置服务 - 调用后端API
import { GetConfig, UpdateConfig, GetAvailableTools } from '@wailsjs/go/main/App';
import type { models } from '@wailsjs/go/models';

export type AppConfig = models.AppConfig;

// 内置工具信息
export interface ToolInfo {
  name: string;
  description: string;
}

export const getConfig = async (): Promise<AppConfig> => {
  return await GetConfig();
};

export const updateConfig = async (config: AppConfig): Promise<string> => {
  return await UpdateConfig(config);
};

// 获取可用的内置工具列表
export const getAvailableTools = async (): Promise<ToolInfo[]> => {
  return await GetAvailableTools();
};
