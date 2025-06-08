package presenter

import (
	"net/http"
)

func APIResponse(statusCode int, message string) map[string]string {
	var apiResponse = make(map[string]string)
	apiResponse["message"] = message
	apiResponse["status"] = http.StatusText(statusCode)
	return apiResponse
}

func ListAPIResponse[T any](list []T, offset int, limit int) map[string]any {
	var apiResponse = make(map[string]any)
	apiResponse["items"] = list
	apiResponse["offset"] = offset
	apiResponse["limit"] = limit

	if len(list) == limit {
		apiResponse["next_offset"] = offset + limit
	} else {
		apiResponse["next_offset"] = nil
	}
	if (offset - limit) < 0 {
		apiResponse["previous_offset"] = nil
	} else {
		apiResponse["previous_offset"] = offset - limit
	}

	return apiResponse
}

func GetAllResponse(list []any, offset int, limit int) map[string]any {
	return nil
}
