package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"log"
)

// SlackFs slack file system
type SlackFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

// NewSlackFsFromEndpoint constructor
func NewSlackFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &SlackFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover files in slack datasource
func (sfs *SlackFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := sfs.getFiles(1); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		sfs.WalkRunning <- false
	}()

	return sfs.FilesChan, sfs.WalkRunning
}

// WalUsers discover users in slack
func (sfs *SlackFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		sfs.getUsers()

		// Slack user scan finished
		sfs.WalkUsersRunning <- false
	}()

	return sfs.UsersChan, sfs.WalkUsersRunning
}

// WalChannels discover channels in slack
func (sfs *SlackFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		sfs.getChannels()

		// Slack channels scan finished
		sfs.WalkChannelsRunning <- false
	}()

	return sfs.ChannelsChan, sfs.WalkChannelsRunning
}

// CreateFile belongs to Fs interface
func (sfs *SlackFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	return sfs.FilesChan
}

// DeleteFile deletes a slack file
func (sfs *SlackFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	return sfs.FilesChan
}

// ShareFile sets a PermalinkPublic available, so everyone with URL has access to the slack file
func (sfs *SlackFs) Update(req file_proto.ShareRequest) chan FileMsg {
	/*	if req.SharePublicly {
			return sfs.shareFilePublicly(req.OriginalId)
		} else {
			r := c.NewRequest(
				globals.DB_SERVICE_NAME,
				"DB.Read",
				&db_proto.ReadRequest{
					Index: req.Index,
					Type:  "file",
					Id:    req.FileId,
				},
			)
			rsp := &db_proto.ReadResponse{}
			if err := c.Call(ctx, r, rsp); err != nil {
				return "", err
			}

			var f *file.KazoupSlackFile
			if err := json.Unmarshal([]byte(rsp.Result), &f); err != nil {
				return "", err
			}

			return sfs.shareFileInsideTeam(f, req.DestinationId)
		}*/
	return sfs.FilesChan
}
