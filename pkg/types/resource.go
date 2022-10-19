package types

type ResourceGroup struct {
	Version   string      `json:"version"`
	Resources []*Resource `json:"resources"`
}

type Resource struct {
	Name      string                 `json:"name"`
	Driver    string                 `json:"driver,omitempty"`
	Source    map[string]interface{} `json:"source,omitempty"`
	Target    map[string]interface{} `json:"target,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	DependsOn []string               `json:"depends_on,omitempty"`
}
