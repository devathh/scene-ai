package openrouterhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devathh/scene-ai/internal/common/config"
	"github.com/devathh/scene-ai/internal/modules/ai/domain/scenario"
	"github.com/google/uuid"
)

type OpenRouterRepository struct {
	cfg    *config.Config
	client *http.Client
	model  string
}

func New(cfg *config.Config, client *http.Client, model string) *OpenRouterRepository {
	return &OpenRouterRepository{
		cfg:    cfg,
		client: client,
		model:  model,
	}
}

func (orr *OpenRouterRepository) GenerateScenario(ctx context.Context, prompt string, authorID uuid.UUID) (*scenario.Scenario, error) {
	requestBytes, err := orr.buildRequestBody(orr.cfg.OpenRouter.ScenarioPromptPath, prompt)
	if err != nil {
		return nil, err
	}

	respBody, err := orr.doRequest(ctx, requestBytes)
	if err != nil {
		return nil, err
	}

	var response response
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) < 1 {
		return nil, errors.New("invalid response from ai: no choices")
	}

	var result responseScenario
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scenario content: %w", err)
	}

	return scenario.New(
		authorID,
		result.Title,
		prompt,
		result.GlobalStylePrompt,
		[]scenario.Scene{},
	)
}

func (orr *OpenRouterRepository) GenerateScenes(ctx context.Context, start *scenario.Scenario, handler func(scenario.Scene)) error {
	startPrompt, err := orr.loadPrompt(orr.cfg.OpenRouter.ScenePromptPath)
	if err != nil {
		return fmt.Errorf("failed to load scene prompt: %w", err)
	}

	messages := []message{
		{Role: "system", Content: startPrompt},
		{Role: "user", Content: orr.promptFromScenario(start)},
	}

	for {
		reqBytes, err := json.Marshal(&request{
			Model:    orr.model,
			Messages: messages,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		respBody, err := orr.doRequest(ctx, reqBytes)
		if err != nil {
			return err
		}

		var response response
		if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if len(response.Choices) < 1 {
			return errors.New("invalid response from ai: no choices")
		}

		choice := response.Choices[0]
		var result responseScene
		if err := json.Unmarshal([]byte(choice.Message.Content), &result); err != nil {
			return fmt.Errorf("failed to unmarshal scene content: %w", err)
		}

		if result.Status == "finished" {
			return nil
		}

		scene, err := scenario.NewScene(
			result.Order,
			result.Title,
			time.Duration(result.DurationSec)*time.Second,
			result.VideoPrompt,
		)
		if err != nil {
			return fmt.Errorf("failed to create scene domain object: %w", err)
		}

		messages = append(messages, message{
			Role:    choice.Message.Role,
			Content: choice.Message.Content,
		})
		messages = append(messages, message{
			Role:    "user",
			Content: "ok",
		})

		handler(scene)
		time.Sleep(1 * time.Second)
	}
}

func (orr *OpenRouterRepository) doRequest(ctx context.Context, bodyBytes []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+orr.cfg.OpenRouter.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := orr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf.Bytes(), nil
}

func (orr *OpenRouterRepository) buildRequestBody(startPromptPath, prompt string) ([]byte, error) {
	startPrompt, err := orr.loadPrompt(startPromptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load scenario prompt: %w", err)
	}

	req := &request{
		Model: orr.model,
		Messages: []message{
			{Role: "system", Content: startPrompt},
			{Role: "user", Content: prompt},
		},
	}

	return json.Marshal(req)
}

func (orr *OpenRouterRepository) promptFromScenario(s *scenario.Scenario) string {
	return fmt.Sprintf("Title of scenario: %s\nScenario prompt: %s\nGlobal style prompt: %s",
		s.Title(),
		s.ScenarioPrompt(),
		s.GlobalStylePrompt(),
	)
}

func (orr *OpenRouterRepository) loadPrompt(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", filepath, err)
	}
	return string(data), nil
}
