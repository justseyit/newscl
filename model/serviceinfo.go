package model



type ServiceState string

const (
	Running ServiceState = "running"
	Waiting ServiceState = "waiting"
)

type ServiceInfo struct {
	SourceCheckFrequency uint `json:"sourceCheckFrequency" bson:"sourceCheckFrequency"`
	SourceInfo []struct {
		SourceID string `json:"sourceID" bson:"sourceID"`
		SourceName string `json:"sourceName" bson:"sourceName"`
		SourceURL string `json:"sourceURL" bson:"sourceURL"`
		SourceFetchFrequency uint `json:"sourceFetchFrequency" bson:"sourceFetchFrequency"`
		LastChecked int `json:"lastChecked" bson:"lastChecked"`
	} `json:"sourceInfo" bson:"sourceInfo"`
	RunningTime int `json:"runningTime" bson:"runningTime"`
	ServiceState ServiceState `json:"serviceState" bson:"serviceState"`
}
