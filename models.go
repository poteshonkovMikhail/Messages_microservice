package main

type Message struct {
	ID        int32  `json:"id"`
	Content   string `json:"content"`
	Processed bool   `json:"processed"`
}

type MessageStatistics struct {
	TotalMessages       int `json:"total_messages"`
	ProcessedMessages   int `json:"processed_messages"`
	UnprocessedMessages int `json:"unprocessed_messages"`
}
