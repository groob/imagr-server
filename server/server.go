package server

import "net/http"

func Serve(repoPath string) {
	http.Handle("/", http.FileServer(http.Dir(repoPath)))
	http.ListenAndServe(":3000", nil)
}
