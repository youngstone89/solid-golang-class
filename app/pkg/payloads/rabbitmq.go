package payloads

import (
	"sync"
	"time"
)

var (
	// 컴파일 타임 타입 변경 체크
	_ Payload = (*RabbitMQPayload)(nil)

	rabbitMQPayloadPool = sync.Pool{
		New: func() interface{} { return new(RabbitMQPayload) },
	}
)

type RabbitMQPayload struct {
	Queue     string                 `json:"queue,omitempty"`
	Value     map[string]interface{} `json:"value,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`

	Index string `json:"index,omitempty"`
	DocID string `json:"doc_id,omitempty"`
	Data  []byte `json:"data,omitempty"`
}

// Clone implements pipeline.Payload.
func (kp *RabbitMQPayload) Clone() Payload {
	newP := rabbitMQPayloadPool.Get().(*RabbitMQPayload)

	return newP
}

// Out implements Payload
func (kp *RabbitMQPayload) Out() (string, string, []byte) {
	return kp.Index, kp.DocID, kp.Data
}

// MarkAsProcessed implements pipeline.Payload
func (p *RabbitMQPayload) MarkAsProcessed() {

	rabbitMQPayloadPool.Put(p)
}