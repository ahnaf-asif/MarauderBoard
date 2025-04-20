package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/ahnafasif/MarauderBoard/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var conversationStore = struct {
	sync.RWMutex
	data map[string][]Message
}{data: make(map[string][]Message)}

func GetOrCreateOllamaSession(ctx *fiber.Ctx) string {
	sessionID := ctx.Cookies("ollama_session")
	if sessionID == "" {
		sessionID = uuid.NewString()
		ctx.Cookie(&fiber.Cookie{
			Name:   "ollama_session",
			Value:  sessionID,
			MaxAge: 60 * 60 * 24, // 24 hour
		})
	}
	return sessionID
}

func GetMessages(sessionID string) []Message {
	conversationStore.RLock()
	defer conversationStore.RUnlock()
	return conversationStore.data[sessionID]
}

func SetMessages(sessionID string, messages []Message) {
	conversationStore.Lock()
	defer conversationStore.Unlock()
	conversationStore.data[sessionID] = messages
}

func cleanOllamaResponse(content string) string {
	content = strings.ReplaceAll(content, "<thinking>", "")
	content = strings.ReplaceAll(content, "</thinking>", "")
	content = strings.ReplaceAll(content, "<think>", "")
	return strings.ReplaceAll(content, "</think>", "")
}

func SendToOllama(model string, messages []Message) (Message, error) {
	payload := map[string]interface{}{
		"model":    model,
		"stream":   false,
		"messages": messages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}

	resp, err := http.Post(configs.OllamaApiUri, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, err
	}

	var ollamaResponse struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(body, &ollamaResponse); err != nil {
		return Message{}, errors.New("invalid response from Ollama")
	}
	ollamaResponse.Message.Content = cleanOllamaResponse(ollamaResponse.Message.Content)
	return ollamaResponse.Message, nil
}

var messageTemplate = template.Must(template.New("message").Parse(`
{{- range . }}
<div class="{{if eq .Role "user"}}text-right{{else}}text-left{{end}} mb-3 last:mb-0">
  <div class="inline-block max-w-[85%] px-4 py-3 rounded-lg 
             {{if eq .Role "user"}}
               bg-indigo-100 text-indigo-900 rounded-br-none
             {{else}}
               bg-gray-100 text-gray-800 rounded-bl-none
             {{end}}">
    {{ .Content }}
  </div>
</div>
{{- end }}`))

type RenderableMessage struct {
	Role    string
	Content template.HTML
}

func RenderMessagesHTML(messages []Message) (string, error) {
	var buf bytes.Buffer

	var renderableMessages []RenderableMessage
	for _, message := range messages {
		renderableMessages = append(renderableMessages, RenderableMessage{
			Role:    message.Role,
			Content: template.HTML(message.Content),
		})
	}
	err := messageTemplate.Execute(&buf, renderableMessages)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

