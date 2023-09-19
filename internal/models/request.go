package models

type Request struct {
	Id      int
	Scheme  string
	Host    string
	Path    string
	Method  string
	Headers string
	Body    string
	Params  string
	Cookies string
}
