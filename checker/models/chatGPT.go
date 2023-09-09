package models

type BodyToRequest struct {
	Model            string `json:"model"`
	Messages         []Msg  `json:"messages"`
	Temperature      int    `json:"temperature"`
	MaxTokens        int    `json:"max_tokens"`
	TopP             int    `json:"top_p"`
	FrequencyPenalty int    `json:"frequency_penalty"`
	PresencePenalty  int    `json:"presence_penalty"`
}

type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToParseAnswer struct {
	Choices []MsgAnswer `json:"choices"`
}

type MsgAnswer struct {
	MSG MSG `json:"message"`
}

type MSG struct {
	Content string `json:"content"`
}

type ToAddInDB struct {
	Feeling string `json:"feeling"`
	Rate    string `json:"rate"`
}
