package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/bulk"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

type TaskHandler struct{}

// Scans subscriber, receive endpoint to crawl it
func (t *TaskHandler) Scans(ctx context.Context, endpoint *proto_datasource.Endpoint) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

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
	_, err = operations.Update(ctx, &proto_operations.UpdateRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	// Publish notification
	if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
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
	wg.Add(3) // Wait for files, users and channels discovery

	// Shared error channel
	errs := make(chan error, 10)

	// Handle errors
	go handleDiscoveryErrors(errs)

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
					errs <- errors.NewDiscoveryError(fc.File, "crawler", fc.Error)
				} else {
					b, err := json.Marshal(fc.File)
					if err != nil {
						errs <- err
					} else {
						if err := bulk.Files(ctx, &crawler.FileMessage{
							Id:     fc.File.GetID(),
							UserId: fc.File.GetUserID(),
							Index:  fc.File.GetIndex(),
							Notify: false,
							Data:   string(b),
						}); err != nil {
							errs <- errors.NewDiscoveryError(fc.File, "crawler bulk file", err)
						}
					}
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
					errs <- errors.NewDiscoveryError(u.User, "crawler", u.Error)
				} else {
					if err := bulk.SlackUsers(ctx, u.User); err != nil {
						errs <- errors.NewDiscoveryError(u.User, "crawler bulk user", err)
					}
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
					errs <- errors.NewDiscoveryError(ch.Channel, "crawler", ch.Error)
				} else {
					if err := bulk.SlackChannels(ctx, ch.Channel); err != nil {
						errs <- errors.NewDiscoveryError(ch.Channel, "crawler bulk channel", err)
					}
				}
			}
		}
	}()

	wg.Wait() // Wait for all listeners

	close(errs)

	return nil
}

func handleDiscoveryErrors(errs chan error) {
	for err := range errs {
		log.Printf("Discovery error: %v", err)
	}
}
