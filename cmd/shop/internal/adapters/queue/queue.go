package queue

type Client struct {
}

func New() (*Client, error) {
	return &Client{}, nil
}
