package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// API handles api calls
type API struct {
	root string
}

// NewAPI creates a new API, root should be set to the root url for the API
func NewAPI(root string) *API {
	api := &API{
		root: root,
	}
	return api
}

func (api *API) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	endpoint := req.URL.Path[len(api.root):]
	method := req.Method

	var err error
	switch endpoint {
	case "channels":
		switch method {
		case "POST":
			err = api.handlePostChannel(res, req)
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	case "messages":
		switch method {
		case "POST":
			err = api.handlePostMessage(res, req)
		default:
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	default:
		http.NotFound(res, req)
		return
	}

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err.Error())
	}
}

func (api *API) handlePostChannel(res http.ResponseWriter, req *http.Request) error {
	return fmt.Errorf("not implemented")
}

func (api *API) handlePostMessage(res http.ResponseWriter, req *http.Request) error {
	return fmt.Errorf("handle post message not implemented")
}
