package builders

import (
	"context"
	"os"
	"testing"

	itunes "github.com/mxpv/podcast"
	"github.com/mxpv/podsync/pkg/api"
	"github.com/mxpv/podsync/pkg/model"
	"github.com/stretchr/testify/require"
)

var (
	vimeoKey = os.Getenv("VIMEO_TEST_API_KEY")
)

func TestQueryVimeoChannel(t *testing.T) {
	builder, err := NewVimeoBuilder(context.Background(), vimeoKey)
	require.NoError(t, err)

	podcast, err := builder.queryChannel(&model.Feed{ItemID: "staffpicks", Quality: api.QualityHigh})
	require.NoError(t, err)

	require.Equal(t, "https://vimeo.com/channels/staffpicks", podcast.Link)
	require.Equal(t, "Vimeo Staff Picks", podcast.Title)
	require.Equal(t, "Vimeo Curation", podcast.IAuthor)
	require.NotEmpty(t, podcast.Description)
	require.NotEmpty(t, podcast.Image)
	require.NotEmpty(t, podcast.IImage)
}

func TestQueryVimeoGroup(t *testing.T) {
	builder, err := NewVimeoBuilder(context.Background(), vimeoKey)
	require.NoError(t, err)

	podcast, err := builder.queryGroup(&model.Feed{ItemID: "motion", Quality: api.QualityHigh})
	require.NoError(t, err)

	require.Equal(t, "https://vimeo.com/groups/motion", podcast.Link)
	require.Equal(t, "Motion Graphic Artists", podcast.Title)
	require.Equal(t, "Danny Garcia", podcast.IAuthor)
	require.NotEmpty(t, podcast.Description)
	require.NotEmpty(t, podcast.Image)
	require.NotEmpty(t, podcast.IImage)
}

func TestQueryVimeoUser(t *testing.T) {
	builder, err := NewVimeoBuilder(context.Background(), vimeoKey)
	require.NoError(t, err)

	podcast, err := builder.queryUser(&model.Feed{ItemID: "motionarray", Quality: api.QualityHigh})
	require.NoError(t, err)

	require.Equal(t, "https://vimeo.com/motionarray", podcast.Link)
	require.Equal(t, "Motion Array", podcast.Title)
	require.Equal(t, "Motion Array", podcast.IAuthor)
	require.NotEmpty(t, podcast.Description)
}

func TestQueryVimeoVideos(t *testing.T) {
	builder, err := NewVimeoBuilder(context.Background(), vimeoKey)
	require.NoError(t, err)

	feed := &itunes.Podcast{}

	err = builder.queryVideos(builder.client.Channels.ListVideo, feed, &model.Feed{ItemID: "staffpicks"})
	require.NoError(t, err)

	require.Equal(t, vimeoDefaultPageSize, len(feed.Items))

	for _, item := range feed.Items {
		require.NotEmpty(t, item.Title)
		require.NotEmpty(t, item.Link)
		require.NotEmpty(t, item.GUID)
		require.NotEmpty(t, item.IDuration)
		require.NotNil(t, item.Enclosure)
		require.NotEmpty(t, item.Enclosure.URL)
		require.True(t, item.Enclosure.Length > 0)
	}
}
