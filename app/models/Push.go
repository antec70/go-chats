package models

type JsonRequest struct {
	Push   *Push  `json:"push"`
	Topic  string `json:"topic" example:"Saint-Petersburg"`
	ChatId string `json:"chatId" example:"305"`
}

type Push struct {
	Type  string `json:"type" example:"1"`
	Image string `json:"image" example:"image.png"`
	Title string `json:"title" example:"New notification"`
	Body  string `json:"body" example:"Check in App"`
}
