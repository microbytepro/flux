package flux

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrHttpRequestNotSupported  = errors.New("webserver: http.request not supported")
	ErrHttpResponseNotSupported = errors.New("webserver: http.responsewriter not supported")
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// MIME types
const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationForm            = "application/x-www-form-urlencoded"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"

	// Ext
	HeaderXRequestId     = "X-Request-Id"
	HeaderValueSeparator = ","
)

// Web interfaces
type (
	// WebMiddleware 定义处理Web请求的中间件函数
	WebMiddleware func(WebRouteHandler) WebRouteHandler

	// WebRouteHandler 定义处理Web请求的处理函数
	WebRouteHandler func(WebContext) error

	// WebRouteHandler 定义Web服务处理异常错误的处理函数
	WebErrorHandler func(err error, ctx WebContext)

	// WebSkipper 用于部分Middleware逻辑，实现忽略部分请求的功能；
	WebSkipper func(ctx WebContext) bool
)

// WebContext 定义封装Web框架的RequestContext的接口；用于RequestHandler，Middleware实现Web请求处理；
type WebContext interface {
	// Method 返回请求的HttpMethod
	Method() string

	// Host 返回请求的Host
	Host() string

	// UserAgent 返回请求的UserAgent
	UserAgent() string

	// Request 返回Http标准Request对象。
	// 如果Web框架不支持标准Request（如fasthttp），返回 ErrHttpRequestNotSupported
	Request() (*http.Request, error)

	// RequestURI 返回请求的URI
	RequestURI() string

	// RequestURL 返回请求对象的URL
	// 注意：部分Web框架返回只读url.URL
	RequestURL() (url *url.URL, readonly bool)

	// RequestBodyReader 返回可重复读取的Reader接口；
	RequestBodyReader() (io.ReadCloser, error)

	// RequestRewrite 修改请求方法和路径；
	RequestRewrite(method string, path string)

	// RequestHeader 返回请求对象的Header
	// 注意：部分Web框架返回只读http.Header
	RequestHeader() (header http.Header, readonly bool)

	// GetRequestHeader 读取请求的Header
	GetRequestHeader(name string) string

	// SetRequestHeader 设置请求的Header的键值对
	SetRequestHeader(name, value string)

	// AddRequestHeader 添加请求指定Name的Header的键值
	AddRequestHeader(name, value string)

	// QueryValues 返回Query查询参数键值对；只读；
	QueryValues() url.Values

	// PathValues 返回动态路径参数的键值对；只读；
	PathValues() url.Values

	// FormValues 返回Form表单参数键值对；只读；
	FormValues() url.Values

	// QueryValues 返回Cookie列表；只读；
	CookieValues() []*http.Cookie

	// QueryValue 查询指定Name的Query参数值
	QueryValue(name string) string

	// PathValue 查询指定Name的动态路径参数值
	PathValue(name string) string

	// FormValue 查询指定Name的表单参数值
	FormValue(name string) string

	// CookieValue 查询指定Name的Cookie对象，并返回是否存在标识
	CookieValue(name string) (cookie *http.Cookie, ok bool)

	// Response 返回Http标准ResponseWriter对象。
	// 如果Web框架不支持标准ResponseWriter（如fasthttp），返回 ErrHttpResponseNotSupported
	Response() (http.ResponseWriter, error)

	// SetResponseWriter 设置响应的ResponseWriter
	SetResponseWriter(w http.ResponseWriter)

	// ResponseHeader 返回响应对象的Header以及是否只读
	// 注意：部分Web框架返回只读http.Header
	ResponseHeader() (header http.Header, readonly bool)

	// ResponseWrite 写入响应状态码和响应数据
	ResponseWrite(statusCode int, contentType string, bytes []byte) error

	// ResponseWriteStream 写入响应状态码和流数据
	ResponseWriteStream(statusCode int, contentType string, reader io.Reader) error

	// ResponseNoContent 返回无数据响应
	ResponseNoContent(statusCode int)

	// ResponseRedirect 返回客户端重定向
	ResponseRedirect(statusCode int, url string)

	// GetResponseHeader 获取已设置的Header键值
	GetResponseHeader(name string) string

	// SetResponseHeader 设置的Header键值
	SetResponseHeader(name, value string)

	// AddResponseHeader 添加指定Name的Header键值
	AddResponseHeader(name, value string)

	// SetValue 设置Context域键值；作用域与请求生命周期相同；
	SetValue(name string, value interface{})

	// GetValue 获取Context域键值；作用域与请求生命周期相同；
	GetValue(name string) interface{}

	// Context 返回具体Web框架实现的RequestContext对象
	Context() interface{}
}

// WebServer 定义Web框架服务器的接口；通过实现此接口来自定义支持不同的Web框架，用于支持不同的Web服务实现。
// 例如默认Web框架为labstack.echo；可以支持git, fasthttp等框架。
type WebServer interface {
	// SetWebErrorHandler 设置Web请求错误处理函数
	SetWebErrorHandler(h WebErrorHandler)

	// SetRouteNotFoundHandler 设置Web路由不存在处理函数
	SetRouteNotFoundHandler(h WebRouteHandler)

	// AddWebInterceptor 添加全局请求拦截器，作用于路由请求前
	AddWebInterceptor(m WebMiddleware)

	// AddWebMiddleware 添加全局中间件函数，作用于路由请求后
	AddWebMiddleware(m WebMiddleware)

	// AddWebRouteHandler 添加请求路由处理函数及其中间件
	AddWebRouteHandler(method, pattern string, h WebRouteHandler, m ...WebMiddleware)

	// AddStdHttpHandler 添加http标准请求路由处理函数及其中间件
	AddStdHttpHandler(method, pattern string, h http.Handler, m ...func(http.Handler) http.Handler)

	// WebServer 返回具体实现的WebServer服务对象，如echo,fasthttp的Server
	WebServer() interface{}

	// WebServer 返回具体实现的WebRouter路由处理对象，如echo,fasthttp的Router
	WebRouter() interface{}

	// Start 启动服务
	Start(addr string) error

	// StartTLS 启动TLS服务
	StartTLS(addr string, certFile, keyFile string) error

	// Shutdown 停止服务
	Shutdown(ctx context.Context) error
}

// WebServerResponseWriter 实现将错误消息和响应数据写入Web服务响应对象
type WebServerResponseWriter interface {
	// WriteError 写入Error错误响应数据到WebServer
	WriteError(webc WebContext, requestId string, header http.Header, error *StateError) error

	// WriteBody 写入Body正常响应数据到WebServer
	WriteBody(webc WebContext, requestId string, header http.Header, status int, body interface{}) error
}

/// Wrapper functions

func WrapHttpHandler(h http.Handler) WebRouteHandler {
	return func(webc WebContext) error {
		// 注意：部分Web框架不支持返回标准Request/Response
		resp, err := webc.Response()
		if nil != err {
			return err
		}
		req, err := webc.Request()
		if nil != err {
			return err
		}
		h.ServeHTTP(resp, req)
		return nil
	}
}