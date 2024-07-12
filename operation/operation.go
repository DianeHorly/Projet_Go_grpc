package operation

// Operation représente une opération effectuée par un appareil
type Operation struct {
	Type string `json:"type"`
	HasSucceeded bool `json:"has_succeeded"`
}

// Device représente un appareil
type Device struct {
    Name             string `json:"device_name"`
    TotalOperations int    `json:"-"`
    FailedOperations int  `json:"-"`
    Operations     []Operation `json:"operations"`
}
