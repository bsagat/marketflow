package domain

type ConnMsg struct {
	Connection string `json:"connection,omitempty"`
	Status     string `json:"status"`
}
