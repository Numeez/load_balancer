package main

import (
	"fmt"
	"log"
	"net/http"
)

// I simulate multiple server by writing the following code

// Run it separately by commenting the main in some other service or you can make a different package for this and run it

func spawnServer(port string, done chan error) {
	mux := http.NewServeMux()
	mux.Handle("/", customMiddleWare(http.HandlerFunc(hello), port))
	log.Println("Server running on port: ", port[1:])
	if err := http.ListenAndServe(port, mux); err != nil {
		done <- err
	}
}

func customMiddleWare(next http.Handler, port string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(w, "Hello from the server with middleware after recovering from the panic\n")
				fmt.Fprintf(w, "This was the error: %v", err)
			}
		}()
		fmt.Fprintf(w, "Response coming from http://localhost:%s\n", port[1:])
		log.Println("From custom middleware: ", r.Header)
		next.ServeHTTP(w, r)
	})
}
func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello from the server with middleware")

}

// func main() {
// 	ports := []string{":8080", ":8081", ":8082", ":8083", ":8084"}
// 	done := make(chan error, 1)
// 	for _, port := range ports {
// 		go spawnServer(port, done)
// 	}
// 	if err := <-done; err != nil {
// 		log.Fatal(err)
// 	}
// }
