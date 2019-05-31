package monitoring

import (
	"fmt"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Setup initialize the opencensus exporter and trace collection
func Setup(addr string, alwaysSample bool) error {
	oce, err := ocagent.NewExporter(
		// TODO: (@odeke-em), enable ocagent-exporter.WithCredentials option.
		ocagent.WithInsecure(),
		ocagent.WithServiceName("c2ae"),
		ocagent.WithAddress(addr),
	)
	if err != nil {
		return fmt.Errorf("failed to create the OpenCensus Agent exporter: %v", err)
	}

	// and now finally register it as a Trace Exporter
	trace.RegisterExporter(oce)
	view.RegisterExporter(oce)

	// setting trace sample rate to 100%
	if alwaysSample {
		trace.ApplyConfig(trace.Config{
			DefaultSampler: trace.AlwaysSample(),
		})
	}

	return nil
}
