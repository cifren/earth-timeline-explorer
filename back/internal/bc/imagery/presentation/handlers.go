// Package presentation exposes the imagery context over HTTP.
package presentation

import (
	"context"
	"errors"
	"net/http"
	"time"

	"frise-chrono/back/internal/bc/imagery/application"
	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/internal/shared/apigen"
	"frise-chrono/back/internal/shared/problem"
)

const (
	availabilityCacheControl = "public, max-age=3600"
	defaultLimit             = 1
)

// Handlers implements the imagery operations of the API.
type Handlers struct {
	drawer *application.Drawer
}

// NewHandlers returns the imagery handlers.
func NewHandlers(drawer *application.Drawer) *Handlers {
	return &Handlers{drawer: drawer}
}

// ListZonesWithImages implements GET /decades/{decade}/zones.
func (h *Handlers) ListZonesWithImages(ctx context.Context, request apigen.ListZonesWithImagesRequestObject) (apigen.ListZonesWithImagesResponseObject, error) {
	decade, err := timespan.NewDecade(request.Decade, time.Now())
	if err != nil {
		return apigen.ListZonesWithImages400ApplicationProblemPlusJSONResponse{
			ProblemApplicationProblemPlusJSONResponse: apigen.ProblemApplicationProblemPlusJSONResponse(problem.New(http.StatusBadRequest, err.Error())),
		}, nil
	}

	zones, err := h.drawer.ZonesWithImages(ctx, decade)
	if errors.Is(err, domain.ErrSourceUnavailable) {
		return apigen.ListZonesWithImages503ApplicationProblemPlusJSONResponse(problem.New(http.StatusServiceUnavailable, "Photos are temporarily unavailable.")), nil
	}
	if err != nil {
		return nil, err
	}

	ids := make([]apigen.ZoneId, 0, len(zones))
	for _, zone := range zones {
		ids = append(ids, apigen.ZoneId(zone))
	}
	cacheControl := availabilityCacheControl
	return apigen.ListZonesWithImages200JSONResponse{
		Body:    apigen.DecadeAvailability{Decade: int(decade), Zones: ids},
		Headers: apigen.ListZonesWithImages200ResponseHeaders{CacheControl: &cacheControl},
	}, nil
}

// DrawImages implements GET /images.
func (h *Handlers) DrawImages(ctx context.Context, request apigen.DrawImagesRequestObject) (apigen.DrawImagesResponseObject, error) {
	decade, err := timespan.NewDecade(request.Params.Decade, time.Now())
	if err != nil {
		return badDrawRequest(err.Error()), nil
	}
	zone, known := zoning.Find(zoning.ID(request.Params.Zone))
	if !known {
		return badDrawRequest("unknown zone " + string(request.Params.Zone)), nil
	}
	limit := defaultLimit
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}
	if limit < 1 || limit > 5 {
		return badDrawRequest("limit must be between 1 and 5"), nil
	}

	images, err := h.drawer.Draw(ctx, zone.ID, decade, limit)
	if errors.Is(err, domain.ErrSourceUnavailable) {
		return apigen.DrawImages503ApplicationProblemPlusJSONResponse(problem.New(http.StatusServiceUnavailable, "Photos are temporarily unavailable.")), nil
	}
	if err != nil {
		return nil, err
	}

	items := make([]apigen.Image, 0, len(images))
	for _, image := range images {
		items = append(items, toAPIImage(image))
	}
	return apigen.DrawImages200JSONResponse{Items: items}, nil
}

func badDrawRequest(detail string) apigen.DrawImages400ApplicationProblemPlusJSONResponse {
	return apigen.DrawImages400ApplicationProblemPlusJSONResponse{
		ProblemApplicationProblemPlusJSONResponse: apigen.ProblemApplicationProblemPlusJSONResponse(problem.New(http.StatusBadRequest, detail)),
	}
}
