package file

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/rossete"
	googledrive "google.golang.org/api/drive/v3"
	"reflect"
	"testing"
	"time"
)

var (
	kazoupGoogleFile = &KazoupGoogleFile{
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
		Original: &googledrive.File{
			Id:          "original_id",
			WebViewLink: "https://a.b.c",
		},
	}
)

func TestKazoupGoogleFile(t *testing.T) {
	var _ File = (*KazoupGoogleFile)(nil)
}

func TestKazoupGoogleFile_PreviewURL(t *testing.T) {
	expected := ""
	result := kazoupGoogleFile.PreviewURL("", "", "", "")

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGoogleFile_GetID(t *testing.T) {
	expected := utils.GetMD5Hash(kazoupGoogleFile.Original.WebViewLink)
	result := kazoupGoogleFile.GetID()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGoogleFile_GetName(t *testing.T) {
	result := kazoupGoogleFile.GetName()

	if kazoupGoogleFile.KazoupFile.Name != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.Name, result)
	}
}

func TestKazoupGoogleFile_GetUserID(t *testing.T) {
	result := kazoupGoogleFile.GetUserID()

	if kazoupGoogleFile.KazoupFile.UserId != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.UserId, result)
	}
}

func TestKazoupGoogleFile_GetIDFromOriginal(t *testing.T) {
	result := kazoupGoogleFile.GetIDFromOriginal()

	if kazoupGoogleFile.Original.Id != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.Original.Id, result)
	}
}

func TestKazoupGoogleFile_GetIndex(t *testing.T) {
	result := kazoupGoogleFile.GetIndex()

	if kazoupGoogleFile.KazoupFile.Index != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.Index, result)
	}
}

func TestKazoupGoogleFile_GetDatasourceID(t *testing.T) {
	result := kazoupGoogleFile.GetDatasourceID()

	if kazoupGoogleFile.KazoupFile.DatasourceId != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.DatasourceId, result)
	}
}

func TestKazoupGoogleFile_GetFileType(t *testing.T) {
	result := kazoupGoogleFile.GetFileType()

	if kazoupGoogleFile.KazoupFile.FileType != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.FileType, result)
	}
}

func TestKazoupGoogleFile_GetPathDisplay(t *testing.T) {
	result := kazoupGoogleFile.GetPathDisplay()

	if "" != result {
		t.Errorf("Expected %v, got %v", "", result)
	}
}

func TestKazoupGoogleFile_GetURL(t *testing.T) {
	result := kazoupGoogleFile.GetURL()

	if kazoupGoogleFile.KazoupFile.URL != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.URL, result)
	}
}

func TestKazoupGoogleFile_GetExtension(t *testing.T) {
	expected := "extension"
	result := kazoupGoogleFile.GetExtension()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGoogleFile_GetBase64(t *testing.T) {
	expected := ""
	result := kazoupGoogleFile.GetBase64()

	if expected != result {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestKazoupGoogleFile_GetModifiedTime(t *testing.T) {
	result := kazoupGoogleFile.GetModifiedTime()

	if kazoupGoogleFile.KazoupFile.Modified != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.Modified, result)
	}
}

func TestKazoupGoogleFile_GetContent(t *testing.T) {
	result := kazoupGoogleFile.GetContent()

	if kazoupGoogleFile.KazoupFile.Content != result {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.Content, result)
	}
}

func TestKazoupGoogleFile_GetOptsTimestamps(t *testing.T) {
	result := kazoupGoogleFile.GetOptsTimestamps()

	if !reflect.DeepEqual(kazoupGoogleFile.KazoupFile.OptsKazoupFile, result) {
		t.Errorf("Expected %v, got %v", kazoupGoogleFile.KazoupFile.OptsKazoupFile, result)
	}
}

func TestKazoupGoogleFile_SetOptsTimestamps(t *testing.T) {
	tim := time.Now()
	opts := &OptsKazoupFile{
		ContentTimestamp:           &tim,
		TagsTimestamp:              &tim,
		AudioTimestamp:             &tim,
		TextAnalyzedTimestamp:      &tim,
		SentimentAnalyzedTimestamp: &tim,
		ThumbnailTimestamp:         &tim,
	}

	kazoupGoogleFile.SetOptsTimestamps(opts)

	if !reflect.DeepEqual(opts, kazoupGoogleFile.KazoupFile.OptsKazoupFile) {
		t.Errorf("Expected %v, got %v", opts, kazoupGoogleFile.KazoupFile.OptsKazoupFile)
	}
}

func TestKazoupGoogleFile_SetHighlight(t *testing.T) {
	highlight := "highlight"

	kazoupGoogleFile.SetHighlight(highlight)

	if highlight != kazoupGoogleFile.KazoupFile.Highlight {
		t.Errorf("Expected %v, got %v", highlight, kazoupGoogleFile.KazoupFile.Highlight)
	}
}

func TestKazoupGoogleFile_SetContentCategory(t *testing.T) {
	s := "content_category"
	categorization := &KazoupCategorization{
		ContentCategory: &s,
	}

	kazoupGoogleFile.SetContentCategory(categorization)

	if !reflect.DeepEqual(categorization, kazoupGoogleFile.KazoupFile.KazoupCategorization) {
		t.Errorf("Expected %v, got %v", categorization, kazoupGoogleFile.KazoupFile.KazoupCategorization)
	}
}

func TestKazoupGoogleFile_SetEntities(t *testing.T) {
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

	kazoupGoogleFile.SetEntities(entities)

	if !reflect.DeepEqual(entities, kazoupGoogleFile.KazoupFile.Entities) {
		t.Errorf("Expected %v, got %v", entities, kazoupGoogleFile.KazoupFile.Entities)
	}
}

func TestKazoupGoogleFile_SetSentiment(t *testing.T) {
	sentiment := &rossete.RosseteSentiment{}

	sentiment.Document.Confidence = 0.55
	sentiment.Document.Label = babbler.Babble()

	kazoupGoogleFile.SetSentiment(sentiment)

	if !reflect.DeepEqual(sentiment, kazoupGoogleFile.KazoupFile.Sentiment) {
		t.Errorf("Expected %v, got %v", sentiment, kazoupGoogleFile.KazoupFile.Sentiment)
	}
}
