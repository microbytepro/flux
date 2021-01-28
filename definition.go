package flux

import (
	"github.com/spf13/cast"
)

type (
	EventType int
)

// 路由元数据事件类型
const (
	EventTypeAdded = iota
	EventTypeUpdated
	EventTypeRemoved
)

const (
	// 从动态Path参数中获取
	ScopePath = "PATH"
	// 查询所有Path参数
	ScopePathMap = "PATH_MAP"
	// 从Query参数中获取
	ScopeQuery      = "QUERY"
	ScopeQueryMulti = "QUERY_MUL"
	// 获取全部Query参数
	ScopeQueryMap = "QUERY_MAP"
	// 只从Form表单参数参数列表中读取
	ScopeForm      = "FORM"
	ScopeFormMulti = "FORM_MUL"
	// 获取Form全部参数
	ScopeFormMap = "FORM_MAP"
	// 只从Query和Form表单参数参数列表中读取
	ScopeParam = "PARAM"
	// 只从Header参数中读取
	ScopeHeader = "HEADER"
	// 获取Header全部参数
	ScopeHeaderMap = "HEADER_MAP"
	// 获取Http Attributes的单个参数
	ScopeAttr = "ATTR"
	// 获取Http Attributes的Map结果
	ScopeAttrs = "ATTRS"
	// 获取Body数据
	ScopeBody = "BODY"
	// 获取Request元数据
	ScopeRequest = "REQUEST"
	// 自动查找数据源
	ScopeValue = "VALUE"
	// 自动查找数据源
	ScopeAuto = "AUTO"
)

const (
	// 原始参数类型：int,long...
	ArgumentTypePrimitive = "PRIMITIVE"
	// 复杂参数类型：POJO
	ArgumentTypeComplex = "COMPLEX"
)

// Support protocols
const (
	ProtoDubbo = "DUBBO"
	ProtoGRPC  = "GRPC"
	ProtoHttp  = "HTTP"
	ProtoEcho  = "ECHO"
)

// ServiceAttributes
const (
	ServiceAttrTagNotDefined uint8 = iota
	ServiceAttrTagRpcProto
	ServiceAttrTagRpcGroup
	ServiceAttrTagRpcVersion
	ServiceAttrTagRpcTimeout
	ServiceAttrTagRpcRetries
)

var ServiceAttrTagNames = map[uint8]string{
	ServiceAttrTagNotDefined: "NotDefined",
	ServiceAttrTagRpcProto:   "RpcProto",
	ServiceAttrTagRpcGroup:   "RpcGroup",
	ServiceAttrTagRpcVersion: "RpcVersion",
	ServiceAttrTagRpcTimeout: "RpcTimeout",
	ServiceAttrTagRpcRetries: "RpcRetries",
}

// EndpointAttributes
const (
	EndpointAttrTagNotDefined uint8 = iota
	EndpointAttrTagAuthorize
)

var EndpointAttrTagNames = map[uint8]string{
	EndpointAttrTagNotDefined: "NotDefined",
	EndpointAttrTagAuthorize:  "Authorize",
}

type (
	// ArgumentLookupFunc 参数值查找函数
	ArgumentLookupFunc func(scope, key string, ctx Context) (MTValue, error)
)

// Argument 定义Endpoint的参数结构元数据
type Argument struct {
	Name      string     `json:"name" yaml:"name"`           // 参数名称
	Type      string     `json:"type" yaml:"type"`           // 参数结构类型
	Class     string     `json:"class" yaml:"class"`         // 参数类型
	Generic   []string   `json:"generic" yaml:"generic"`     // 泛型类型
	HttpName  string     `json:"httpName" yaml:"httpName"`   // 映射Http的参数Key
	HttpScope string     `json:"httpScope" yaml:"httpScope"` // 映射Http参数值域
	Fields    []Argument `json:"fields" yaml:"fields"`       // 子结构字段
	// helper
	ValueLoader   func() MTValue     `json:"-"`
	LookupFunc    ArgumentLookupFunc `json:"-"`
	ValueResolver MTValueResolver    `json:"-"`
}

// Attribute 定义服务的属性信息
type Attribute struct {
	Tag   uint8       `json:"tag" yaml:"tag"`
	Name  string      `json:"name" yaml:"name"`
	Value interface{} `json:"value" yaml:"value"`
}

func (a Attribute) ValueString() string {
	return cast.ToString(a.Value)
}

func (a Attribute) ValueInt() int {
	return cast.ToInt(a.Value)
}

func (a Attribute) ValueBool() bool {
	return cast.ToBool(a.Value)
}

// EmbeddedAttributes
type EmbeddedAttributes struct {
	Attributes []Attribute `json:"attributes" yaml:"attributes"`
}

func (c EmbeddedAttributes) AttrByTag(tag uint8) Attribute {
	return c.AttrLookup(func(attr Attribute) bool {
		return tag == attr.Tag
	})
}

func (c EmbeddedAttributes) AttrByName(name string) Attribute {
	return c.AttrLookup(func(attr Attribute) bool {
		return name == attr.Name
	})
}

func (c EmbeddedAttributes) AttrLookup(tf func(Attribute) bool) Attribute {
	for _, attr := range c.Attributes {
		if tf(attr) {
			return attr
		}
	}
	return Attribute{}
}

// EmbeddedExtensions
type EmbeddedExtensions struct {
	Extensions map[string]interface{} `json:"extensions" yaml:"extensions"` // 扩展信息
}

func (e EmbeddedExtensions) Ext(name string) (interface{}, bool) {
	v, ok := e.Extensions[name]
	return v, ok
}

func (e EmbeddedExtensions) ExtString(name string) string {
	v, ok := e.Extensions[name]
	if ok {
		return cast.ToString(v)
	} else {
		return ""
	}
}

func (e EmbeddedExtensions) ExtBool(name string) bool {
	v, ok := e.Extensions[name]
	if ok {
		return cast.ToBool(v)
	} else {
		return false
	}
}

func (e EmbeddedExtensions) ExtInt(name string) int {
	v, ok := e.Extensions[name]
	if ok {
		return cast.ToInt(v)
	} else {
		return 0
	}
}

// BackendService 定义连接上游目标服务的信息
type BackendService struct {
	AliasId    string     `json:"aliasId" yaml:"aliasId"`       // Service别名
	ServiceId  string     `json:"serviceId" yaml:"serviceId"`   // Service的标识ID
	Scheme     string     `json:"scheme" yaml:"scheme"`         // Service侧URL的Scheme
	RemoteHost string     `json:"remoteHost" yaml:"remoteHost"` // Service侧的Host
	Interface  string     `json:"interface" yaml:"interface"`   // Service侧的URL
	Method     string     `json:"method" yaml:"method"`         // Service侧的方法
	Arguments  []Argument `json:"arguments" yaml:"arguments"`   // Service侧的参数结构
	// Extends
	EmbeddedAttributes `yaml:",inline"`
	EmbeddedExtensions `yaml:",inline"`
	RpcProto           string `json:"rpcProto" yaml:"rpcProto"`     // Deprecated Service侧的协议
	RpcGroup           string `json:"rpcGroup" yaml:"rpcGroup"`     // Deprecated Service侧的接口分组
	RpcVersion         string `json:"rpcVersion" yaml:"rpcVersion"` // Deprecated Service侧的接口版本
	RpcTimeout         string `json:"rpcTimeout" yaml:"rpcTimeout"` // Deprecated Service侧的调用超时
	RpcRetries         string `json:"rpcRetries" yaml:"rpcRetries"` // Deprecated Service侧的调用重试
}

func (b BackendService) AttrRpcProto() string {
	return b.AttrByTag(ServiceAttrTagRpcProto).ValueString()
}

func (b BackendService) AttrRpcTimeout() string {
	return b.AttrByTag(ServiceAttrTagRpcTimeout).ValueString()
}

func (b BackendService) AttrRpcGroup() string {
	return b.AttrByTag(ServiceAttrTagRpcGroup).ValueString()
}

func (b BackendService) AttrRpcVersion() string {
	return b.AttrByTag(ServiceAttrTagRpcVersion).ValueString()
}

func (b BackendService) AttrRpcRetries() string {
	return b.AttrByTag(ServiceAttrTagRpcRetries).ValueString()
}

// IsValid 判断服务配置是否有效；Interface+Method不能为空；
func (b BackendService) IsValid() bool {
	return "" != b.Interface && "" != b.Method
}

// HasArgs 判定是否有参数
func (b BackendService) HasArgs() bool {
	return len(b.Arguments) > 0
}

// ServiceID 构建标识当前服务的ID
func (b BackendService) ServiceID() string {
	return b.Interface + ":" + b.Method
}

// Endpoint 定义前端Http请求与后端RPC服务的端点元数据
type Endpoint struct {
	Application        string         `json:"application" yaml:"application"` // 所属应用名
	Version            string         `json:"version" yaml:"version"`         // 端点版本号
	HttpPattern        string         `json:"httpPattern" yaml:"httpPattern"` // 映射Http侧的UriPattern
	HttpMethod         string         `json:"httpMethod" yaml:"httpMethod"`   // 映射Http侧的Method
	Service            BackendService `json:"service" yaml:"service"`         // 上游服务
	Permission         BackendService `json:"permission" yaml:"permission"`   // Deprecated 权限验证定义
	Permissions        []string       `json:"permissions" yaml:"permissions"` // 多组权限验证服务ID列表
	EmbeddedAttributes `yaml:",inline"`
	EmbeddedExtensions `yaml:",inline"`
}

func (e Endpoint) PermissionServiceIds() []string {
	ids := make([]string, 0, 1+len(e.Permissions))
	if e.Permission.IsValid() {
		ids = append(ids, e.Permission.ServiceId)
	}
	ids = append(ids, e.Permissions...)
	return ids
}

func (e Endpoint) IsValid() bool {
	return "" != e.HttpMethod && "" != e.HttpPattern && e.Service.IsValid()
}

func (e Endpoint) AttrAuthorize() bool {
	return e.AttrByTag(EndpointAttrTagAuthorize).ValueBool()
}

// HttpEndpointEvent  定义从注册中心接收到的Endpoint数据变更
type HttpEndpointEvent struct {
	EventType EventType
	Endpoint  Endpoint
}

// BackendServiceEvent  定义从注册中心接收到的Service定义数据变更
type BackendServiceEvent struct {
	EventType EventType
	Service   BackendService
}

func EnsureServiceAttribute(tag uint8, name string) (uint8, string) {
	return EnsureAttribute(tag, name, ServiceAttrTagNames)
}

func EnsureEndpointAttribute(tag uint8, name string) (uint8, string) {
	return EnsureAttribute(tag, name, EndpointAttrTagNames)
}

func EnsureAttribute(tag uint8, name string, mapping map[uint8]string) (uint8, string) {
	if tag > 0 && name == "" {
		name = mapping[tag]
	}
	if name != "" && tag <= 0 {
		for t, n := range mapping {
			if name == n {
				tag = t
			}
		}
	}
	return tag, name
}
