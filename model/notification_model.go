package model

type NotificationModel struct {
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Token       string            `json:"token,omitempty"`
	Tokens      []string          `json:"tokens,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	ImageURL    string            `json:"imageUrl,omitempty"`
	ClickAction string            `json:"clickAction,omitempty"`
}
