package types

type Variable struct {
	Name   *YAMLNodeMetadata[string] `json:"name" validate:"required,unique"`
	Value  *YAMLNodeMetadata[string] `json:"value" validate:"required_if=Random false"`
	Once   *YAMLNodeMetadata[bool]   `json:"once"`
	Random *YAMLNodeMetadata[bool]   `json:"random"`
}

type EnvGroup struct {
	Name      *YAMLNodeMetadata[string] `json:"name" validate:"required"`
	CloneFrom *YAMLNodeMetadata[string] `json:"clone_from" validate:"required"`
}

type Build struct {
	Name       *YAMLNodeMetadata[string]                            `json:"name" validate:"required"`
	Context    *YAMLNodeMetadata[string]                            `json:"context" validate:"dir"`
	Method     *YAMLNodeMetadata[string]                            `json:"method" validate:"required,oneof=pack docker registry"`
	Builder    *YAMLNodeMetadata[string]                            `json:"builder" validate:"required_if=Method pack"`
	Buildpacks []*YAMLNodeMetadata[string]                          `json:"buildpacks"`
	Dockerfile *YAMLNodeMetadata[string]                            `json:"dockerfile" validate:"required_if=Method docker"`
	Image      *YAMLNodeMetadata[string]                            `json:"image" validate:"required_if=Method registry"`
	Env        map[*YAMLNodeMetadata[string]]*YAMLNodeMetadata[any] `json:"env"`
	EnvGroups  []*YAMLNodeMetadata[string]                          `json:"env_groups"`
}

type AddonResource struct {
	Name      *YAMLNodeMetadata[string]                            `json:"name" validate:"required,unique"`
	DependsOn []*YAMLNodeMetadata[string]                          `json:"depends_on"`
	Type      *YAMLNodeMetadata[string]                            `json:"type" validate:"required"`
	ChartURL  *YAMLNodeMetadata[string]                            `json:"chart_url" validate:"url"`
	Version   *YAMLNodeMetadata[string]                            `json:"version"`
	Deploy    map[*YAMLNodeMetadata[string]]*YAMLNodeMetadata[any] `json:"deploy"`
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
	Name      *YAMLNodeMetadata[string]                            `json:"name" validate:"required,unique"`
	DependsOn []*YAMLNodeMetadata[string]                          `json:"depends_on"`
	Type      *YAMLNodeMetadata[string]                            `json:"type" validate:"required,oneof=web worker job"`
	ChartURL  *YAMLNodeMetadata[string]                            `json:"chart_url" validate:"url"`
	Version   *YAMLNodeMetadata[string]                            `json:"version"`
	Deploy    map[*YAMLNodeMetadata[string]]*YAMLNodeMetadata[any] `json:"deploy"`
	Build     map[*YAMLNodeMetadata[string]]*YAMLNodeMetadata[any] `json:"build"`
}

type PorterYAML struct {
	Version   *YAMLNodeMetadata[string] `json:"version" validate:"required"`
	Variables []*Variable               `json:"variables"`
	EnvGroups []*EnvGroup               `json:"env_groups"`
	Builds    []*Build                  `json:"builds"`
	Apps      []*AppResource            `json:"apps"`
	Addons    []*AddonResource          `json:"addons"`
}
