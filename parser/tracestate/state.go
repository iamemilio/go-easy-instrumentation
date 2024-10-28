package tracestate

import (
	"github.com/dave/dst/dstutil"
	"github.com/newrelic/go-easy-instrumentation/internal/codegen"
)

// State stores the current state of the tracing process.
type State struct {
	definedTxn    bool
	agentVariable string
	txnVariable   string
}

// Main creates a new State object for tracing a main function.
// We know the agent must be initialized in the main function.
//
// The agentVariable is the name of the agent variable in the main function.
func Main(agentVariable string) *State {
	return &State{
		definedTxn:    false,
		agentVariable: agentVariable,
		txnVariable:   codegen.DefaultTransactionVariable,
	}
}

func DownstreamFunction(txnVariableName string) *State {
	return &State{
		txnVariable: txnVariableName,
	}
}

// CreateTransactionIfNeeded creates a transaction in the line before the current cursor position if:
//  1. The agent variable is in scope
//  2. The cursor is in a function body
//
// Setting endAfterStatement to true will wrap the current cursor position with a transaction by inserting a transaction end statement after the cursor.
// This is useful for tracing a function call statement.
func (tc *State) CreateTransactionIfNeeded(c *dstutil.Cursor, functionName string, endAfterStatement bool) {
	if tc.agentVariable != "" && c.Index() > 0 {
		c.InsertBefore(codegen.StartTransaction(tc.agentVariable, tc.txnVariable, functionName, tc.definedTxn))
		tc.definedTxn = true
		if endAfterStatement {
			c.InsertAfter(codegen.EndTransaction(tc.txnVariable))
		}
	}
}

func (tc *State) TransactionVariable() string {
	return tc.txnVariable
}

func (tc *State) AgentVariable() string {
	return tc.agentVariable
}

func (tc *State) DownstreamFunction() *State {
	return &State{
		txnVariable: tc.txnVariable,
	}
}
