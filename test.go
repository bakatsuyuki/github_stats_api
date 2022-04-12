package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
)

func main() {
	fmt.Println("ここ")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://dummy.url.com/user", strings.NewReader(""))

	authGitHubApp(w, r)
	fmt.Println(w.Result())
}
