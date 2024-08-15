package main

import (
	"reflect"
	"testing"

	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
)

func TestInstrumentationManager_AddImport(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect bool
	}{
		{
			name: "AddImport",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{}}},
				currentPackage: "foo",
			},
			args:   args{path: "bar"},
			expect: true,
		},
		{
			name: "AddImport_nil_check",
			fields: fields{
				packages: map[string]*PackageState{"foo": {importsAdded: map[string]bool{}}},
			},
			args:   args{path: "bar"},
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}

			defer panicRecovery(t)
			m.AddImport(tt.args.path)

			if m.packages["foo"].importsAdded["bar"] != true && tt.expect {
				t.Errorf("AddImport failed to add import bar to package foo, got: %+v", m.packages["foo"].importsAdded)
			}
			if tt.expect == false && len(m.packages["foo"].importsAdded) != 0 {
				t.Errorf("AddImport added import bar to package foo, got: %+v", m.packages["foo"].importsAdded)
			}

		})
	}
}

func TestInstrumentationManager_GetImports(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "GetImports_one_import",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{"bar": true}}},
				currentPackage: "foo",
			},
			args: args{fileName: "foo"},
			want: []string{"bar"},
		},
		{
			name: "GetImports_empty",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{}}},
				currentPackage: "foo",
			},
			args: args{fileName: "foo"},
			want: []string{},
		},
		{
			name: "GetImports_nil_check",
			fields: fields{
				packages: map[string]*PackageState{"foo": {importsAdded: map[string]bool{}}},
			},
			args: args{fileName: "foo"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			if got := m.GetImports(tt.args.fileName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstrumentationManager.GetImports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstrumentationManager_CreateFunctionDeclaration(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		decl *dst.FuncDecl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect bool
	}{
		{
			name: "CreateFunctionDeclaration",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{}, tracedFuncs: map[string]*tracedFunction{}}},
				currentPackage: "foo",
			},
			args:   args{decl: &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}}},
			expect: true,
		},
		{
			name: "CreateFunctionDeclaration_nil_check",
			fields: fields{
				packages: map[string]*PackageState{"foo": {importsAdded: map[string]bool{}, tracedFuncs: map[string]*tracedFunction{}}},
			},
			args:   args{decl: &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}}},
			expect: false,
		},
		{
			name: "CreateFunctionDeclaration_already_exists",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{}, tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args:   args{decl: &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}}},
			expect: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			defer panicRecovery(t)
			m.CreateFunctionDeclaration(tt.args.decl)

			if tt.expect {
				if m.packages["foo"].tracedFuncs["bar"] == nil {
					t.Errorf("CreateFunctionDeclaration failed to add new function bar to package foo, got: %+v", m.packages["foo"].tracedFuncs)
				}
				if len(m.packages["foo"].tracedFuncs) != 1 {
					t.Errorf("CreateFunctionDeclaration must not create a new entry if one already exists with that function name: %+v", m.packages["foo"].tracedFuncs)
				}
			}
			if !tt.expect {
				_, ok := m.packages["foo"].tracedFuncs["bar"]
				if ok {
					t.Errorf("CreateFunctionDeclaration added function bar to package foo when it should not have: %+v", m.packages["foo"].tracedFuncs)
				}
			}
		})
	}
}

func TestInstrumentationManager_UpdateFunctionDeclaration(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		decl *dst.FuncDecl
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		updates bool
	}{
		{
			name: "UpdateFunctionDeclaration",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {importsAdded: map[string]bool{}, tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args:    args{decl: &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}}},
			updates: true,
		},
		{
			name: "UpdateFunctionDeclaration_nil_check",
			fields: fields{
				packages: map[string]*PackageState{"foo": {importsAdded: map[string]bool{}, tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
			},
			args:    args{decl: &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}}},
			updates: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}

			defer panicRecovery(t)
			m.UpdateFunctionDeclaration(tt.args.decl)

			if tt.updates && reflect.DeepEqual(m.packages["foo"].tracedFuncs["bar"].body, tt.args.decl) == false {
				t.Errorf("UpdateFunctionDeclaration failed to update function bar to package foo, got: %+v", m.packages["foo"].tracedFuncs)
			}

			if !tt.updates && reflect.DeepEqual(m.packages["foo"].tracedFuncs["bar"].body, tt.args.decl) == true {
				t.Errorf("UpdateFunctionDeclaration updated function bar to package foo when it should not have: %+v", m.packages["foo"].tracedFuncs)
			}
		})
	}
}

// What if there are two instrumentable function invocations in a statement?
func TestInstrumentationManager_GetPackageFunctionInvocation(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		node dst.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *invocationInfo
	}{
		{
			name: "basic_passing_case",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {}},
				currentPackage: "foo",
			},
			args: args{node: &dst.CallExpr{Fun: &dst.Ident{Name: "bar", Path: "foo"}}},
			want: &invocationInfo{packageName: "foo", functionName: "bar", call: &dst.CallExpr{Fun: &dst.Ident{Name: "bar", Path: "foo"}}},
		},
		{
			name: "empty_path_passes",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {}},
				currentPackage: "foo",
			},
			args: args{node: &dst.CallExpr{Fun: &dst.Ident{Name: "bar"}}},
			want: &invocationInfo{packageName: "foo", functionName: "bar", call: &dst.CallExpr{Fun: &dst.Ident{Name: "bar"}}},
		},
		{
			name: "finds_call_in_complex_node",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {}},
				currentPackage: "foo",
			},
			args: args{node: &dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Name: "Sprintf", Path: "fmt"}, Args: []dst.Expr{&dst.CallExpr{Fun: &dst.Ident{Name: "bar"}}}}}},
			want: &invocationInfo{packageName: "foo", functionName: "bar", call: &dst.CallExpr{Fun: &dst.Ident{Name: "bar"}}},
		},
		{
			name: "ignore_functions_not_in_package",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {}},
				currentPackage: "foo",
			},
			args: args{node: &dst.CallExpr{Fun: &dst.Ident{Name: "bar", Path: "fmt"}}},
			want: nil,
		},
		{
			name: "ignore_block_statements",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {}},
				currentPackage: "foo",
			},
			args: args{node: &dst.BlockStmt{List: []dst.Stmt{&dst.ExprStmt{X: &dst.CallExpr{Fun: &dst.Ident{Name: "bar"}}}}}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			defer panicRecovery(t)
			got := m.GetPackageFunctionInvocation(tt.args.node)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInstrumentationManager_AddTxnArgumentToFunctionDecl(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		decl       *dst.FuncDecl
		txnVarName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *dst.FuncDecl
		wantRequireTxn bool
	}{
		{
			name: "simple_passing_case",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{
				decl:       &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}, Type: &dst.FuncType{Params: &dst.FieldList{}}},
				txnVarName: "txn",
			},
			want: &dst.FuncDecl{
				Name: &dst.Ident{Name: "bar"},
				Type: &dst.FuncType{
					Params: &dst.FieldList{
						List: []*dst.Field{{
							Names: []*dst.Ident{dst.NewIdent("txn")},
							Type: &dst.StarExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent("newrelic"),
									Sel: dst.NewIdent("Transaction"),
								},
							},
						}},
					},
				},
			},
			wantRequireTxn: true,
		},
		{
			name: "simple_case_nil_params",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{
				decl:       &dst.FuncDecl{Name: &dst.Ident{Name: "bar"}, Type: &dst.FuncType{Params: nil}},
				txnVarName: "txn",
			},
			want: &dst.FuncDecl{
				Name: &dst.Ident{Name: "bar"},
				Type: &dst.FuncType{
					Params: &dst.FieldList{
						List: []*dst.Field{{
							Names: []*dst.Ident{dst.NewIdent("txn")},
							Type: &dst.StarExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent("newrelic"),
									Sel: dst.NewIdent("Transaction"),
								},
							},
						}},
					},
				},
			},
			wantRequireTxn: true,
		},
		{
			name: "nil_function_declaration",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{
				decl:       nil,
				txnVarName: "txn",
			},
			want:           nil,
			wantRequireTxn: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			defer panicRecovery(t)
			m.AddTxnArgumentToFunctionDecl(tt.args.decl, tt.args.txnVarName)
			assert.Equal(t, tt.want, tt.args.decl)
			assert.Equal(t, tt.wantRequireTxn, m.packages[m.currentPackage].tracedFuncs["bar"].requiresTxn)
		})
	}
}

func TestInstrumentationManager_ShouldInstrumentFunction(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		inv *invocationInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "function_should_be_instrumented",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: true,
		},
		{
			name: "nil_invocation",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{inv: nil},
			want: false,
		},
		{
			name: "already_instrumented",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {traced: true}}}},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: false,
		},
		{
			name: "package_not_found",
			fields: fields{
				packages:       map[string]*PackageState{},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			defer panicRecovery(t)
			assert.Equal(t, tt.want, m.ShouldInstrumentFunction(tt.args.inv))
		})
	}
}

func TestInstrumentationManager_RequiresTransactionArgument(t *testing.T) {
	type fields struct {
		userAppPath       string
		diffFile          string
		appName           string
		agentVariableName string
		currentPackage    string
		packages          map[string]*PackageState
	}
	type args struct {
		inv *invocationInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "requres_txn",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {requiresTxn: true}}}},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: true,
		},
		{
			name: "nil_invocation",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {}}}},
				currentPackage: "foo",
			},
			args: args{inv: nil},
			want: false,
		},
		{
			name: "does_not_require_txn",
			fields: fields{
				packages:       map[string]*PackageState{"foo": {tracedFuncs: map[string]*tracedFunction{"bar": {requiresTxn: false}}}},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: false,
		},
		{
			name: "package_not_found",
			fields: fields{
				packages:       map[string]*PackageState{},
				currentPackage: "foo",
			},
			args: args{inv: &invocationInfo{packageName: "foo", functionName: "bar"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InstrumentationManager{
				userAppPath:       tt.fields.userAppPath,
				diffFile:          tt.fields.diffFile,
				appName:           tt.fields.appName,
				agentVariableName: tt.fields.agentVariableName,
				currentPackage:    tt.fields.currentPackage,
				packages:          tt.fields.packages,
			}
			defer panicRecovery(t)
			assert.Equal(t, tt.want, m.RequiresTransactionArgument(tt.args.inv))
		})
	}
}
