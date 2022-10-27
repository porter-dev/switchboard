package types

type Variable struct {
	Name   *YAMLNode[string] `json:"name" validate:"required,unique"`
	Value  *YAMLNode[string] `json:"value" validate:"required_if=Random false"`
	Once   *YAMLNode[bool]   `json:"once"`
	Random *YAMLNode[bool]   `json:"random"`
	Length *YAMLNode[uint]   `json:"length"`
}

type EnvGroup struct {
	Name      *YAMLNode[string] `json:"name" validate:"required"`
	CloneFrom *YAMLNode[string] `json:"clone_from" validate:"required"`
}

type Build struct {
	Name       *YAMLNode[string]                       `json:"name" validate:"required"`
	Context    *YAMLNode[string]                       `json:"context" validate:"dir"`
	Method     *YAMLNode[string]                       `json:"method" validate:"required,oneof=pack docker registry"`
	Builder    *YAMLNode[string]                       `json:"builder" validate:"required_if=Method pack"`
	Buildpacks []*YAMLNode[string]                     `json:"buildpacks"`
	Dockerfile *YAMLNode[string]                       `json:"dockerfile" validate:"required_if=Method docker"`
	Image      *YAMLNode[string]                       `json:"image" validate:"required_if=Method registry"`
	Env        map[*YAMLNode[string]]*YAMLNode[string] `json:"env"`
	EnvGroups  []*YAMLNode[string]                     `json:"env_groups"`
}

type AddonResource struct {
	Name      *YAMLNode[string]                    `json:"name" validate:"required,unique"`
	DependsOn []*YAMLNode[string]                  `json:"depends_on"`
	Type      *YAMLNode[string]                    `json:"type" validate:"required"`
	ChartURL  *YAMLNode[string]                    `json:"chart_url" validate:"url"`
	Version   *YAMLNode[string]                    `json:"version"`
	Deploy    map[*YAMLNode[string]]*YAMLNode[any] `json:"deploy"`
}

// type AppResourceBuild struct {
// 	Name       string         `json:"name" mapstructure:"name"`
// 	Context    string         `json:"context" mapstructure:"context"`
// 	Method     string         `json:"method" mapstructure:"method"`
// 	Builder    string         `json:"builder" mapstructure:"builder"`
// 	Buildpacks []string       `json:"buildpacks" mapstructure:"buildpacks"`
// 	Dockerfile string         `json:"dockerfile" mapstructure:"dockerfile"`
// 	Image      string         `json:"image" mapstructure:"image"`
// 	Env        map[string]any `json:"env" mapstructure:"env"`
// 	EnvGroups  []string       `json:"env_groups" mapstructure:"env_groups"`
// 	Ref        string         `json:"ref" mapstructure:"ref"`
// }

type AppResource struct {
	Name      *YAMLNode[string]                    `json:"name" validate:"required,unique"`
	DependsOn []*YAMLNode[string]                  `json:"depends_on"`
	Type      *YAMLNode[string]                    `json:"type" validate:"required,oneof=web worker job"`
	ChartURL  *YAMLNode[string]                    `json:"chart_url" validate:"url"`
	Version   *YAMLNode[string]                    `json:"version"`
	Deploy    map[*YAMLNode[string]]*YAMLNode[any] `json:"deploy"`
	Build     map[*YAMLNode[string]]*YAMLNode[any] `json:"build"`
}

type PorterYAML struct {
	Version   *YAMLNode[string]           `json:"version" validate:"required"`
	Variables *YAMLNode[[]*Variable]      `json:"variables"`
	EnvGroups *YAMLNode[[]*EnvGroup]      `json:"env_groups"`
	Builds    *YAMLNode[[]*Build]         `json:"builds"`
	Apps      *YAMLNode[[]*AppResource]   `json:"apps"`
	Addons    *YAMLNode[[]*AddonResource] `json:"addons"`
}
