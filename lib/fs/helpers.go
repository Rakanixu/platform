package fs

import (
	crawler_proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/file"
)

type FileMeta struct {
	File  file.File
	Error error
}

func NewFileMeta(file file.File, err error) FileMeta {
	return FileMeta{
		File:  file,
		Error: err,
	}
}

// UserMsg is a helper for discovering users in file system
// Right now is just for slack users, but can be abstracted easily by creating users lib (interface the User)
// This can be done in the same way that FileMeta
type UserMsg struct {
	User  *crawler_proto.SlackUserMessage
	Error error
}

func NewUserMsg(user *crawler_proto.SlackUserMessage, err error) UserMsg {
	return UserMsg{
		User:  user,
		Error: err,
	}
}

// ChannelMsg is a helper for discovering channels in file system
// As UserMsg, can be abstracted to use a interface in the future
type ChannelMsg struct {
	Channel *crawler_proto.SlackChannelMessage
	Error   error
}

func NewChannelMsg(channel *crawler_proto.SlackChannelMessage, err error) ChannelMsg {
	return ChannelMsg{
		Channel: channel,
		Error:   err,
	}
}
