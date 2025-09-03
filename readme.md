# Go HTTP Load Balancer

This is a simple **round-robin HTTP load balancer** implemented in Go.  
It forwards incoming requests to a pool of backend servers and returns the response back to the client.

---

## ✅ Features

- Round-robin load balancing across multiple servers.  
- Forwards HTTP method, headers, and body to backend.  
- Copies status, headers, and body back to client.  
- Thread-safe server rotation using `sync.Mutex`.  
- Backend request timeout support.  

---

## 📂 How It Works

1. Clients send requests to the load balancer (listening on port `7777`).  
2. The load balancer selects the **next backend server** in round-robin fashion.  
3. The request is proxied to that server.  
4. The response from the backend is written back to the client.  

---

## 🚀 Example

```go
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
