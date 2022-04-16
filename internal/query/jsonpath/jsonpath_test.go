package jsonpath_test

import (
	"testing"

	"github.com/porter-dev/switchboard/internal/query/jsonpath"

	"github.com/stretchr/testify/assert"
)

func TestNestedJSONPath(t *testing.T) {
	res, err := jsonpath.GetResult(map[string]interface{}{
		"nest-1": map[string]interface{}{
			"nest-2": map[string]string{
				"nest-3": "resultvalue",
			},
		},
	}, "{ .nest-1.nest-2.nest-3 }")

	assert.NoError(t, err, "nested JSON path should not throw error")

	assert.Equal(t, "resultvalue", res, "nested JSON path works as expected")
}

func TestJSONPathConcat(t *testing.T) {
	res, err := jsonpath.GetResult(map[string]interface{}{
		"nest-1": map[string]interface{}{
			"nest-2": map[string]string{
				"nest-3": "resultvalue",
			},
		},
	}, "{ .nest-1.nest-2.nest-3 }{'testing'}")

	assert.NoError(t, err, "nested JSON path should not throw error")

	assert.Equal(t, "resultvaluetesting", res, "concat JSON path works as expected")
}
