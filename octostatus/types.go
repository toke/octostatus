package main

type JobInfo struct {
	JobDetail JobDetail   `json:"job"`
	Progress  JobProgress `json:"progress"`
	State     string      `json:"state"`
}

type JobDetail struct {
	AveragePrintTime   float64                `json:"averagePrintTime"`
	EstimatedPrintTime float64                `json:"estimatedPrintTime"`
	LastPrintTime      float64                `json:"lastPrintTime"`
	File               JobFile                `json:"file"`
	Tool               map[string]JobFilament `json:"filament"`
}

type JobFile struct {
	Date   int    `json:"date"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Path   string `json:"path"`
	Size   int32  `json:"size"`
}

type JobFilament struct {
	Length float64 `json:"length"`
	Volume float64 `json:"volume"`
}

type JobProgress struct {
	ETA                 string  `json:"ETA,omitempty"`
	Completion          float64 `json:"completion,omitempty"`
	Filepos             int32   `json:"filepos"`
	PrintTime           int32   `json:"printTime"`
	PrintTimeLeft       int32   `json:"printTimeLeft"`
	PrintTimeLeftOrigin string  `json:"printTimeLeftOrigin,omitempty"`
	PrintTimeLeftString string  `json:"printTimeLeftString,omitempty"`
}
