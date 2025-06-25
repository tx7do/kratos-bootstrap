package elasticsearch

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
)

// ParseErrorMessage 解析 Elasticsearch 错误消息
func ParseErrorMessage(body io.ReadCloser) (*ErrorResponse, error) {
	defer body.Close()

	var errorResponse ErrorResponse
	if err := json.NewDecoder(body).Decode(&errorResponse); err != nil {
		return nil, ErrUnmarshalResponse
	}

	return &errorResponse, nil
}

// MergeOptions 合并 Elasticsearch 索引的映射和设置
func MergeOptions(mapping, settings string) (string, error) {
	codec := encoding.GetCodec("json")

	body := make(map[string]interface{})

	if mapping != "" {
		var mappingObj map[string]interface{}
		if err := codec.Unmarshal([]byte(mapping), &mappingObj); err != nil {
			log.Errorf("failed to unmarshal mapping: %v", err)
			return "", err
		}
		if existingMappings, ok := mappingObj["mappings"]; ok {
			body["mappings"] = existingMappings
		} else {
			body["mappings"] = mappingObj
		}
	}

	if settings != "" {
		var settingsObj map[string]interface{}
		if err := codec.Unmarshal([]byte(settings), &settingsObj); err != nil {
			log.Errorf("failed to unmarshal settings: %v", err)
			return "", err
		}
		// 检查 settings 是否包含 settings 字段
		if existingSettings, ok := settingsObj["settings"]; ok {
			body["settings"] = existingSettings
		} else {
			body["settings"] = settingsObj
		}
	}

	bodyBytes, err := codec.Marshal(body)
	if err != nil {
		log.Errorf("failed to marshal request body: %v", err)
		return "", err
	}

	return string(bodyBytes), nil
}

func ParseQueryString(query string) []string {
	codec := encoding.GetCodec("json")

	var err error
	queryMap := make(map[string]string)
	if err = codec.Unmarshal([]byte(query), &queryMap); err == nil {
		var queries []string
		for k, v := range queryMap {
			queries = append(queries, k+":"+v)
		}
		return queries
	}

	var queryMapArray []map[string]string
	if err = codec.Unmarshal([]byte(query), &queryMapArray); err == nil {
		var queries []string
		for _, item := range queryMapArray {
			for k, v := range item {
				queries = append(queries, k+":"+v)
			}
		}
		return queries
	}

	return nil
}

func MakeQueryString(andQuery, orQuery string) string {
	a := ParseQueryString(andQuery)
	o := ParseQueryString(orQuery)

	if len(a) == 0 && len(o) == 0 {
		return ""
	}

	if len(a) > 0 && len(o) == 0 {
		return strings.Join(a, " AND ")
	} else if len(a) == 0 && len(o) > 0 {
		return strings.Join(o, " OR ")
	} else if len(a) > 0 && len(o) > 0 {
		return strings.Join(a, " AND ") + " AND (" + strings.Join(o, " OR ") + ")"
	} else {
		return strings.Join(a, " AND ") + " AND " + strings.Join(o, " OR ")
	}
}
