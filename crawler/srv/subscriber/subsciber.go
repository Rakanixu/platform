package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	crawler_proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_conn "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
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

	// Update token in DB
	if err := db_conn.UpdateFileSystemAuth(c.Client, globals.NewSystemContext(), endpoint.Id, auth); err != nil {
		return err
	}

	// Publish crawler started, or is just going to start..
	if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.CrawlerStartedTopic, &crawler_proto.CrawlerStartedMessage{
		UserId:       endpoint.UserId,
		DatasourceId: endpoint.Id,
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

				// Clear index (files that no longer exists, rename, etc..)
				if err := globals.ClearIndex(c.Client, endpoint); err != nil {
					log.Println("ERROR clearing index after scan", err)
				}

				msg := &crawler.CrawlerFinishedMessage{
					DatasourceId: endpoint.Id,
				}
				// Publish crawling process has finished
				if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
					//return err
					log.Println("Error publishin crawler finished", err)
				}
				close(filesChan)
				close(filesRunning)
				wg.Done()
				return
			// Channel receives File to be indexed by Elastic Search
			case fc := <-filesChan:
				if fc.Error != nil {
					log.Println("Error discovering file", fc.Error)
				}

				if err := file.IndexAsync(c.Client, fc.File, globals.FilesTopic, fc.File.GetIndex(), false); err != nil {
					log.Println("Error indexing async file", err)
				}
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

				if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.SlackUsersTopic, u.User)); err != nil {
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

				if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.SlackChannelsTopic, ch.Channel)); err != nil {
					log.Println("Error indexing channel", err)
				}
			}
		}
	}()

	wg.Wait() // Wait for all listeners

	return nil
}
