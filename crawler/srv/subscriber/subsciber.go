package subscriber

import (
	"encoding/json"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_conn "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

type Crawler struct {
	Client client.Client
}

// Scans subscriber, receive endpoint to crawl it
func (c *Crawler) Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	cfs, err := fs.NewFsFromEndpoint(endpoint)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := cfs.Authorize()
	if err != nil {
		return err
	}

	// Set time for starting scan, crawler running  and update datasource
	endpoint.Token = auth
	endpoint.CrawlerRunning = true
	endpoint.LastScanStarted = time.Now().Unix()
	b, err := json.Marshal(endpoint)
	if err != nil {
		return err
	}

	// Update token in DB
	_, err = db_conn.UpdateFromDB(c.Client, ctx, &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	// Publish notification
	if err := c.Client.Publish(ctx, c.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
		Info:   "Scan started on " + endpoint.Url + " datasource.",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: endpoint.UserId,
		Data:   string(b),
	})); err != nil {
		return err
	}

	// Receive users
	usersChan, usersRunning := cfs.WalkUsers()

	// Receive channels
	channelsChan, channelsRunning := cfs.WalkChannels()

	// Receive files founded by FileSystem
	filesChan, filesRunning := cfs.Walk()

	var wg sync.WaitGroup
	wg.Add(3) //Wait for our 3 goroutines (file system listeners)

	// Listen to files
	go func() {
		for {
			select {
			// Channel receives signal cralwer has finished
			case <-filesRunning:
				time.Sleep(time.Second * 8)
				close(filesChan)
				close(filesRunning)
				wg.Done()
				return
			// Channel receives File to be indexed by Elastic Search
			case fc := <-filesChan:
				if fc.Error != nil {
					log.Println("Error discovering file", fc.Error)
				}

				if err := file.IndexAsync(ctx, c.Client, fc.File, globals.FilesTopic, fc.File.GetIndex(), false); err != nil {
					log.Println("Error indexing async file", err)
				}

				time.Sleep(globals.DISCOVERY_DELAY_MS)
			}
		}
	}()

	// Listen to users
	go func() {
		for {
			select {
			// Signal walk users finished
			case <-usersRunning:
				close(usersChan)
				close(usersRunning)
				wg.Done()
				return
			// Channel receives users to be indexed
			case u := <-usersChan:
				if u.Error != nil {
					log.Println("Error discovering user", u.Error)
				}

				if err := c.Client.Publish(ctx, c.Client.NewPublication(globals.SlackUsersTopic, u.User)); err != nil {
					log.Println("Error indexing user", err)
				}
			}
		}
	}()

	// Listen to channels
	go func() {
		for {
			select {
			// Signal walk channels finished
			case <-channelsRunning:
				close(channelsChan)
				close(channelsRunning)
				wg.Done()
				return
			// Channel receives channels to be indexed
			case ch := <-channelsChan:
				if ch.Error != nil {
					log.Println("Error discovering channel", ch.Error)
				}

				if err := c.Client.Publish(ctx, c.Client.NewPublication(globals.SlackChannelsTopic, ch.Channel)); err != nil {
					log.Println("Error indexing channel", err)
				}
			}
		}
	}()

	wg.Wait() // Wait for all listeners

	return nil
}
