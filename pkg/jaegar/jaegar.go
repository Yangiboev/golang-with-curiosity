package jaegar

import (
	"io"

	"github.com/Yangiboev/golang-with-curiosity/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaegar(cfg config.Config) (opentracing.Tracer, io.Closer, error) {
	jaegarCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaegar.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaegar.LogSpans,
			LocalAgentHostPort: cfg.Jaegar.Host,
		},
	}
	return jaegarCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
