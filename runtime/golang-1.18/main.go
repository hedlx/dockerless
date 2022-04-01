package dockerless

import (
	"fmt"
	"net/http"
)

type LambdaT func(w http.ResponseWriter, req *http.Request)

func Lambda(lambda LambdaT) {
	http.HandleFunc("", lambda)
	http.ListenAndServe(":3000", nil)
}
