package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"log"
	"net/http"
	"strings"
)

// getDirChildren get children from directory
func (bfs *BoxFs) getDirChildren(id string, offset, limit int) error {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/?offset=%d&limit=%d", globals.BoxFoldersEndpoint, id, offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var bdc *box.BoxDirContents
	if err := json.NewDecoder(rsp.Body).Decode(&bdc); err != nil {
		return err
	}

	for _, v := range bdc.ItemCollection.Entries {
		if v.Type == "folder" {
			// Push found directories into the queue to be crawled
			bfs.Directories <- v.ID
		} else {
			// File discovered, but need to retrieve more info about the file
			if err := bfs.getMetadataFromFile(v.ID); err != nil {
				return err
			}
		}
	}

	if bdc.ItemCollection.TotalCount > bdc.ItemCollection.Offset+bdc.ItemCollection.Limit {
		bfs.getDirChildren(
			id,
			bdc.ItemCollection.Offset+bdc.ItemCollection.Limit,
			bdc.ItemCollection.Limit,
		)
	}

	return nil
}

// getMetadataFromFile retrieves more info about discovered files in box
func (bfs *BoxFs) getMetadataFromFile(id string) error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxFileMetadataEndpoint+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var fm *box.BoxFileMeta
	if err := json.NewDecoder(rsp.Body).Decode(&fm); err != nil {
		return err
	}

	f := file.NewKazoupFileFromBoxFile(*fm, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

	if err := bfs.generateThumbnail(fm, f.ID); err != nil {
		log.Println(err)
	}

	bfs.FilesChan <- NewFileMsg(f, nil)

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage with kazoup id
func (bfs *BoxFs) generateThumbnail(fm *box.BoxFileMeta, id string) error {
	name := strings.Split(fm.Name, ".")

	if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
		// Download file from Box, so connector is globals.Box
		bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
		if err != nil {
			return err
		}

		rc, err := bcs.Download(fm.ID)
		if err != nil {
			return errors.New("ERROR downloading box file")
		}

		rd, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for box file")
		}

		// Upload file to GoogleCloudStorage, so connector is globals.GoogleCloudStorage
		ncs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.GoogleCloudStorage)
		if err != nil {
			return err
		}

		if err := ncs.Upload(rd, id); err != nil {
			return err
		}
	}

	return nil
}

// token returns authorization header token as string
func (bfs *BoxFs) token() string {
	return "Bearer " + bfs.Endpoint.Token.AccessToken
}