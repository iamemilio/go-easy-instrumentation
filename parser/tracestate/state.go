package tracestate

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/newrelic/go-easy-instrumentation/internal/codegen"
	"github.com/newrelic/go-easy-instrumentation/parser/tracestate/traceobject"
)

// State stores the current state of the tracing process.
type State struct {
	main          bool
	definedTxn    bool
	agentVariable string
	txnVariable   string
	object        traceobject.TraceObject
}

// Main creates a new State object for tracing a main function.
// We know the agent must be initialized in the main function.
//
// The agentVariable is the name of the agent variable in the main function.
func Main(agentVariable string) *State {
	return &State{
		main:          true,
		agentVariable: agentVariable,
		txnVariable:   codegen.DefaultTransactionVariable,
		object:        traceobject.NewTransaction(),
	}
}

func DownstreamFunction(transactionVariable string) *State {
	return &State{
		txnVariable: transactionVariable,
		object:      traceobject.NewTransaction(),
	}
}

// CreateTransactionIfNeeded creates a transaction in the line before the current cursor position if:
//  1. The agent variable is in scope
//  2. The cursor is in a function body
//  3. We are in the main method
//
// Setting endAfterStatement to true will wrap the current cursor position with a transaction by inserting a transaction end statement after the cursor.
// This is useful for tracing a function call statement.
func (tc *State) CreateTransactionIfNeeded(c *dstutil.Cursor, functionName string, endAfterStatement bool) {
	if tc.main && tc.agentVariable != "" && c.Index() > 0 {
		c.InsertBefore(codegen.StartTransaction(tc.agentVariable, tc.txnVariable, functionName, tc.definedTxn))
		tc.definedTxn = true
		if endAfterStatement {
			c.InsertAfter(codegen.EndTransaction(tc.txnVariable))
		}
	}
}

// TransactionVariable returns the name of the transaction variable.
func (tc *State) TransactionVariable() string {
	return tc.txnVariable
}

// AgentVariable returns the name of the agent variable.
// This may return an empty string if no agent variable is in scope.
func (tc *State) AgentVariable() string {
	return tc.agentVariable
}

// DownstreamFunction returns a new State object for tracing a downstream function.
func (tc *State) DownstreamFunction() *State {
	return &State{
		txnVariable: tc.txnVariable,
		object:      tc.object,
	}
}

func (tc *State) AddTracingToCall(pkg *decorator.Package, call *dst.CallExpr, async bool) {
	tc.object.AddToCall(pkg, call, tc.txnVariable, async)
}

func (tc *State) AddTracingToFunctionDecl(pkg *decorator.Package, decl *dst.FuncDecl) {
	tc.object.AddToFuncDecl(pkg, decl)
}

func (tc *State) AddTracingToFunctionLiteral(pkg *decorator.Package, lit *dst.FuncLit) {
	tc.object.AddToFuncLit(pkg, lit)
}

func (tc *State) AssignTransactionVariable(variableName string) dst.Stmt {
	if tc.txnVariable != "" {
		return tc.object.AssignTransactionVariable(variableName)
	}
	return nil
}
