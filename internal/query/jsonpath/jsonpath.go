package jsonpath

import (
	"fmt"

	"k8s.io/client-go/util/jsonpath"
)

func GetResult(data map[string]interface{}, query string) (interface{}, error) {
	js := jsonpath.New("query")

	err := js.Parse(query)

	if err != nil {
		return nil, err
	}

	results, err := js.FindResults(data)

	if err != nil {
		return nil, err
	}

	validResults := make([]interface{}, 0)

	for _, result := range results {
		for _, r := range result {
			// if this cannot be interfaced, throw an error
			if r.CanInterface() {
				validResults = append(validResults, r.Interface())
			}
		}
	}

	if len(validResults) == 1 {
		return validResults[0], nil
	} else if len(validResults) > 1 {
		var joinedRes string

		for _, validRes := range validResults {
			joinedRes += fmt.Sprintf("%v", validRes)
		}

		return joinedRes, nil
	}

	return nil, fmt.Errorf("no query result")
}
