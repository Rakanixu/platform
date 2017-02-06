package context

/* This package defines all different context key - value pairs used on platform */

const (
	CTX_SCOPE_INTERNAL = "internal"
)

type UserIdCtxKey struct{}

type UserIdCtxValue string

type ScopeCtxKey struct{}

type ScopeCtxValue string
