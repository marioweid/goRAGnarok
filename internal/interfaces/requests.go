package interfaces

type GenerateRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}
