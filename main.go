package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServerInfo struct {
	serverURL []string
	mu        sync.Mutex
	count     int
}

func newServerInfo() *ServerInfo {
	return &ServerInfo{
		serverURL: []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8083", "http://localhost:8084"},
		mu:        sync.Mutex{},
		count:     0,
	}
}

func (s *ServerInfo) getNextServer() string {
	var serverURL string
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == len(s.serverURL)-1 {
		s.count = 0
		serverURL = s.serverURL[s.count]
	} else {
		serverURL = s.serverURL[s.count]
		s.count += 1
	}
	return serverURL
}

func (s *ServerInfo) HandlerLoadBalancer(w http.ResponseWriter, r *http.Request) {
	serverUrl := s.getNextServer()
	request, err := makeProxyRequest(r, serverUrl)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	err = writeResponse(w, response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	serverInfo := newServerInfo()
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(serverInfo.HandlerLoadBalancer))
	log.Println("Load Balancer listening on port: 7777")
	if err := http.ListenAndServe(":7777", router); err != nil {
		log.Fatal(err)
	}

}

func makeProxyRequest(req *http.Request, url string) (*http.Request, error) {
	proxyReq, err := http.NewRequest(req.Method, url+req.RequestURI, req.Body)
	if err != nil {
		return nil, err
	}
	
	proxyReq.Header = req.Header.Clone()
	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)
	return proxyReq, nil
}

func writeResponse(w http.ResponseWriter, resp *http.Response) error {
	for name, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(name, v)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
