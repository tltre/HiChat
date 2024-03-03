package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseBody define the standard response from API
/* the params are:
+ SC: http status code
+ Code: internal code
+ Msg: error message
+ Data: some value return from api
+ Rows: DataBase rows data
+ Total: length of Rows
*/
type ResponseBody struct {
	SC    int
	Code  int
	Msg   string
	Data  map[string]string
	Rows  interface{}
	Total int
}

// SendNormalResp Unified return normal information
func SendNormalResp(w http.ResponseWriter, msg string, data map[string]string, rows interface{}, total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	r := ResponseBody{
		SC:    http.StatusOK,
		Code:  0,
		Msg:   msg,
		Data:  data,
		Rows:  rows,
		Total: total,
	}
	ret, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, string(ret))
}

// SendErrorResp return error message
func SendErrorResp(w http.ResponseWriter, sc int, msg string, data map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)
	r := ResponseBody{
		SC:    sc,
		Code:  -1,
		Msg:   msg,
		Data:  data,
		Rows:  nil,
		Total: 0,
	}
	ret, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, string(ret))
}
