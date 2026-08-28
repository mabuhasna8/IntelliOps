package contracts

import "time"

type AuditEvent struct {
    SchemaVersion string `json:"schema_version"`

    EventID   string    `json:"event_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`

    ActorID      string `json:"actor_id,omitempty"`
    ResourceType string `json:"resource_type"`
    ResourceID   string `json:"resource_id"`

    Metadata map[string]any `json:"metadata,omitempty"`
}

