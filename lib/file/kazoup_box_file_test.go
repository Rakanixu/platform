package file

import (
	"github.com/kazoup/platform/lib/rossete"
	"github.com/tjarratt/babble"
	"reflect"
	"testing"
	"time"
)

var (
	babbler       = babble.NewBabbler()
	kazoupBoxFile = &KazoupBoxFile{
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

func TestKazoupBoxFile(t *testing.T) {
	var _ File = (*KazoupBoxFile)(nil)
}

func TestKazoupBoxFile_PreviewURL(t *testing.T) {
	expected := "https://api.box.com/2.0/files/original_id/thumbnail.png?min_height=256&min_width=256"
	result := kazoupBoxFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupBoxFile_GetID(t *testing.T) {
	result := kazoupBoxFile.GetID()

	if kazoupBoxFile.KazoupFile.ID != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.ID, result)
	}
}

func TestKazoupBoxFile_GetName(t *testing.T) {
	result := kazoupBoxFile.GetName()

	if kazoupBoxFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.Name, result)
	}
}

func TestKazoupBoxFile_GetUserID(t *testing.T) {
	result := kazoupBoxFile.GetUserID()

	if kazoupBoxFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.UserId, result)
	}
}

func TestKazoupBoxFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupBoxFile.GetIDFromOriginal()

	if kazoupBoxFile.OriginalID != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.OriginalID, result)
	}
}

func TestKazoupBoxFile_GetIndex(t *testing.T) {
	result := kazoupBoxFile.GetIndex()

	if kazoupBoxFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.Index, result)
	}
}

func TestKazoupBoxFile_GetDatasourceID(t *testing.T) {
	result := kazoupBoxFile.GetDatasourceID()

	if kazoupBoxFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupBoxFile_GetFileType(t *testing.T) {
	result := kazoupBoxFile.GetFileType()

	if kazoupBoxFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.FileType, result)
	}
}

func TestKazoupBoxFile_GetURL(t *testing.T) {
	result := kazoupBoxFile.GetURL()

	if kazoupBoxFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.URL, result)
	}
}

func TestKazoupBoxFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupBoxFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupBoxFile_GetModifiedTime(t *testing.T) {
	result := kazoupBoxFile.GetModifiedTime()

	if kazoupBoxFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.Modified, result)
	}
}

func TestKazoupBoxFile_GetContent(t *testing.T) {
	result := kazoupBoxFile.GetContent()

	if kazoupBoxFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.Content, result)
	}
}

func TestKazoupBoxFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupBoxFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupBoxFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupBoxFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupBoxFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupBoxFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupBoxFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupBoxFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupBoxFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupBoxFile.SetHighlight(highlight)

	if highlight != kazoupBoxFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupBoxFile.KazoupFile.Highlight)
	}
}

func TestKazoupBoxFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupBoxFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupBoxFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupBoxFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupBoxFile_SetEntities(t *testing.T) {
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

	kazoupBoxFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupBoxFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupBoxFile.KazoupFile.Entities)
	}
}

func TestKazoupBoxFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupBoxFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupBoxFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupBoxFile.KazoupFile.Sentiment)
	}
}
