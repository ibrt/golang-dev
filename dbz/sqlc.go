package dbz

// SQLConfig describes the SQLC configuration.
type SQLConfig struct {
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
