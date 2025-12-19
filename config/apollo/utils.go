package apollo

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

func format(ns string) string {
	arr := strings.Split(ns, ".")
	suffix := arr[len(arr)-1]
	if len(arr) <= 1 || suffix == properties {
		return json
	}
	if _, ok := formats[suffix]; !ok {
		// fallback
		return json
	}

	return suffix
}

// resolve convert kv pair into one map[string]interface{} by split key into different
// map level. such as: app.name = "application" => map[app][name] = "application"
func resolve(key string, value any, target map[string]any) {
	// expand key "aaa.bbb" into map[aaa]map[bbb]interface{}
	keys := strings.Split(key, ".")
	last := len(keys) - 1
	cursor := target

	for i, k := range keys {
		if i == last {
			cursor[k] = value
			break
		}

		// not the last key, be deeper
		v, ok := cursor[k]
		if !ok {
			// create a new map
			deeper := make(map[string]any)
			cursor[k] = deeper
			cursor = deeper
			continue
		}

		// current exists, then check existing value type, if it's not map
		// that means duplicate keys, and at least one is not map instance.
		if cursor, ok = v.(map[string]any); !ok {
			log.Warnf("duplicate key: %v\n", strings.Join(keys[:i+1], "."))
			break
		}
	}
}

// genKey got the key of config.KeyValue pair.
// eg: namespace.ext with subKey got namespace.subKey
func genKey(ns, sub string) string {
	arr := strings.Split(ns, ".")
	if len(arr) == 1 {
		if ns == "" {
			return sub
		}

		return ns + "." + sub
	}

	suffix := arr[len(arr)-1]
	_, ok := formats[suffix]
	if ok {
		return strings.Join(arr[:len(arr)-1], ".") + "." + sub
	}

	return ns + "." + sub
}
