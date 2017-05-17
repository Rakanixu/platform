package file

import (
	"github.com/kazoup/platform/lib/rossete"
	"reflect"
	"testing"
	"time"
)

var (
	kazoupSlackFile = &KazoupSlackFile{
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
	}
)

func TestKazoupSlackFile(t *testing.T) {
	var _ File = (*KazoupSlackFile)(nil)
}

func TestKazoupSlackFile_PreviewURL(t *testing.T) {
	expected := ""
	result := kazoupSlackFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupSlackFile_GetID(t *testing.T) {
	result := kazoupSlackFile.GetID()

	if kazoupSlackFile.KazoupFile.ID != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.ID, result)
	}
}

func TestKazoupSlackFile_GetName(t *testing.T) {
	result := kazoupSlackFile.GetName()

	if kazoupSlackFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.Name, result)
	}
}

func TestKazoupSlackFile_GetUserID(t *testing.T) {
	result := kazoupSlackFile.GetUserID()

	if kazoupSlackFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.UserId, result)
	}
}

func TestKazoupSlackFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupSlackFile.GetIDFromOriginal()

	if kazoupSlackFile.OriginalID != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.OriginalID, result)
	}
}

func TestKazoupSlackFile_GetIndex(t *testing.T) {
	result := kazoupSlackFile.GetIndex()

	if kazoupSlackFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.Index, result)
	}
}

func TestKazoupSlackFile_GetDatasourceID(t *testing.T) {
	result := kazoupSlackFile.GetDatasourceID()

	if kazoupSlackFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupSlackFile_GetFileType(t *testing.T) {
	result := kazoupSlackFile.GetFileType()

	if kazoupSlackFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.FileType, result)
	}
}

func TestKazoupSlackFile_GetURL(t *testing.T) {
	result := kazoupSlackFile.GetURL()

	if kazoupSlackFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.URL, result)
	}
}

func TestKazoupSlackFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupSlackFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupSlackFile_GetModifiedTime(t *testing.T) {
	result := kazoupSlackFile.GetModifiedTime()

	if kazoupSlackFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.Modified, result)
	}
}

func TestKazoupSlackFile_GetContent(t *testing.T) {
	result := kazoupSlackFile.GetContent()

	if kazoupSlackFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.Content, result)
	}
}

func TestKazoupSlackFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupSlackFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupSlackFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupSlackFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupSlackFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupSlackFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupSlackFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupSlackFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupSlackFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupSlackFile.SetHighlight(highlight)

	if highlight != kazoupSlackFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupSlackFile.KazoupFile.Highlight)
	}
}

func TestKazoupSlackFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupSlackFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupSlackFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupSlackFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupSlackFile_SetEntities(t *testing.T) {
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

	kazoupSlackFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupSlackFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupSlackFile.KazoupFile.Entities)
	}
}

func TestKazoupSlackFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupSlackFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupSlackFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupSlackFile.KazoupFile.Sentiment)
	}
}
