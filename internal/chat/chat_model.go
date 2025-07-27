package chat

type ChatMessage struct {
	Sender  string `json:"sender" validate:"required,alphanum"`
	Message string `json:"message" validate:"required,min=1,max=250"`
}
