package nettest

type Result struct {
	SoftwareFrom string `json:"softfrom"`
	SoftwareTo   string `json:"softto"`
	IpFrom       string `json:"ipfrom"`
	IpTo         string `json:"ipto"`
	Port         int    `json:"port"`
	Status       bool   `json:"status"`
}

type Results struct {
	TaskID    string    `json: "taskid"`
	AllStatus []*Result `json: "allstatus"`
}
