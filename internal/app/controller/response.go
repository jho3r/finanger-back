package controller

type (
	// Error is the struct for the error response.
	Error struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	// Success is the struct for the success response.
	Success struct {
		Message string `json:"message"`
	}

	// Data is the struct for the data response.
	Data struct {
		Data interface{} `json:"data"`
	}
)
