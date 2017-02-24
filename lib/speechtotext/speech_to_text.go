package speechtotext

import (
	speech "cloud.google.com/go/speech/apiv1beta1"
	//"azul3d.org/audio.v1"
	"fmt"
	"github.com/youpy/go-wav"
	"golang.org/x/net/context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	aop, err := client.AsyncRecognize(context.Background(), &speechpb.AsyncRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:   speechpb.RecognitionConfig_FLAC,
			SampleRate: 16000,
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
				log.Println("Error polling audio content: %v", err)
			}
			if op.Done {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if op.GetError() != nil {
			log.Println("Error extracting audio content: %v", err)
		}

		if op.GetResponse() != nil {
			sttc.AudioContent = string(op.GetResponse().Value)
		}

		wg.Done()
	}()

	wg.Wait()

	return sttc, nil
}

// ummm
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

	log.Println("----", resp)
	log.Println("ERR", err)

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil, nil
}

func Wav(rc *os.File) {
	defer rc.Close()
	/*	infile_path := flag.String("infile", "", "wav file to read")
		flag.Parse()

		file, _ := os.Open(*infile_path)*/
	reader := wav.NewReader(rc)

	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			fmt.Printf("L/R: %d/%d\n", reader.IntValue(sample, 0), reader.IntValue(sample, 1))
		}
	}

}
