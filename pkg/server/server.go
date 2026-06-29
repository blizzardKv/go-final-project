package server

import (
	"fmt"
	"net/http"
	"os"
)

const defaultPort = 7540

func port() int {
	if env := os.Getenv("TODO_PORT"); env != "" {
		var p int
		if _, err := fmt.Sscanf(env, "%d", &p); err == nil {
			return p
		}
	}
	return defaultPort
}

func Start(webDir string) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	addr := fmt.Sprintf(":%d", port())
	return http.ListenAndServe(addr, nil)
}
