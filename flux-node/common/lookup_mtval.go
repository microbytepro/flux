package common

import (
	"errors"
	"github.com/bytepowered/flux/flux-node"
	"github.com/bytepowered/flux/flux-pkg"
	"net/url"
	"strings"
)

// LookupExpr 搜索LookupExpr表达式指定域的值。
func LookupMTValueByExpr(expr string, ctx flux.Context) (interface{}, error) {
	if "" == expr || nil == ctx {
		return nil, errors.New("empty lookup expr, or context is nil")
	}
	scope, key, ok := fluxpkg.LookupParseExpr(expr)
	if !ok {
		return "", errors.New("illegal lookup expr: " + expr)
	}
	mtv, err := LookupMTValue(scope, key, ctx)
	if nil != err {
		return "", err
	}
	return mtv.Value, nil
}

// 默认实现查找MTValue
func LookupMTValue(scope, key string, ctx flux.Context) (value flux.MTValue, err error) {
	if "" == scope || "" == key {
		return flux.WrapObjectMTValue(nil), errors.New("lookup empty scope or key, scope: " + scope + ", key: " + key)
	}
	if nil == ctx {
		return flux.WrapObjectMTValue(nil), errors.New("lookup nil context")
	}
	req := ctx.Request()
	switch strings.ToUpper(scope) {
	case flux.ScopePath:
		return flux.WrapStringMTValue(req.PathVar(key)), nil
	case flux.ScopePathMap:
		return flux.WrapStrValuesMapMTValue(req.PathVars()), nil
	case flux.ScopeQuery:
		return flux.WrapStringMTValue(req.QueryVar(key)), nil
	case flux.ScopeQueryMulti:
		return flux.WrapStrListMTValue(req.QueryVars()[key]), nil
	case flux.ScopeQueryMap:
		return flux.WrapStrValuesMapMTValue(req.QueryVars()), nil
	case flux.ScopeForm:
		return flux.WrapStringMTValue(req.FormVar(key)), nil
	case flux.ScopeFormMap:
		return flux.WrapStrValuesMapMTValue(req.FormVars()), nil
	case flux.ScopeFormMulti:
		return flux.WrapStrListMTValue(req.FormVars()[key]), nil
	case flux.ScopeHeader:
		return flux.WrapStringMTValue(req.HeaderVar(key)), nil
	case flux.ScopeHeaderMap:
		return flux.WrapStrValuesMapMTValue(req.HeaderVars()), nil
	case flux.ScopeAttr:
		v, _ := ctx.GetAttribute(key)
		return flux.WrapObjectMTValue(v), nil
	case flux.ScopeAttrs:
		return flux.WrapStrMapMTValue(ctx.Attributes()), nil
	case flux.ScopeBody:
		reader, err := req.BodyReader()
		return flux.MTValue{Value: reader, MediaType: req.HeaderVar(flux.HeaderContentType)}, err
	case flux.ScopeParam:
		v, _ := fluxpkg.LookupByProviders(key, req.QueryVars, req.FormVars)
		return flux.WrapStringMTValue(v), nil
	case flux.ScopeRequest:
		switch strings.ToLower(key) {
		case "method":
			return flux.WrapStringMTValue(ctx.Method()), nil
		case "uri":
			return flux.WrapStringMTValue(ctx.URI()), nil
		default:
			return flux.WrapStringMTValue(""), nil
		}
	case flux.ScopeAuto:
		fallthrough
	default:
		if v, ok := fluxpkg.LookupByProviders(key, req.PathVars, req.QueryVars, req.FormVars, func() url.Values {
			return url.Values(req.HeaderVars())
		}); ok {
			return flux.WrapStringMTValue(v), nil
		}
		if v, ok := ctx.GetAttribute(key); ok {
			return flux.WrapObjectMTValue(v), nil
		}
		return flux.WrapObjectMTValue(nil), nil
	}
}
