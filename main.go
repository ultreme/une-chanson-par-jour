package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/rs/rest-layer-mem"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
)

var (
	user = schema.Schema{
		Description: `The user object`,
		Fields: schema.Fields{
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,

			"name": {
				Required:   true,
				Filterable: true,
				Validator: &schema.String{
					MaxLen: 100,
				},
			},
		},
	}

	season = schema.Schema{
		Description: `The season object`,
		Fields: schema.Fields{
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,

			"name": {
				Required:   true,
				Filterable: true,
				Validator: &schema.String{
					MaxLen: 100,
				},
			},
		},
	}

	song = schema.Schema{
		Description: `The song object`,
		Fields: schema.Fields{
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,

			"user": {
				Required:   true,
				Filterable: true,
				Validator: &schema.Reference{
					Path: "users",
				},
			},

			"season": {
				Required:   true,
				Filterable: true,
				Validator: &schema.Reference{
					Path: "seasons",
				},
			},

			"name": {
				Required: true,
				Validator: &schema.String{
					MaxLen: 100,
				},
			},

			"url": {
				Required: true,
				Validator: &schema.String{
					MaxLen: 4096,
				},
			},
		},
	}

	vote = schema.Schema{
		Description: `the vote object`,
		Fields: schema.Fields{
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,

			"auth": {
				Required:   true,
				Filterable: true,
				Validator: &schema.Reference{
					Path: "users",
				},
			},
			"song": {
				Required:   true,
				Filterable: true,
				Validator: &schema.Reference{
					Path: "songs",
				},
			},
		},
	}
)

func main() {
	// Create a REST API resource index
	index := resource.NewIndex()
	index.Bind("users", user, mem.NewHandler(), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})
	index.Bind("seasons", season, mem.NewHandler(), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})
	index.Bind("songs", song, mem.NewHandler(), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	// Create API HTTP handler for the resource graph
	api, err := rest.NewHandler(index)
	if err != nil {
		log.Fatalf("Invalid API configuration: %s", err)
	}

	// Init a xhandler chain (see https://github.com/rs/xhandler)
	c := xhandler.Chain{}
	c.UseC(xhandler.CloseHandler)
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))
	c.UseC(xlog.NewHandler(xlog.Config{}))
	c.UseC(xaccess.NewHandler())
	c.UseC(cors.New(cors.Options{OptionsPassthrough: true}).HandlerC)

	// Bind the API under /api/ path
	http.Handle("/api/", http.StripPrefix("/api/", c.Handler(api)))

	// Serve it
	log.Print("Serving API on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
