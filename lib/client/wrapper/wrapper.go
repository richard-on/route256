package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

func NewRequest[Req any, Resp any](ctx context.Context, url string, method string, req Req, resp Resp) (Resp, error) {
	rawJSON, err := json.Marshal(req)
	if err != nil {
		return resp, errors.Wrap(err, "marshaling json")
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(rawJSON))
	if err != nil {
		return resp, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return resp, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&resp)
	if err != nil {
		return resp, errors.Wrap(err, "decoding json")
	}

	return resp, nil
}
