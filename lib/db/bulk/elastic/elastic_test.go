package elastic

import (
	//"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	kazoup_context "github.com/kazoup/platform/lib/context"
	//"github.com/kazoup/platform/lib/file"
	//"github.com/kazoup/platform/lib/slack"
	"github.com/kazoup/platform/lib/wrappers"
	//"github.com/micro/go-micro"
	"golang.org/x/net/context"
	//elib "gopkg.in/olivere/elastic.v5"
	"testing"
)

const (
	TEST_USER_ID        = "test_user"
	TEST_INDEX_FILES    = "index_test_files"
	TEST_INDEX_USERS    = "index_test_users"
	TEST_INDEX_CHANNELS = "index_test_channels"
)

var (
	e = &elastic{
		FilesChannel:         make(chan *filesChannel),
		SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
		SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
	}
	srv = wrappers.NewKazoupService("test")
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestElastic_Init(t *testing.T) {
	result := e.Init(srv)

	if nil != result {
		t.Errorf("Expected %v, got: %v", nil, result)
	}

	if !e.Client.IsRunning() {
		t.Error("Client not running")
	}
}

/*

func TestElastic_Files(t *testing.T) {
	result := e.Files(ctx, &crawler.FileMessage{})

	if nil != result {
		t.Error("FilesChannel blocked")
	}
}

func TestElastic_SlackUsers(t *testing.T) {
	result := e.SlackUsers(ctx, &crawler.SlackUserMessage{})

	if nil != result {
		t.Error("SlackUsersChannel blocked")
	}
}

func TestElastic_SlackChannels(t *testing.T) {
	result := e.SlackChannels(ctx, &crawler.SlackChannelMessage{})

	if nil != result {
		t.Error("SlackChannelsChannel blocked")
	}
}

func Test_processFiles(t *testing.T) {
	f := file.NewKazoupFileFromMockFile()
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	e = &elastic{
		FilesChannel:         make(chan *filesChannel),
		SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
		SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
	}
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL("http://elasticsearch:9200"),
		elib.SetBasicAuth("", ""),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		t.Fatal(err)
	}
	e.BulkFilesProcessor, err = e.Client.BulkProcessor().
		Workers(1).
		BulkActions(1).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	testData := []struct {
		filesChannel *filesChannel
	}{
		{
			&filesChannel{
				FileMessage: &crawler.FileMessage{
					Id:     "1",
					UserId: TEST_USER_ID,
					Index:  TEST_INDEX_FILES,
					Notify: true,
					Data:   string(b),
				},
				Ctx: micro.NewContext(ctx, srv),
			},
		},
		{
			&filesChannel{
				FileMessage: &crawler.FileMessage{
					Id:     "2",
					UserId: TEST_USER_ID,
					Index:  TEST_INDEX_FILES,
					Data:   string(b),
				},
				Ctx: micro.NewContext(ctx, srv),
			},
		},
		{
			&filesChannel{
				FileMessage: &crawler.FileMessage{
					Id:     "3",
					UserId: TEST_USER_ID,
					Index:  TEST_INDEX_FILES,
					Data:   string(b),
				},
				Ctx: micro.NewContext(ctx, srv),
			},
		},
		{
			&filesChannel{
				FileMessage: &crawler.FileMessage{
					Id:     "4",
					UserId: TEST_USER_ID,
					Index:  TEST_INDEX_FILES,
					Data:   string(b),
				},
				Ctx: micro.NewContext(ctx, srv),
			},
		},
		{
			&filesChannel{
				FileMessage: &crawler.FileMessage{
					Id:     "5",
					UserId: TEST_USER_ID,
					Index:  TEST_INDEX_FILES,
					Data:   string(b),
				},
				Ctx: micro.NewContext(ctx, srv),
			},
		},
	}

	go processFiles(e)

	for _, tt := range testData {
		e.FilesChannel <- tt.filesChannel
	}

	// Data may not be yet, but is eventually consistent
*/
/*	time.Sleep(time.Millisecond * 100)*/ /*

 */
/*	count, err := e.Client.Count(TEST_INDEX_FILES).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if int64(len(testData)) != count {
		t.Errorf("Expected: %v, got: %v", len(testData), count)
	}*/ /*


	_, err = e.Client.DeleteIndex(TEST_INDEX_FILES).Do(context.TODO())
	if err != nil {
		t.Error("Error cleaning up index from running ElasticSearch instance")
	}
}

func Test_processSlackUsers(t *testing.T) {
	u := &slack.ESSlackUser{
		UserID:   TEST_USER_ID,
		TeamID:   "team_id",
		UserName: "user_name",
	}
	b, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}

	e = &elastic{
		FilesChannel:         make(chan *filesChannel),
		SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
		SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
	}
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL("http://elasticsearch:9200"),
		elib.SetBasicAuth("", ""),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		t.Fatal(err)
	}
	e.BulkProcessor, err = e.Client.BulkProcessor().
		Workers(1).
		BulkActions(1).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	testData := []struct {
		userChannel *crawler.SlackUserMessage
	}{
		{
			&crawler.SlackUserMessage{
				Id:    "1",
				Index: TEST_INDEX_USERS,
				Data:  string(b),
			},
		},
		{
			&crawler.SlackUserMessage{
				Id:    "2",
				Index: TEST_INDEX_USERS,
				Data:  string(b),
			},
		},
	}

	go processSlackUsers(e)

	for _, tt := range testData {
		e.SlackUsersChannel <- tt.userChannel
	}

	// Data may not be yet, but is eventually consistent
	// time.Sleep(time.Millisecond * 100)
*/
/*
	count, err := e.Client.Count(TEST_INDEX_USERS).Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if int64(len(testData)) != count {
		t.Errorf("Expected: %v, got: %v", len(testData), count)
	}*/ /*


 */
/*	_, err = e.Client.DeleteIndex(TEST_INDEX_USERS).Do(context.TODO())
	if err != nil {
		t.Error("Error cleaning up index from running ElasticSearch instance", err)
	}*/ /*

}

func Test_processSlackChannels(t *testing.T) {
	u := &slack.ESSlackChannel{
		ChannelID:   "channel_id",
		ChannelName: "channel_name",
	}
	b, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}

	e = &elastic{
		FilesChannel:         make(chan *filesChannel),
		SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
		SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
	}
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL("http://elasticsearch:9200"),
		elib.SetBasicAuth("", ""),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		t.Fatal(err)
	}
	e.BulkProcessor, err = e.Client.BulkProcessor().
		Workers(1).
		BulkActions(1).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	testData := []struct {
		channelChannel *crawler.SlackChannelMessage
	}{
		{
			&crawler.SlackChannelMessage{
				Id:    "1",
				Index: TEST_INDEX_CHANNELS,
				Data:  string(b),
			},
		},
		{
			&crawler.SlackChannelMessage{
				Id:    "2",
				Index: TEST_INDEX_CHANNELS,
				Data:  string(b),
			},
		},
	}

	go processSlackChannels(e)

	for _, tt := range testData {
		e.SlackChannelsChannel <- tt.channelChannel
	}
}
*/
