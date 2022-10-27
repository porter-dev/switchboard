package types

type Variable struct {
	Name   *YAMLNode[string] `yaml:"name" validate:"required,unique"`
	Value  *YAMLNode[string] `yaml:"value" validate:"required_if=Random false"`
	Once   *YAMLNode[bool]   `yaml:"once"`
	Random *YAMLNode[bool]   `yaml:"random"`
	Length *YAMLNode[uint]   `yaml:"length"`
}

type EnvGroup struct {
	Name      *YAMLNode[string] `yaml:"name" validate:"required"`
	CloneFrom *YAMLNode[string] `yaml:"clone_from" validate:"required"`
}

type Build struct {
	Name       *YAMLNode[string]                       `yaml:"name" validate:"required"`
	Context    *YAMLNode[string]                       `yaml:"context" validate:"dir"`
	Method     *YAMLNode[string]                       `yaml:"method" validate:"required,oneof=pack docker registry"`
	Builder    *YAMLNode[string]                       `yaml:"builder" validate:"required_if=Method pack"`
	Buildpacks []*YAMLNode[string]                     `yaml:"buildpacks"`
	Dockerfile *YAMLNode[string]                       `yaml:"dockerfile" validate:"required_if=Method docker"`
	Image      *YAMLNode[string]                       `yaml:"image" validate:"required_if=Method registry"`
	Env        map[*YAMLNode[string]]*YAMLNode[string] `yaml:"env"`
	EnvGroups  []*YAMLNode[string]                     `yaml:"env_groups"`
}

type AddonResource struct {
	Name      *YAMLNode[string]                    `yaml:"name" validate:"required,unique"`
	DependsOn []*YAMLNode[string]                  `yaml:"depends_on"`
	Type      *YAMLNode[string]                    `yaml:"type" validate:"required"`
	ChartURL  *YAMLNode[string]                    `yaml:"chart_url" validate:"url"`
	Version   *YAMLNode[string]                    `yaml:"version"`
	Deploy    map[*YAMLNode[string]]*YAMLNode[any] `yaml:"deploy"`
}

// type AppResourceBuild struct {
// 	Name       string         `yaml:"name" mapstructure:"name"`
// 	Context    string         `yaml:"context" mapstructure:"context"`
// 	Method     string         `yaml:"method" mapstructure:"method"`
// 	Builder    string         `yaml:"builder" mapstructure:"builder"`
// 	Buildpacks []string       `yaml:"buildpacks" mapstructure:"buildpacks"`
// 	Dockerfile string         `yaml:"dockerfile" mapstructure:"dockerfile"`
// 	Image      string         `yaml:"image" mapstructure:"image"`
// 	Env        map[string]any `yaml:"env" mapstructure:"env"`
// 	EnvGroups  []string       `yaml:"env_groups" mapstructure:"env_groups"`
// 	Ref        string         `yaml:"ref" mapstructure:"ref"`
// }

type AppResource struct {
	Name      *YAMLNode[string]                    `yaml:"name" validate:"required,unique"`
	DependsOn []*YAMLNode[string]                  `yaml:"depends_on"`
	Type      *YAMLNode[string]                    `yaml:"type" validate:"required,oneof=web worker job"`
	ChartURL  *YAMLNode[string]                    `yaml:"chart_url" validate:"url"`
	Version   *YAMLNode[string]                    `yaml:"version"`
	Deploy    map[*YAMLNode[string]]*YAMLNode[any] `yaml:"deploy"`
	Build     map[*YAMLNode[string]]*YAMLNode[any] `yaml:"build"`
}

type PorterYAML struct {
	Version   *YAMLNode[string]           `yaml:"version" validate:"required"`
	Variables *YAMLNode[[]*Variable]      `yaml:"variables"`
	EnvGroups *YAMLNode[[]*EnvGroup]      `yaml:"env_groups"`
	Builds    *YAMLNode[[]*Build]         `yaml:"builds"`
	Apps      *YAMLNode[[]*AppResource]   `yaml:"apps"`
	Addons    *YAMLNode[[]*AddonResource] `yaml:"addons"`
}
