package webhooks

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/yourdudeken/mpesa-sdk/client"
	"github.com/yourdudeken/mpesa-sdk/types"
)

type EventType string

const (
	EventSTKCallback       EventType = "stk:callback"
	EventB2CResult         EventType = "b2c:result"
	EventB2BResult         EventType = "b2b:result"
	EventReversalResult    EventType = "reversal:result"
	EventTransactionStatus EventType = "transaction:status"
	EventAccountBalance    EventType = "account:balance"
	EventC2BValidation     EventType = "c2b:validation"
)

type WebhookHandler func(eventType EventType, payload interface{})

type Manager struct {
	mu       sync.RWMutex
	handlers map[EventType][]WebhookHandler
}

func NewManager() *Manager {
	return &Manager{
		handlers: make(map[EventType][]WebhookHandler),
	}
}

func (m *Manager) On(event EventType, handler WebhookHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[event] = append(m.handlers[event], handler)
}

func (m *Manager) Off(event EventType, handler WebhookHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	handlers := m.handlers[event]
	filtered := make([]WebhookHandler, 0, len(handlers))
	for _, h := range handlers {
		if fmt.Sprintf("%p", h) != fmt.Sprintf("%p", handler) {
			filtered = append(filtered, h)
		}
	}
	m.handlers[event] = filtered
}

func (m *Manager) Emit(event EventType, payload interface{}) {
	m.mu.RLock()
	handlers := m.handlers[event]
	m.mu.RUnlock()

	for _, handler := range handlers {
		func(h WebhookHandler) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[mpesa-sdk] Webhook handler panic: %v\n", r)
				}
			}()
			h(event, payload)
		}(handler)
	}
}

func (m *Manager) HandleSTKCallback(body json.RawMessage) {
	var payload types.STKCallbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("[mpesa-sdk] Failed to parse STK callback: %v\n", err)
		return
	}
	result := client.ParseSTKCallback(payload)
	m.Emit(EventSTKCallback, result)
}

func (m *Manager) HandleResultCallback(body json.RawMessage) {
	var result types.MpesaResult
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[mpesa-sdk] Failed to parse result callback: %v\n", err)
		return
	}

	switch {
	case result.Result.ResultParameters != nil:
		for _, p := range result.Result.ResultParameters.ResultParameter {
			if p.Key == "AccountBalance" {
				m.Emit(EventAccountBalance, result)
				return
			}
			if p.Key == "TransactionStatus" {
				m.Emit(EventTransactionStatus, result)
				return
			}
		}
		m.Emit(EventB2CResult, result)
	default:
		m.Emit(EventB2CResult, result)
	}
}

func VerifySignature(payload []byte, signature, secret string) bool {
	return client.VerifySignature(string(payload), signature, secret)
}
