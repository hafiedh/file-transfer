package constants

const (
	MESSAGE_SUCCESS                = "Success"
	MESSAGE_PARTIAL_SUCCESS        = "Partial Success"
	MESSAGE_STILL_PROCESS          = "is being process"
	MESSAGE_FAILED                 = "Something went wrong"
	MESSAGE_INVALID_REQUEST_FORMAT = "Invalid Request Format"
	MESSAGE_UNAUTHORIZED           = "Unauthorized"
	MESSAGE_FORBIDDEN              = "Forbidden"
	MESSAGE_CONFLICT               = "Conflict"
	MESSAGE_PROGRESSION            = "Progression"
)

type DefaultResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Status  int         `json:"status"`
}

type PaginationData struct {
	Page        uint `json:"page"`
	TotalPages  uint `json:"totalPages"`
	TotalItems  uint `json:"totalItems"`
	Limit       uint `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

type PaginationResponseData struct {
	Results        interface{} `json:"results"`
	PaginationData `json:"pagination"`
}
