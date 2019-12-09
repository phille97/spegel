package proxy

import (
	"fmt"
	"io"
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
		headResp, err := p.client.Head(nodeUrl.String())
		if err == nil && headResp.StatusCode == http.StatusOK {
			req, err := http.NewRequest(r.Method, nodeUrl.String(), r.Body)
			if err != nil {
				continue
			}
			copyHeaders(req.Header, r.Header)

			resp, err := p.client.Do(req)
			if err != nil {
				continue
			}

			copyHeaders(w.Header(), resp.Header)
			w.WriteHeader(resp.StatusCode)

			io.Copy(w, resp.Body)
			resp.Body.Close()

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "no backends availible")
	return
}

func copyHeaders(dst, src http.Header) {
	for k := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
