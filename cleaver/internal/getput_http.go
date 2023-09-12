package internal

import (
	"io"
	"net/http"
	"net/url"
	"path"

	"git.softndit.com/collector/backend/cleaver"
	"github.com/BurntSushi/toml"
)

func init() {
	RegisterGetter("http", NewHTTPGetter)
}

type httpGetter struct {
	baseURL url.URL
}

// NewHTTPGetter TBD
func NewHTTPGetter(config *toml.Primitive) (ImageGetter, error) {
	var data struct {
		BaseURL string
	}
	if err := toml.PrimitiveDecode(*config, &data); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(data.BaseURL)
	if err != nil {
		return nil, err
	}

	return &httpGetter{*baseURL}, nil
}

func (h *httpGetter) Get(src url.URL) ([]byte, error) {
	if src.Scheme == "" {
		if h.baseURL.Scheme != "" {
			src.Scheme = h.baseURL.Scheme
		} else {
			src.Scheme = "http"
		}
	}
	if src.Host == "" && h.baseURL.Host != "" {
		src.Host = h.baseURL.Host
	}
	src.Path = path.Join(h.baseURL.Path, src.Path)

	resp, err := http.Get(src.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, cleaver.ErrCategoryGetNotFound.New("not found")
	default:
		return nil, cleaver.ErrCategoryGet.Newf("can't get data, http status", resp.StatusCode)
	}
	return body, nil
}
