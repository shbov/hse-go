package Structs

type Response struct {
	Message string `json:"outputString"`
	Code    int    `json:"-"`
}
