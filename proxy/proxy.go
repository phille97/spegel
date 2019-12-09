package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"
)

type Proxy struct {
	client *http.Client
	Nodes  []url.URL
}

func NewProxy() *Proxy {
	client := &http.Client{}

	return &Proxy{
		client: client,
		Nodes:  []url.URL{},
	}
}

func (p *Proxy) Update(nodes []url.URL) {
	p.Nodes = nodes
}

func (p *Proxy) HandleProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	for _, node := range p.Nodes {
		nodeUrl := node
		nodeUrl.Path = path.Join("/get/", vars["path"])
		resp, err := p.client.Head(nodeUrl.String())
		if err == nil {
			if r.Method == "HEAD" {
				w.WriteHeader(resp.StatusCode)
				return
			} else if r.Method == "GET" && resp.StatusCode == http.StatusOK {
				http.Redirect(w, r, nodeUrl.String(), http.StatusFound)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "no backends availible")
	return
}
