package errorx

import "net/http"

func init() {
	RegisterBatch([]CodeSpec{
		{CodeOK, "success", http.StatusOK, "Success", false},
		{CodeInvalidParams, "invalid parameters", http.StatusBadRequest, "Invalid request parameters", false},
		{CodeUnauthorized, "unauthorized", http.StatusUnauthorized, "Authentication required", false},
		{CodeForbidden, "forbidden", http.StatusForbidden, "Access denied", false},
		{CodeNotFound, "not found", http.StatusNotFound, "Resource not found", false},
		{CodeInternal, "internal server error", http.StatusInternalServerError, "Internal server error", true},
		{CodeUnavailable, "service unavailable", http.StatusServiceUnavailable, "Service unavailable", true},
	})
}
