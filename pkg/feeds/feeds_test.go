//go:generate mockgen -source=feeds.go -destination=feeds_mock_test.go -package=feeds

package feeds

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/mxpv/podsync/pkg/api"
	"github.com/mxpv/podsync/pkg/model"
	"github.com/mxpv/podsync/pkg/queue"
)

var feed = &model.Feed{
	HashID:   "123",
	ItemID:   "xyz",
	Provider: api.ProviderVimeo,
	LinkType: api.LinkTypeChannel,
	PageSize: 50,
	Quality:  api.QualityHigh,
	Format:   api.FormatVideo,
	Episodes: []*model.Item{
		{ID: "1", Title: "Title", Description: "Description"},
	},
}

func TestService_CreateFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockstorage(ctrl)
	db.EXPECT().SaveFeed(gomock.Any()).Times(1).Return(nil)

	gen, _ := NewIDGen()

	builder := NewMockBuilder(ctrl)
	builder.EXPECT().Build(gomock.Any()).Times(1).Return(nil)

	s := Service{
		generator: gen,
		storage:   db,
		builders:  map[api.Provider]Builder{api.ProviderYoutube: builder},
	}

	req := &api.CreateFeedRequest{
		URL:      "youtube.com/channel/123",
		PageSize: 50,
		Quality:  api.QualityHigh,
		Format:   api.FormatVideo,
	}

	hashID, err := s.CreateFeed(req, &api.Identity{})
	require.NoError(t, err)
	require.NotEmpty(t, hashID)
}

func TestService_makeFeed(t *testing.T) {
	req := &api.CreateFeedRequest{
		URL:      "youtube.com/channel/123",
		PageSize: 1000,
		Quality:  api.QualityLow,
		Format:   api.FormatAudio,
	}

	gen, _ := NewIDGen()

	s := Service{
		generator: gen,
	}

	feed, err := s.makeFeed(req, &api.Identity{})
	require.NoError(t, err)
	require.Equal(t, 50, feed.PageSize)
	require.Equal(t, api.QualityHigh, feed.Quality)
	require.Equal(t, api.FormatVideo, feed.Format)

	feed, err = s.makeFeed(req, &api.Identity{FeatureLevel: api.ExtendedFeatures})
	require.NoError(t, err)
	require.Equal(t, 150, feed.PageSize)
	require.Equal(t, api.QualityLow, feed.Quality)
	require.Equal(t, api.FormatAudio, feed.Format)

	feed, err = s.makeFeed(req, &api.Identity{FeatureLevel: api.ExtendedPagination})
	require.NoError(t, err)
	require.Equal(t, 600, feed.PageSize)
	require.Equal(t, api.QualityLow, feed.Quality)
	require.Equal(t, api.FormatAudio, feed.Format)
}

func TestService_QueryFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockstorage(ctrl)
	db.EXPECT().GetFeed("123").Times(1).Return(nil, nil)

	s := Service{storage: db}
	_, err := s.QueryFeed("123")
	require.NoError(t, err)
}

func TestService_BuildFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(feed.HashID).Times(1).Return(feed, nil)

	q := NewMockSender(ctrl)
	q.EXPECT().Add(gomock.Eq(&queue.Item{
		ID:       feed.HashID,
		URL:      feed.ItemURL,
		Start:    1,
		Count:    feed.PageSize,
		LastID:   feed.LastID,
		LinkType: feed.LinkType,
		Format:   string(feed.Format),
		Quality:  string(feed.Quality),
	})).Times(1)

	s := Service{storage: stor, sender: q}

	_, err := s.BuildFeed(feed.HashID)
	require.NoError(t, err)
}

func TestService_WrongID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(gomock.Any()).Times(1).Return(nil, errors.New("not found"))

	s := &Service{storage: stor}

	_, err := s.BuildFeed("invalid_feed_id")
	require.Error(t, err)
}

func TestService_GetMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetMetadata(feed.HashID).Times(1).Return(feed, nil)

	s := &Service{storage: stor}

	m, err := s.GetMetadata(feed.HashID)
	require.NoError(t, err)
	require.EqualValues(t, 0, m.Downloads)
}
