package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"path/filepath"
	"testing"
)

func TestCreateSource_succeeds(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		path := "../"
		absolutePath, err := filepath.Abs(path)
		require.Nil(t, err)

		repository := mock_domain.NewMockSourceRepository(controller)
		repository.EXPECT().Insert(gomock.Any()).DoAndReturn(func(source domain.Source) (string, error) {
			assert.Equal(t, "source1", source.Name)
			assert.Equal(t, absolutePath, source.Path)
			return "id", nil
		})

		err = CreateSourceProvider(repository)("source1", path)
		assert.Nil(t, err)
	})
}

func TestCreateSource_failsForInvalidPath(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		pathToFile := "../managesources_test.go"
		nonExistentPath := "/nonexistent"

		err := CreateSourceProvider(nil)("source1", pathToFile)
		assert.Equal(t, ErrInvalidPath, err)
		err = CreateSourceProvider(nil)("source1", nonExistentPath)
		assert.Equal(t, ErrInvalidPath, err)
	})
}

func TestFetchAllSources(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		source1 := domain.Source{Name: "source1"}
		source2 := domain.Source{Name: "source2"}

		repository := mock_domain.NewMockSourceRepository(controller)
		repository.EXPECT().FindAll().Return([]domain.Source{source1, source2}, nil)

		sources, err := FetchAllSourcesProvider(repository)()
		assert.Nil(t, err)
		assert.Contains(t, sources, source1)
		assert.Contains(t, sources, source2)
	})
}

func TestDeleteSource(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		repository := mock_domain.NewMockSourceRepository(controller)
		repository.EXPECT().DeleteByID("id1").Return(nil)
		err := DeleteSourceProvider(repository)("id1")
		assert.Nil(t, err)
	})
}
