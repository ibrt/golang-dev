package dbz

import (
	"fmt"
	"path/filepath"

	"github.com/ibrt/golang-utils/filez"
	"github.com/ibrt/golang-utils/hashz"
	"github.com/ibrt/golang-utils/jsonz"

	"github.com/ibrt/golang-dev/dbz/internal/assets"
	"github.com/ibrt/golang-dev/gtz"
)

// SQLCConfig describes the SQLC configuration.
type SQLCConfig struct {
	Version string              `json:"version,omitempty"`
	Plugins []*SQLCConfigPlugin `json:"plugins,omitempty"`
	SQL     []*SQLCConfigSQL    `json:"sql,omitempty"`
}

// SQLCConfigPlugin describes part of the SQLC configuration.
type SQLCConfigPlugin struct {
	Name string                `json:"name,omitempty"`
	WASM *SQLCConfigPluginWASM `json:"wasm,omitempty"`
}

// SQLCConfigPluginWASM describes part of the SQLC configuration.
type SQLCConfigPluginWASM struct {
	URL    string `json:"url,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

// SQLCConfigSQL describes part of the SQLC configuration.
type SQLCConfigSQL struct {
	Engine   string                  `json:"engine,omitempty"`
	Schema   string                  `json:"schema,omitempty"`
	Queries  string                  `json:"queries,omitempty"`
	Database *SQLCConfigSQLDatabase  `json:"database,omitempty"`
	Codegen  []*SQLCConfigSQLCodegen `json:"codegen,omitempty"`
	Rules    []string                `json:"rules"`
}

// SQLCConfigSQLDatabase describes part of the SQLC configuration.
type SQLCConfigSQLDatabase struct {
	Managed bool   `json:"managed"`
	URI     string `json:"uri,omitempty"`
}

// SQLCConfigSQLCodegen describes part of the SQLC configuration.
type SQLCConfigSQLCodegen struct {
	Plugin  string                       `json:"plugin,omitempty"`
	Out     string                       `json:"out,omitempty"`
	Options *SQLCConfigSQLCodegenOptions `json:"options,omitempty"`
}

// SQLCConfigSQLCodegenOptions describes part of the SQLC configuration.
type SQLCConfigSQLCodegenOptions struct {
	BuildTags                   string                                 `json:"build_tags,omitempty"`
	EmitAllEnumValues           bool                                   `json:"emit_all_enum_values"`
	EmitDbTags                  bool                                   `json:"emit_db_tags"`
	EmitEmptySlices             bool                                   `json:"emit_empty_slices"`
	EmitEnumValidMethod         bool                                   `json:"emit_enum_valid_method"`
	EmitExactTableNames         bool                                   `json:"emit_exact_table_names"`
	EmitExportedQueries         bool                                   `json:"emit_exported_queries"`
	EmitInterface               bool                                   `json:"emit_interface"`
	EmitJSONTags                bool                                   `json:"emit_json_tags"`
	EmitMethodsWithDBArgument   bool                                   `json:"emit_methods_with_db_argument"`
	EmitParamsStructPointers    bool                                   `json:"emit_params_struct_pointers"`
	EmitPointersForNullTypes    bool                                   `json:"emit_pointers_for_null_types"`
	EmitPreparedQueries         bool                                   `json:"emit_prepared_queries"`
	EmitResultStructPointers    bool                                   `json:"emit_result_struct_pointers"`
	EmitSQLAsComment            bool                                   `json:"emit_sql_as_comment"`
	InflectionExcludeTableNames []string                               `json:"inflection_exclude_table_names,omitempty"`
	JSONTagsCaseStyle           string                                 `json:"json_tags_case_style,omitempty"`
	JSONTagsIDUppercase         bool                                   `json:"json_tags_id_uppercase"`
	OmitSQLCVersion             bool                                   `json:"omit_sqlc_version"`
	OmitUnusedStructs           bool                                   `json:"omit_unused_structs"`
	OutputBatchFileName         string                                 `json:"output_batch_file_name,omitempty"`
	OutputCopyFromFileName      string                                 `json:"output_copyfrom_file_name,omitempty"`
	OutputDBFileName            string                                 `json:"output_db_file_name,omitempty"`
	OutputFilesSuffix           string                                 `json:"output_files_suffix,omitempty"`
	OutputModelsFileName        string                                 `json:"output_models_file_name,omitempty"`
	OutputQuerierFileName       string                                 `json:"output_querier_file_name,omitempty"`
	Overrides                   []*SQLCConfigSQLCodegenOptionsOverride `json:"overrides,omitempty"`
	Package                     string                                 `json:"package,omitempty"`
	QueryParameterLimit         int32                                  `json:"query_parameter_limit"`
	Rename                      map[string]string                      `json:"rename,omitempty"`
	SQLDriver                   string                                 `json:"sql_driver,omitempty"`
	SQLPackage                  string                                 `json:"sql_package,omitempty"`
}

// SQLCConfigSQLCodegenOptionsOverride describes part of the SQLC configuration.
type SQLCConfigSQLCodegenOptionsOverride struct {
	Column      string                                     `json:"column,omitempty"`
	DBType      string                                     `json:"db_type,omitempty"`
	Nullable    bool                                       `json:"nullable"`
	Unsigned    bool                                       `json:"unsigned"`
	GoType      *SQLCConfigSQLCodegenOptionsOverrideGoType `json:"go_type,omitempty"`
	GoStructTag string                                     `json:"go_struct_tag,omitempty"`
}

// SQLCConfigSQLCodegenOptionsOverrideGoType describes part of the SQLC configuration.
type SQLCConfigSQLCodegenOptionsOverrideGoType struct {
	Import  string `json:"import,omitempty"`
	Package string `json:"package,omitempty"`
	Type    string `json:"type,omitempty"`
	Pointer bool   `json:"pointer"`
	Slice   bool   `json:"slice"`
}

// SQLCGeneratorParams describes parameters.
type SQLCGeneratorParams struct {
	BuildDirPath   string `validate:"required"`
	SchemaDirPath  string `validate:"required"`
	QueriesDirPath string `validate:"required"`
	PostgresURL    string `validate:"required"`
	OutDirPath     string `validate:"required"`
	OutPackageName string `validate:"required"`
}

// GetPluginFilePath returns the plugin file path.
func (p *SQLCGeneratorParams) GetPluginFilePath() string {
	return filepath.Join(p.BuildDirPath, "plugin.wasm.gz")
}

// GetConfigFilePath returns the config file path.
func (p *SQLCGeneratorParams) GetConfigFilePath() string {
	return filepath.Join(p.BuildDirPath, "plugin.wasm.gz")

}

// SQLCGenerator generates database bindings with SQLC.
type SQLCGenerator struct {
	params *SQLCGeneratorParams
	config *SQLCConfig
}

// MustNewSQLCGenerator initializes a new SQLCGenerator.
func MustNewSQLCGenerator(params *SQLCGeneratorParams) *SQLCGenerator {
	return &SQLCGenerator{
		params: params,
		config: &SQLCConfig{
			Version: "2",
			Plugins: []*SQLCConfigPlugin{
				{
					Name: "plugin",
					WASM: &SQLCConfigPluginWASM{
						URL:    fmt.Sprintf("file://%v", params.GetPluginFilePath()),
						SHA256: hashz.MustHashSHA256(filez.MustReadFile(params.GetPluginFilePath())),
					},
				},
			},
			SQL: []*SQLCConfigSQL{
				{
					Engine:  "postgresql",
					Schema:  params.SchemaDirPath,
					Queries: params.QueriesDirPath,
					Database: &SQLCConfigSQLDatabase{
						Managed: false,
						URI:     params.PostgresURL,
					},
					Codegen: []*SQLCConfigSQLCodegen{
						{
							Plugin: "plugin",
							Out:    params.OutDirPath,
							Options: &SQLCConfigSQLCodegenOptions{
								BuildTags:                "",
								EmitAllEnumValues:        true,
								EmitEnumValidMethod:      true,
								EmitExactTableNames:      true,
								EmitInterface:            true,
								EmitJSONTags:             true,
								EmitParamsStructPointers: true,
								EmitPointersForNullTypes: true,
								EmitResultStructPointers: true,
								EmitSQLAsComment:         true,
								JSONTagsCaseStyle:        "camel",
								OmitSQLCVersion:          true,
								OutputBatchFileName:      "batch.gen.go",
								OutputCopyFromFileName:   "copyfrom.gen.go",
								OutputDBFileName:         "impl.gen.go",
								OutputFilesSuffix:        ".gen",
								OutputModelsFileName:     "models.gen.go",
								OutputQuerierFileName:    "iface.gen.go",
								Package:                  params.OutPackageName,
								QueryParameterLimit:      0,
								Rename:                   make(map[string]string, 0),
								SQLPackage:               "pgx/v5",
								Overrides: []*SQLCConfigSQLCodegenOptionsOverride{
									{
										DBType:   "uuid",
										Nullable: false,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Type: "string",
										},
									},
									{
										DBType:   "uuid",
										Nullable: true,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Type:    "string",
											Pointer: true,
										},
									},
									{
										DBType:   "pg_catalog.timestamptz",
										Nullable: false,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Import: "time",
											Type:   "Time",
										},
									},
									{
										DBType:   "pg_catalog.timestamptz",
										Nullable: true,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Import:  "time",
											Type:    "Time",
											Pointer: true,
										},
									},
									{
										DBType:   "timestamptz",
										Nullable: false,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Import: "time",
											Type:   "Time",
										},
									},
									{
										DBType:   "timestamptz",
										Nullable: true,
										GoType: &SQLCConfigSQLCodegenOptionsOverrideGoType{
											Import:  "time",
											Type:    "Time",
											Pointer: true,
										},
									},
								},
							},
						},
					},
					Rules: []string{
						"sqlc/db-prepare",
					},
				},
			},
		},
	}
}

// SetRename sets a rename rule.
func (c *SQLCGenerator) SetRename(k, v string) *SQLCGenerator {
	c.config.SQL[0].Codegen[0].Options.Rename[k] = v
	return c
}

// MergeRenames merges the given rename rules.
func (c *SQLCGenerator) MergeRenames(m map[string]string) *SQLCGenerator {
	for k, v := range m {
		c.SetRename(k, v)
	}
	return c
}

// AddOverride adds an override to the config.
func (c *SQLCGenerator) AddOverride(overrides ...*SQLCConfigSQLCodegenOptionsOverride) *SQLCGenerator {
	c.config.SQL[0].Codegen[0].Options.Overrides = append(c.config.SQL[0].Codegen[0].Options.Overrides, overrides...)
	return c
}

// MustOutput the config to disk.
func (c *SQLCGenerator) MustOutput() {
	filez.MustWriteFile(c.params.GetConfigFilePath(), 0777, 0666, jsonz.MustMarshalPretty(c.config))
}

// MustGenerate generates the SQLC bindings.
func (c *SQLCGenerator) MustGenerate() {
	filez.MustPrepareDir(c.params.BuildDirPath, 0777)

	filez.MustWriteFile(c.params.GetPluginFilePath(), 0777, 0666, assets.SQLCPluginWASMGZEmbed)
	filez.MustWriteFile(c.params.GetConfigFilePath(), 0777, 0666, jsonz.MustMarshalPretty(c.config))

	gtz.GoToolSQLC.MustRun("vet", "-f", c.params.GetConfigFilePath())
	gtz.GoToolSQLC.MustRun("generate", "-f", c.params.GetConfigFilePath())
}
