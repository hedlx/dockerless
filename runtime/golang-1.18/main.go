package doless

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Method  string `json:"method"`
	Payload []byte `json:"payload"`
}

type Error struct {
	Reason  string  `json:"reason"`
	Details *string `json:"details,omitempty"`
}

type LambdaT func(req *Request) (int, string)

func Handler(lambda LambdaT) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			rec := recover()

			if rec == nil {
				return
			}

			ret := &Error{
				Reason: "Failed to handle request",
			}

			if err, isError := rec.(error); isError {
				details := err.Error()
				ret.Details = &details
			}

			if sErr, isString := rec.(string); isString {
				ret.Details = &sErr
			}

			resp, _ := json.Marshal(ret)

			w.WriteHeader(500)
			fmt.Fprint(w, string(resp))
		}()

		defer req.Body.Close()
		payload, err := io.ReadAll(req.Body)

		if err != nil {
			w.WriteHeader(500)
			resp, _ := json.Marshal(&Error{
				Reason: err.Error(),
			})
			fmt.Fprint(w, string(resp))
			return
		}

		status, resp := lambda(&Request{
			Method:  req.Method,
			Payload: payload,
		})

		w.WriteHeader(status)
		fmt.Fprint(w, resp)
	}
}

func Lambda(lambda LambdaT) {
	http.HandleFunc("/", Handler(lambda)) // TODO: must be generic handler
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})
	http.ListenAndServe(":3000", nil)
}
