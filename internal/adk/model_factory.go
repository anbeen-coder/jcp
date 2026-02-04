package adk

import (
	"context"
	"fmt"

	"github.com/run-bigpig/jcp/internal/adk/openai"
	"github.com/run-bigpig/jcp/internal/models"

	go_openai "github.com/sashabaranov/go-openai"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

// ModelFactory 模型工厂，根据配置创建对应的 adk model
type ModelFactory struct{}

// NewModelFactory 创建模型工厂
func NewModelFactory() *ModelFactory {
	return &ModelFactory{}
}

// CreateModel 根据 AI 配置创建对应的模型
func (f *ModelFactory) CreateModel(ctx context.Context, config *models.AIConfig) (model.LLM, error) {
	switch config.Provider {
	case models.AIProviderGemini:
		return f.createGeminiModel(ctx, config)
	case models.AIProviderOpenAI:
		return f.createOpenAIModel(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

// createGeminiModel 创建 Gemini 模型
func (f *ModelFactory) createGeminiModel(ctx context.Context, config *models.AIConfig) (model.LLM, error) {
	clientConfig := &genai.ClientConfig{
		APIKey: config.APIKey,
	}

	if config.BaseURL != "" {
		clientConfig.Backend = genai.BackendGeminiAPI
	}

	return gemini.NewModel(ctx, config.ModelName, clientConfig)
}

// createOpenAIModel 创建 OpenAI 兼容模型
func (f *ModelFactory) createOpenAIModel(config *models.AIConfig) (model.LLM, error) {
	openaiCfg := go_openai.DefaultConfig(config.APIKey)

	if config.BaseURL != "" {
		openaiCfg.BaseURL = config.BaseURL
	}

	return openai.NewOpenAIModel(config.ModelName, openaiCfg), nil
}
