package file

import (
	"github.com/kazoup/platform/lib/rossete"
	"reflect"
	"testing"
	"time"
)

var (
	kazoupGmailFile = &KazoupGmailFile{
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

func TestKazoupGmailFile(t *testing.T) {
	var _ File = (*KazoupGmailFile)(nil)
}

func TestKazoupGmailFile_PreviewURL(t *testing.T) {
	expected := "https://mail.google.com/mail/u/original_id"
	result := kazoupGmailFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGmailFile_GetID(t *testing.T) {
	result := kazoupGmailFile.GetID()

	if kazoupGmailFile.KazoupFile.ID != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.ID, result)
	}
}

func TestKazoupGmailFile_GetName(t *testing.T) {
	result := kazoupGmailFile.GetName()

	if kazoupGmailFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.Name, result)
	}
}

func TestKazoupGmailFile_GetUserID(t *testing.T) {
	result := kazoupGmailFile.GetUserID()

	if kazoupGmailFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.UserId, result)
	}
}

func TestKazoupGmailFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupGmailFile.GetIDFromOriginal()

	if kazoupGmailFile.OriginalID != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.OriginalID, result)
	}
}

func TestKazoupGmailFile_GetIndex(t *testing.T) {
	result := kazoupGmailFile.GetIndex()

	if kazoupGmailFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.Index, result)
	}
}

func TestKazoupGmailFile_GetDatasourceID(t *testing.T) {
	result := kazoupGmailFile.GetDatasourceID()

	if kazoupGmailFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupGmailFile_GetFileType(t *testing.T) {
	result := kazoupGmailFile.GetFileType()

	if kazoupGmailFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.FileType, result)
	}
}

func TestKazoupGmailFile_GetURL(t *testing.T) {
	result := kazoupGmailFile.GetURL()

	if kazoupGmailFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.URL, result)
	}
}

func TestKazoupGmailFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupGmailFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGmailFile_GetModifiedTime(t *testing.T) {
	result := kazoupGmailFile.GetModifiedTime()

	if kazoupGmailFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.Modified, result)
	}
}

func TestKazoupGmailFile_GetContent(t *testing.T) {
	result := kazoupGmailFile.GetContent()

	if kazoupGmailFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.Content, result)
	}
}

func TestKazoupGmailFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupGmailFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupGmailFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupGmailFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupGmailFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupGmailFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupGmailFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupGmailFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupGmailFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupGmailFile.SetHighlight(highlight)

	if highlight != kazoupGmailFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupGmailFile.KazoupFile.Highlight)
	}
}

func TestKazoupGmailFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupGmailFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupGmailFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupGmailFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupGmailFile_SetEntities(t *testing.T) {
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

	kazoupGmailFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupGmailFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupGmailFile.KazoupFile.Entities)
	}
}

func TestKazoupGmailFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupGmailFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupGmailFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupGmailFile.KazoupFile.Sentiment)
	}
}
