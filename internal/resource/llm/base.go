package llm

type LLM interface {
	Invoke(prompt string) (string, error)
	Stream(prompt string) (chan string, chan error, error)
}
