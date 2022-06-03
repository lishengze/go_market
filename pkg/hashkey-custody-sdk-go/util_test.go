package jadepoolsaas

import (
	"encoding/json"
	"net/http"
)

const (
	TestAppKey    = "gznXiKEInYdAiITtk55KyUk3"
	TestAppSecret = "gZJHdgNYlywjdS815T8feXoPfmY9K6KCBRuPs8q3f2tvEWnzN5S58OJjRraY5YQE"
)

func writeSuccessResponse(w http.ResponseWriter, data map[string]interface{}) (int, error) {
	sign, err := signHMACSHA256(data, TestAppSecret)
	if err != nil {
		return 0, err
	}

	result, err := json.Marshal(&map[string]interface{}{
		"code":    0,
		"data":    data,
		"message": "success",
		"sign":    sign,
	})
	if err != nil {
		return 0, err
	}

	return w.Write(result)
}
