package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"cloud.google.com/go/spanner"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/kelseyhightower/envconfig"
	"github.com/sinmetal/gcpmetadata"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type EnvConfig struct {
	SpannerDatabase string `required:"true"`
}

var spannerClient *spanner.Client

func main() {
	ctx := context.Background()

	var env EnvConfig
	if err := envconfig.Process("stonefoot", &env); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("ENV_CONFIG %+v\n", env)

	if gcpmetadata.OnGCP() {
		projectID, err := gcpmetadata.GetProjectID()
		if err != nil {
			panic(err)
		}
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: projectID,
		})
		if err != nil {
			panic(err)
		}
		trace.RegisterExporter(exporter)
	}

	Initialize(ctx, &env)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	ucon.ListenAndServe(fmt.Sprintf(":%s", port))
}

func Initialize(ctx context.Context, config *EnvConfig) {
	SetUpSpanner(ctx, config.SpannerDatabase)
	SetUpUcon()
}

func SetUpSpanner(ctx context.Context, database string) {
	sc, err := spanner.NewClient(ctx, database)
	if err != nil {
		panic(err)
	}
	spannerClient = sc
}

func SetUpUcon() {
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	swPlugin := swagger.NewPlugin(&swagger.Options{
		Object: &swagger.Object{
			Info: &swagger.Info{
				Title:   "StoneFoot",
				Version: "v1",
			},
			Schemes: []string{"http", "https"},
		},
		DefinitionNameModifier: func(refT reflect.Type, defName string) string {
			if strings.HasSuffix(defName, "JSON") {
				return defName[:len(defName)-4]
			}
			return defName
		},
	})
	ucon.Plugin(swPlugin)

	SetUpTweetAPI(swPlugin)

	http.Handle("/api/", &ochttp.Handler{
		Handler:        ucon.DefaultMux,
		Propagation:    &propagation.HTTPFormat{},
		FormatSpanName: formatSpanName,
	})
}

func formatSpanName(r *http.Request) string {
	return fmt.Sprintf("/stonefoot/%s", r.URL.Path)
}
