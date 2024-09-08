package llm

type UserMessage struct {
	Content string `json:"content"`
}

type DmResponseMessage struct {
	UserInput          string `json:"user-input"`
	RawAdjudicate      string `json:"raw-adjudicate"`
	AdjudicateThoughts string `json:"adjudicate-thoughts"`
	Description        string `json:"description"`
	RawActionEncode    string `json:"raw-action-encode"`
	Actions            string `json:"actions"`
}
