package web

import (
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"net/http"
	"testing"
)

func TestCreateSource_succeeds(t *testing.T) {
	handler := NewSourceHandler(func(name string, path string) error {
		assert.Equal(t, "name1", name)
		assert.Equal(t, "path1", path)
		return nil
	}, nil, nil, nil)

	response, err := send(handler, "/api/v1/sources", http.MethodPost, &createSourceRequest{Name: "name1", Path: "path1"})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestFetchAllSources_succeeds(t *testing.T) {
	source1 := domain.Source{ID: "id1", Name: "name1", Path: "path1"}
	source2 := domain.Source{ID: "id2", Name: "name2", Path: "path2"}
	handler := NewSourceHandler(nil, func() ([]domain.Source, error) {
		return []domain.Source{source1, source2}, nil
	}, nil, nil)

	response, err := send(handler, "/api/v1/sources", http.MethodGet, nil)
	assert.Nil(t, err)

	var body getAllSourcesResponse
	parseResponse(t, response, &body)

	assertResponse := func(source domain.Source, response getSourceResponse) {
		assert.EqualValues(t, source.ID, response.ID)
		assert.EqualValues(t, source.Name, response.Name)
		assert.EqualValues(t, source.Path, response.Path)
	}
	assertResponse(source1, body.Sources[0])
	assertResponse(source2, body.Sources[1])
}

func TestSyncSource_succeeds(t *testing.T) {
	handler := NewSourceHandler(nil, nil, nil, func(sourceID string) error {
		assert.Equal(t, "sourceID", sourceID)
		return nil
	})

	response, err := send(handler, "/api/v1/sources/sourceID/sync", http.MethodPut, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDeleteSource_succeeds(t *testing.T) {
	handler := NewSourceHandler(nil, nil, func(id string) error {
		assert.Equal(t, "sourceID", id)
		return nil
	}, nil)

	response, err := send(handler, "/api/v1/sources/sourceID/", http.MethodDelete, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
