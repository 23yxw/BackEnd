package utils

func JsonResponse(errorCode int, results interface{}, msg string, extraInfo string) (ret map[string]interface{}) {

	ret = map[string]interface{}{
		"errorCode": errorCode,
		"data":      results,
		"msg":       msg,
		"extraInfo": extraInfo,
	}
	return
}
