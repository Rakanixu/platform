package speechtotext

import (
	"bytes"
	speech "cloud.google.com/go/speech/apiv1beta1"
	"fmt"
	"github.com/golang/protobuf/proto"
	normalize_text "github.com/kazoup/platform/lib/normalization/text"
	"golang.org/x/net/context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type SpeechToText interface {
	Content() string
}

type SpeechToTextContent struct {
	AudioContent string
}

func (stt *SpeechToTextContent) Content() string {
	return stt.AudioContent
}

func AsyncContent(uri string) (SpeechToText, error) {
	// Creates a client.
	client, err := speech.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	aop, err := client.AsyncRecognize(context.Background(), &speechpb.AsyncRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			SampleRate:   8000,
			LanguageCode: "en-GB",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: uri,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	sttc := &SpeechToTextContent{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var op *longrunningpb.Operation
		opC := longrunningpb.NewOperationsClient(client.Connection())

		for {
			op, err = opC.GetOperation(context.Background(), &longrunningpb.GetOperationRequest{
				Name: aop.Name(),
			})
			if err != nil {
				log.Printf("Error polling audio content: %v", err)
				return
			}
			if op.Done {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if op.GetError() != nil {
			log.Printf("Error extracting audio content: %v", err)
		}
		if op.GetResponse() != nil {
			var res speechpb.AsyncRecognizeResponse
			if err := proto.Unmarshal(op.GetResponse().Value, &res); err != nil {
				log.Printf("Error unmarshalling speech to text response: %v", err)
			}

			var buffer bytes.Buffer
			for _, result := range res.Results {
				for _, alt := range result.Alternatives {
					n, err := normalize_text.Normalize(alt.Transcript)
					if err != nil {
						log.Printf("Error normalizing audio content: %v", err)
					}
					_, err = buffer.WriteString(n)
					if err != nil {
						log.Printf("Error concatenating audio content transcripts: %v", err)
					}
					_, err = buffer.WriteString(" ")
					if err != nil {
						log.Printf("Error concatenating audio content transcripts: %v", err)
					}
				}
			}

			sttc.AudioContent = buffer.String()
		}

		wg.Done()
	}()

	wg.Wait()

	return sttc, nil
}

func Content(rc io.ReadCloser) (SpeechToText, error) {
	defer rc.Close()

	// Creates a client.
	client, err := speech.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	// Reads the audio file into memory.
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	// Detects speech in the audio file.
	resp, err := client.SyncRecognize(context.Background(), &speechpb.SyncRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:   speechpb.RecognitionConfig_ENCODING_UNSPECIFIED,
			SampleRate: 16000,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: b},
		},
	})

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil, nil
}
