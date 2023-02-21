package loms

const (
	StockEndpoint       = "/stocks"
	CreateOrderEndpoint = "/createOrder"
)

type Client struct {
	url      string
	urlStock string
	urlOrder string
}

func New(url string) *Client {
	return &Client{
		url:      url,
		urlStock: url + StockEndpoint,
		urlOrder: url + CreateOrderEndpoint,
	}
}
