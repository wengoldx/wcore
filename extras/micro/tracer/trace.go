package tracer

import (
	"context"
	"encoding/json"
	bctx "github.com/astaxie/beego/context"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/wengoldx/wing/logger"
	"io"
	"net/http"
	"time"
)

// NewTrace init service tracer
func NewTrace(service string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   "http://jaeger:14268/api/traces",
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	logger.I("init jaeger client config ")
	return tracer, closer, err
}

// TraceMethod Intercept spanid from context use to micro srv
func TraceMethod(ctx context.Context, method string) opentracing.Span {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	var sp opentracing.Span
	spID, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	sp = opentracing.StartSpan(method, opentracing.ChildOf(spID))
	return sp
}

// TraceWrapp create web service trace middleware for beego
func TraceWrapp(bct *bctx.Context) {
	sp := opentracing.GlobalTracer().StartSpan(bct.Request.URL.Path)
	tracer := opentracing.GlobalTracer()
	md := make(map[string]string)
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(bct.Request.Header))
	if err == nil {
		sp = opentracing.GlobalTracer().StartSpan(bct.Request.URL.Path, opentracing.ChildOf(spanCtx))
		tracer = sp.Tracer()
		logger.I("extract spandid from http header")
	}
	defer sp.Finish()

	if err := tracer.Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {

	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	bct.Input.SetData("Tracer-context", ctx)

	bct.Request.ParseForm()
	params, _ := json.Marshal(bct.Request.Form)
	sp.SetTag("Method", bct.Request.Method)
	sp.SetTag("URL", bct.Request.URL.EscapedPath())
	sp.SetTag("Params", string(params)+string(bct.Input.RequestBody))
}

// TraceHTTP tracer http request,for rewrite micro api plugins
func TraceHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			logger.I("extract spanid err：", err)
		}
		sp := opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
		defer sp.Finish()

		if err := opentracing.GlobalTracer().Inject(sp.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header)); err != nil {

		}

		sct := &StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
		h.ServeHTTP(sct.WrappedResponseWriter(), r)

		sp.SetTag("Method", r.Method)
		sp.SetTag("URL", r.URL.EscapedPath())
		sp.SetTag("Code", uint16(sct.Status))
	})
}

// HTTPToMicro http request context tranfer to micro context
// return spanid and context.context
func HTTPToMicro(req *http.Request, method string) (context.Context, opentracing.Span) {
	md := make(map[string]string)
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {

	}
	sp := opentracing.GlobalTracer().StartSpan(method, opentracing.ChildOf(spanCtx))
	if err := opentracing.GlobalTracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {

	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	return ctx, sp
}

// ContextWithSpan tranfer beego context to context
func ContextWithSpan(c *bctx.Context) context.Context {
	v := c.Input.GetData("Tracer-context")
	if v == nil {
		return context.TODO()
	}
	ctx := v.(context.Context)
	return ctx
}
