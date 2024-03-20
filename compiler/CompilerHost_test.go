package main

import (
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"os"
	"testing"
)

const py_extractor_json = `{	"idl_filename": "py_extractor",	"idl_extension": ".json",	"idl_filename_with_extension": "py_extractor.json",	"idl_full_path": "py_extractor.json",	"metaffi_guest_lib": "py_extractor_MetaFFIGuest",	"target_language": "python3",	"modules": [		{			"name": "py_extractor",			"comment": "",			"tags": {},			"functions": [],			"classes": [				{					"name": "variable_info",					"comment": "",					"tags": {},					"constructors": [],					"release": {						"name": "Releasevariable_info",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_variable_info_Releasevariable_info",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [],					"fields": [						{							"name": "name",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_name",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_variable_info_get_name",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "name",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "type",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_type",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_variable_info_get_type",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "type",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "is_getter",							"type": "bool",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_is_getter",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_variable_info_get_is_getter",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "is_getter",										"type": "bool",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "is_setter",							"type": "bool",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_is_setter",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_variable_info_get_is_setter",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "is_setter",										"type": "bool",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						}					]				},				{					"name": "parameter_info",					"comment": "",					"tags": {},					"constructors": [],					"release": {						"name": "Releaseparameter_info",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_parameter_info_Releaseparameter_info",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [],					"fields": [						{							"name": "name",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_name",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_parameter_info_get_name",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "name",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "type",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_type",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_parameter_info_get_type",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "type",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "is_default_value",							"type": "bool",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_is_default_value",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_parameter_info_get_is_default_value",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "is_default_value",										"type": "bool",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "is_optional",							"type": "bool",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_is_optional",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_parameter_info_get_is_optional",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "is_optional",										"type": "bool",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "kind",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_kind",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_parameter_info_get_kind",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "kind",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						}					]				},				{					"name": "function_info",					"comment": "",					"tags": {},					"constructors": [],					"release": {						"name": "Releasefunction_info",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_function_info_Releasefunction_info",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [],					"fields": [						{							"name": "name",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_name",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_function_info_get_name",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "name",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "comment",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_comment",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_function_info_get_comment",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "comment",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "parameters",							"type": "handle_array",							"type_alias": "parameter_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_parameters",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_function_info_get_parameters",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "parameters",										"type": "handle_array",										"type_alias": "parameter_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						},						{							"name": "return_values",							"type": "string8_array",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_return_values",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_function_info_get_return_values",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "return_values",										"type": "string8_array",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						}					]				},				{					"name": "class_info",					"comment": "",					"tags": {},					"constructors": [],					"release": {						"name": "Releaseclass_info",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_class_info_Releaseclass_info",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [],					"fields": [						{							"name": "name",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_name",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_class_info_get_name",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "name",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "comment",							"type": "string8",							"type_alias": "",							"comment": "",							"tags": {},							"dimensions": 0,							"getter": {								"name": "get_comment",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_class_info_get_comment",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "comment",										"type": "string8",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"instance_required": true							},							"setter": null						},						{							"name": "fields",							"type": "handle_array",							"type_alias": "variable_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_fields",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_class_info_get_fields",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "fields",										"type": "handle_array",										"type_alias": "variable_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						},						{							"name": "methods",							"type": "handle_array",							"type_alias": "function_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_methods",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_class_info_get_methods",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "methods",										"type": "handle_array",										"type_alias": "function_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						}					]				},				{					"name": "py_info",					"comment": "",					"tags": {},					"constructors": [],					"release": {						"name": "Releasepy_info",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_py_info_Releasepy_info",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [],					"fields": [						{							"name": "globals",							"type": "handle_array",							"type_alias": "variable_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_globals",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_py_info_get_globals",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "globals",										"type": "handle_array",										"type_alias": "variable_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						},						{							"name": "functions",							"type": "handle_array",							"type_alias": "function_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_functions",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_py_info_get_functions",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "functions",										"type": "handle_array",										"type_alias": "function_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						},						{							"name": "classes",							"type": "handle_array",							"type_alias": "class_info",							"comment": "",							"tags": {},							"dimensions": 1,							"getter": {								"name": "get_classes",								"comment": "",								"tags": {},								"function_path": {									"entrypoint_function": "EntryPoint_py_info_get_classes",									"metaffi_guest_lib": "py_extractor_MetaFFIGuest",									"module": "py_extractor"								},								"parameters": [									{										"name": "this_instance",										"type": "handle",										"type_alias": "",										"comment": "",										"tags": {},										"dimensions": 0									}								],								"return_values": [									{										"name": "classes",										"type": "handle_array",										"type_alias": "class_info",										"comment": "",										"tags": {},										"dimensions": 1									}								],								"instance_required": true							},							"setter": null						}					]				},				{					"name": "py_extractor",					"comment": "",					"tags": {},					"constructors": [						{							"name": "py_extractor",							"comment": "",							"tags": {},							"function_path": {								"entrypoint_function": "EntryPoint_py_extractor_py_extractor",								"metaffi_guest_lib": "py_extractor_MetaFFIGuest",								"module": "py_extractor"							},							"parameters": [								{									"name": "filename",									"type": "string8",									"type_alias": "",									"comment": "",									"tags": {},									"dimensions": 0								}							],							"return_values": [								{									"name": "new_instance",									"type": "handle",									"type_alias": "",									"comment": "",									"tags": {},									"dimensions": 0								}							]						}					],					"release": {						"name": "Releasepy_extractor",						"comment": "Releases object",						"tags": {},						"function_path": {							"entrypoint_function": "EntryPoint_py_extractor_Releasepy_extractor",							"metaffi_guest_lib": "py_extractor_MetaFFIGuest",							"module": "py_extractor"						},						"parameters": [							{								"name": "this_instance",								"type": "handle",								"type_alias": "",								"comment": "",								"tags": {},								"dimensions": 0							}						],						"return_values": []					},					"methods": [						{							"name": "extract",							"comment": "",							"tags": {},							"function_path": {								"entrypoint_class": "py_extractor",								"entrypoint_function": "EntryPoint_py_extractor_extract",								"metaffi_guest_lib": "py_extractor_MetaFFIGuest",								"module": "py_extractor"							},							"parameters": [								{									"name": "this_instance",									"type": "handle",									"type_alias": "",									"comment": "",									"tags": {},									"dimensions": 0								}							],							"return_values": [								{									"name": "info",									"type": "handle",									"type_alias": "py_info",									"comment": "",									"tags": {},									"dimensions": 0								}							],							"instance_required": true						}					],					"fields": []				}			],			"globals": [],			"external_resources": []		}	]}`

// --------------------------------------------------------------------
func TestPyExtractorHost(t *testing.T) {

	def, err := IDL.NewIDLDefinitionFromJSON(py_extractor_json)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir("temp_host", 0700)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		err = os.RemoveAll("temp_host")
		if err != nil {
			t.Fatal(err)
			return
		}
	}()

	cmp := NewHostCompiler()
	err = cmp.Compile(def, "temp_host", "", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
}

//--------------------------------------------------------------------
