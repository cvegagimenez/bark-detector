package otel

import (
	"context"
	"errors"
	"log"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	stdoutmetric "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"

	//"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	provider "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var (
	meter       metric.Meter
	barkValue   float64
	ID          string
)

func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil

		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	res, err := newResource()
	if err != nil {
		handleErr(err)
		return
	}

	meterProvider, err := newMeterProvider(ctx, res)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)

	otel.SetMeterProvider(meterProvider)

	meter = otel.GetMeterProvider().Meter("bark-detector")

	err = setupMetrics()
	if err != nil {
		handleErr(err)
		return
	}

	return
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("bark-detector"),
		))
}

func newMeterProvider(ctx context.Context, res *resource.Resource) (*provider.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	//metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`http\.server\.(request|response)\.size`)
	var dropMetricsView provider.View = func(i provider.Instrument) (provider.Stream, bool) {
		s := provider.Stream{
			Name:        i.Name,
			Description: i.Description,
			Unit:        i.Unit,
			Aggregation: provider.AggregationDrop{},
		}

		if re.MatchString(i.Name) {
			return s, true
		}

		return s, false
	}

	meterProvider := provider.NewMeterProvider(
		provider.WithResource(res),
		provider.WithReader(provider.NewPeriodicReader(metricExporter)),
		provider.WithView(
			dropMetricsView,
		),
	)

	return meterProvider, nil
}

func setupMetrics() error {
    barkGauge, err := meter.Float64ObservableGauge(
        "bark_power",
        metric.WithDescription("Maximum bark power detected"),
        metric.WithUnit("RMS"),
    )
    if err != nil {
        return err
    }

	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			observer.ObserveFloat64(
				barkGauge, 
				barkValue, 
				metric.WithAttributes(
					attribute.String("sensorID", ID),
			))
			return nil
		},
		barkGauge,
	)
	if err != nil {
		return err
	}
	return nil
}

func RecordBarkPower(power float64, sensorID string) {
    barkValue = power
	ID = sensorID
	log.Printf("Recorded bark power: %f", power)
}