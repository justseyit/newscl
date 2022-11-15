package model

var ServInfo *serviceInfo

type serviceInfo struct {
	LastPostTime   string `json:"lastPostTime"`
	WaitingForPost bool   `json:"waitingForPost"`
	PostInterval   int    `json:"postInterval"` // in seconds
	PostURL        string `json:"postURL"`
	SucceededPost  int    `json:"succeededPost"`
	FailedPost     int    `json:"failedPost"`
}


func InitServiceInfo() {
	ServInfo = &serviceInfo{}
}