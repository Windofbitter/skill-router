package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Println("Skill Router running at http://localhost:9527")
	http.ListenAndServe(":9527", nil)
}
