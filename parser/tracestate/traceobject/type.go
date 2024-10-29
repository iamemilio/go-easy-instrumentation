package traceobject

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// TraceObject is an object that contains New Relic tracing in the form of a transaction.
// Transactions can be injected into various object types that may require different
// methods of retrieval.
//
// This interface defines a standard set of behaviors that all objects containing a transaction
// must implement for the underlying transaction to be usable for tracing.
type TraceObject interface {
	// AddToCall adds a trace object to a call expression, passing it as an argument
	// to the function being invoked in the call.
	AddToCall(pkg *decorator.Package, call *dst.CallExpr, transactionVariable string, async bool)

	// AddToFuncDecl adds a trace object to a function declaration as a parameter, so that
	// trace objects can be passed in calls to this function.
	//
	// Make sure that the package passed is from the same package that the function is defined in.
	AddToFuncDecl(pkg *decorator.Package, decl *dst.FuncDecl)

	// AddToFuncLit adds a trace object to a function literal definition as a parameter, so that
	// trace objects can be passed in calls to this function literal.
	//
	// Make sure that the package passed is from the same package that the function literal is defined in.
	AddToFuncLit(pkg *decorator.Package, lit *dst.FuncLit)

	// AssignTransactionVariable fetches the transaction from the trace object and assigns it to a variable
	AssignTransactionVariable(variableName string) dst.Stmt
}
