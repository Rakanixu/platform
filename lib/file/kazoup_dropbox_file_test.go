package file

import (
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/rossete"
	"reflect"
	"testing"
	"time"
)

var (
	kazoupDropboxFile = &KazoupDropboxFile{
		KazoupFile: KazoupFile{
			ID:           "KazoupFile_ID",
			UserId:       babbler.Babble(),
			Name:         "KazoupFile_name.extension",
			URL:          babbler.Babble(),
			Modified:     time.Now(),
			FileSize:     1000,
			IsDir:        false,
			Category:     babbler.Babble(),
			MimeType:     babbler.Babble(),
			Depth:        0,
			FileType:     babbler.Babble(),
			LastSeen:     0,
			Access:       babbler.Babble(),
			DatasourceId: babbler.Babble(),
			Index:        babbler.Babble(),
			Content:      babbler.Babble(),
			Highlight:    babbler.Babble(),
		},
		Original: &dropbox.DropboxFile{
			ID: "original_id",
		},
	}
)

func TestKazoupDropboxFile(t *testing.T) {
	var _ File = (*KazoupDropboxFile)(nil)
}

func TestKazoupDropboxFile_PreviewURL(t *testing.T) {
	expected := ""
	result := kazoupDropboxFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupDropboxFile_GetID(t *testing.T) {
	result := kazoupDropboxFile.GetID()

	if kazoupDropboxFile.KazoupFile.ID != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.ID, result)
	}
}

func TestKazoupDropboxFile_GetName(t *testing.T) {
	result := kazoupDropboxFile.GetName()

	if kazoupDropboxFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.Name, result)
	}
}

func TestKazoupDropboxFile_GetUserID(t *testing.T) {
	result := kazoupDropboxFile.GetUserID()

	if kazoupDropboxFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.UserId, result)
	}
}

func TestKazoupDropboxFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupDropboxFile.GetIDFromOriginal()

	if kazoupDropboxFile.Original.ID != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.Original.ID, result)
	}
}

func TestKazoupDropboxFile_GetIndex(t *testing.T) {
	result := kazoupDropboxFile.GetIndex()

	if kazoupDropboxFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.Index, result)
	}
}

func TestKazoupDropboxFile_GetDatasourceID(t *testing.T) {
	result := kazoupDropboxFile.GetDatasourceID()

	if kazoupDropboxFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupDropboxFile_GetFileType(t *testing.T) {
	result := kazoupDropboxFile.GetFileType()

	if kazoupDropboxFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.FileType, result)
	}
}

func TestKazoupDropboxFile_GetPathDisplay(t *testing.T) {
	result := kazoupDropboxFile.GetPathDisplay()

	if "" != result {
		t.Errorf("Expected %v, got %v", "", result)
	}
}

func TestKazoupDropboxFile_GetURL(t *testing.T) {
	result := kazoupDropboxFile.GetURL()

	if kazoupDropboxFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.URL, result)
	}
}

func TestKazoupDropboxFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupDropboxFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupDropboxFile_GetBase64(t *testing.T) {
	expected := ""
	result := kazoupDropboxFile.GetBase64()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupDropboxFile_GetModifiedTime(t *testing.T) {
	result := kazoupDropboxFile.GetModifiedTime()

	if kazoupDropboxFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.Modified, result)
	}
}

func TestKazoupDropboxFile_GetContent(t *testing.T) {
	result := kazoupDropboxFile.GetContent()

	if kazoupDropboxFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.Content, result)
	}
}

func TestKazoupDropboxFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupDropboxFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupDropboxFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupDropboxFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupDropboxFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupDropboxFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupDropboxFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupDropboxFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupDropboxFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupDropboxFile.SetHighlight(highlight)

	if highlight != kazoupDropboxFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupDropboxFile.KazoupFile.Highlight)
	}
}

func TestKazoupDropboxFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupDropboxFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupDropboxFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupDropboxFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupDropboxFile_SetEntities(t *testing.T) {
	entities := &rossete.RosseteEntities{
		Entities: []rossete.RosseteEntity{
			{
				Type:       babbler.Babble(),
				Mention:    babbler.Babble(),
				Normalized: babbler.Babble(),
				Count:      1,
				EntityID:   babbler.Babble(),
			},
		},
	}

	kazoupDropboxFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupDropboxFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupDropboxFile.KazoupFile.Entities)
	}
}

func TestKazoupDropboxFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupDropboxFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupDropboxFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupDropboxFile.KazoupFile.Sentiment)
	}
}
