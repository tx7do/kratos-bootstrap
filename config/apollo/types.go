package apollo

const (
	yaml       = "yaml"
	yml        = "yml"
	json       = "json"
	properties = "properties"
)

var formats map[string]struct{}

func init() {
	formats = make(map[string]struct{})

	formats[yaml] = struct{}{}
	formats[yml] = struct{}{}
	formats[json] = struct{}{}
	formats[properties] = struct{}{}
}
