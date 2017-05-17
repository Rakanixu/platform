package file

import (
	"github.com/kazoup/platform/lib/rossete"
	"reflect"
	"testing"
	"time"
)

var (
	kazoupOneDriveFile = &KazoupOneDriveFile{
		KazoupFile: KazoupFile{
			ID:           "KazoupFile_ID",
			OriginalID:   "original_id",
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
	}
)

func TestKazoupOneDriveFile(t *testing.T) {
	var _ File = (*KazoupOneDriveFile)(nil)
}

func TestKazoupOneDriveFile_PreviewURL(t *testing.T) {
	expected := kazoupOneDriveFile.PreviewUrl
	result := kazoupOneDriveFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupOneDriveFile_GetID(t *testing.T) {
	expected := kazoupOneDriveFile.ID
	result := kazoupOneDriveFile.GetID()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupOneDriveFile_GetName(t *testing.T) {
	result := kazoupOneDriveFile.GetName()

	if kazoupOneDriveFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.Name, result)
	}
}

func TestKazoupOneDriveFile_GetUserID(t *testing.T) {
	result := kazoupOneDriveFile.GetUserID()

	if kazoupOneDriveFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.UserId, result)
	}
}

func TestKazoupOneDriveFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupOneDriveFile.GetIDFromOriginal()

	if kazoupOneDriveFile.OriginalID != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.OriginalID, result)
	}
}

func TestKazoupOneDriveFile_GetIndex(t *testing.T) {
	result := kazoupOneDriveFile.GetIndex()

	if kazoupOneDriveFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.Index, result)
	}
}

func TestKazoupOneDriveFile_GetDatasourceID(t *testing.T) {
	result := kazoupOneDriveFile.GetDatasourceID()

	if kazoupOneDriveFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupOneDriveFile_GetFileType(t *testing.T) {
	result := kazoupOneDriveFile.GetFileType()

	if kazoupOneDriveFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.FileType, result)
	}
}

func TestKazoupOneDriveFile_GetURL(t *testing.T) {
	result := kazoupOneDriveFile.GetURL()

	if kazoupOneDriveFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.URL, result)
	}
}

func TestKazoupOneDriveFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupOneDriveFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupOneDriveFile_GetModifiedTime(t *testing.T) {
	result := kazoupOneDriveFile.GetModifiedTime()

	if kazoupOneDriveFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.Modified, result)
	}
}

func TestKazoupOneDriveFile_GetContent(t *testing.T) {
	result := kazoupOneDriveFile.GetContent()

	if kazoupOneDriveFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.Content, result)
	}
}

func TestKazoupOneDriveFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupOneDriveFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupOneDriveFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupOneDriveFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupOneDriveFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupOneDriveFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupOneDriveFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupOneDriveFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupOneDriveFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupOneDriveFile.SetHighlight(highlight)

	if highlight != kazoupOneDriveFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupOneDriveFile.KazoupFile.Highlight)
	}
}

func TestKazoupOneDriveFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupOneDriveFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupOneDriveFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupOneDriveFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupOneDriveFile_SetEntities(t *testing.T) {
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

	kazoupOneDriveFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupOneDriveFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupOneDriveFile.KazoupFile.Entities)
	}
}

func TestKazoupOneDriveFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupOneDriveFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupOneDriveFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupOneDriveFile.KazoupFile.Sentiment)
	}
}
