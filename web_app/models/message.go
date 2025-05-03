package models

// Canal 传入 Kafka的消息结构
type CanalMessage struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	IdDdl    bool                     `json:"idDdl"`
	Data     []map[string]interface{} `json:"data"`
}
