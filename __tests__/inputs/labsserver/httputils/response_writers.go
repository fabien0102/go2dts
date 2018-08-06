package httputils

import (
	"encoding/json"
	"net/http"
)

type jsonErrorMessage struct {
	ErrorMessages []string            `json:"errors"`
	Fields        map[string][]string `json:"fieldErrors"`
}

// ServerError writes a 500 error response, optionally, if debug is true, it will write the error
// message to the response body
func ServerError(w http.ResponseWriter, err error, debug ...bool) {
	statusText := http.StatusText(http.StatusInternalServerError)
	if len(debug) > 0 && debug[0] {
		statusText = err.Error()
	}
	http.Error(w, statusText, http.StatusInternalServerError)
}

// JSONError constructs and writes a standard json error message payload with
// the supplied statuscode.  Any call to `http.Error` can be swapped to `jsonError`
// without any modification. It optionally accepts zero or more `map[string]string`
// that are merged into a `Fields` value in the final payload.
func JSONError(w http.ResponseWriter, message string, statuscode int, fieldMessages ...map[string]error) {
	msg := jsonErrorMessage{
		ErrorMessages: []string{message},
		Fields:        map[string][]string{},
	}

	for _, messages := range fieldMessages {
		for k, v := range messages {
			msg.Fields[k] = append(msg.Fields[k], v.Error())
		}
	}

	err := JSONResponseWithCode(w, statuscode, msg)
	if err != nil {
		http.Error(w, "json serializtion error", http.StatusInternalServerError)
		return
	}

	return
}

// JSONServerError returns a 500 error server error with a json body of the form
//
// {"message": statusText}
//
// where the statusText is the err.Error() string when debug == true
// and when debug == false it is the standard "Internal Server Error"
func JSONServerError(w http.ResponseWriter, err error, debug ...bool) {
	statusText := http.StatusText(http.StatusInternalServerError)
	if len(debug) > 0 && debug[0] {
		statusText = err.Error()
	}
	JSONError(w, statusText, http.StatusInternalServerError)
}

// JSONResponse is a simple http response writer that only requires arguments
// from the stdlib, it writes the data as json to the response body with a 200 status code
func JSONResponse(w http.ResponseWriter, data interface{}) error {
	return JSONResponseWithCode(w, http.StatusOK, data)
}

// JSONResponseWithCode is a simple http response writer that writes a JSON a payload with the
// specified status code.
func JSONResponseWithCode(w http.ResponseWriter, statuscode int, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	_, err = w.Write(payload)
	return err
}
