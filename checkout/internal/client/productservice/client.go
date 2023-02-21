package productservice

const (
	GetProductEndpoint = "/get_product"
	ListSKUEndpoint    = "/list_skus"
)

type Client struct {
	url        string
	urlProduct string
	urlList    string
	token      string
}

func New(url, token string) *Client {
	return &Client{
		url:        url,
		urlProduct: url + GetProductEndpoint,
		urlList:    url + ListSKUEndpoint,
		token:      token,
	}
}
