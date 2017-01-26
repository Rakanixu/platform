package tests

import (
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"net/http"
	"testing"
	"time"
)

const (
	NUM_DS_CREATED   = 6
	BOX_URL          = "box://" + USER_ID
	DROPBOX_URL      = "dropbox://" + USER_ID
	GOOGLE_DRIVE_URL = "googledrive://" + USER_ID
	ONE_DRIVE_URL    = "onedrive://" + USER_ID
	GMAIL_URL        = "gmail://" + USER_ID
	SLACK_URL        = "slack://" + USER_ID
)

var datasources []proto.Endpoint

var datasources_test_data = testTable{
	// Box
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + BOX_URL + `",
				"token": {
					"access_token": "fQO1ykhzJU7ig2KQP9wW6NMnzFYAN4Ox",
					"token_type": "bearer",
					"refresh_token": "T0mJ1ywOW1q5CzNZ9gXkNG9iaEiCJItpXBhlRScONPyUk2O7kjfIf5CvSMnCvM9P",
					"expiry": 1485186804
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Dropbox
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + DROPBOX_URL + `",
				"token": {
					"access_token": "jEG_xTrcB7AAAAAAAAAACu1VNyeRFSo0IbRWK-OmhOrivvwuXG8fyOVLyOD2SKoz",
					"token_type": "bearer",
					"expiry": -62135596800
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// GoogleDrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + GOOGLE_DRIVE_URL + `",
				"token": {
					"access_token": "ya29.GlvcA4v06BpeFIgNQz7SDRXNa_97fww5mrpuDHEwPojo78dJuRUK4G9tiBPpzcIeF1yXQkHxAE_vsvuCifJwkJTtlXb71OKecqwMWF5lteTK14tMoK1WZXpm9fUV",
					"token_type": "Bearer",
					"refresh_token": "1/lumkR2KI6aoCkiBz7M4TVcA-2HvQkgYlo5q-lcXQbkQ",
					"expiry": 1485188549
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Onedrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + ONE_DRIVE_URL + `",
				"token": {
					"access_token": "EwAAA61DBAAUGCCXc8wU/zFu9QnLdZXy+YnElFkAAfVOBh/VfpTPt4LIErXT2GEsp03M5ySRNIqX/DprHf3VBtrizDMwiRGqwskP06YVPRiN0xefENrceIOrDr2CnD/NXaD72JUErBz0AM6JtakX0OW/Ida4bztVWreX8zx7A0u057ipipajz39aoBzEQymxD1iVi6HRCqAZ6/53iOpx1k1e/1dhwtfwCEL9UfeHXXSAkBl02MLK/vXqKo+F+4AVrfrvmxIjsLESWpLa7ufS+RbAQOMDbuWOTiwgJtDjJQ9ZMi8c1JJFttLHMZR1WfmVXvqAqGOGTjA9jMUnW74enaS1MnxtBMxCXq4Yjlk/dKe+4lWQTam4rUgssAxddCEDZgAACMF/essAiY2A0AFoyyTfE247XNtP0RGC3GS3dYa9AutciTXwedhqx8lG3ImgrEQFcSlIuB4HJr766En2HNvBqyQyUr3N+wfqrUI5abdvKi4+ZDC+0Vfwkjo4q+JHG4QCa9Zkgtl7t5MFHgZ1sJyFLKCyaBAj/5rvd51rEnMteAoH4inZIvt4McbjsFlqwR6ZuLNvDyc2EbSJ88gWq3/r2gRBbUFh0Oz1gVzWnrl8xLzJ7cjh1vXpCfdQG0Ov7OIZyOBEycqGjM2UnhodojXJ3NJw97g/xWqJDCYRlsSwtGtdqagLI9KLbnpP+n0fUM8U9ejBkAA1OOccsB/g4sfAWzQa8ERLjS9NfStqgDUvKlnEsIlnPJHZyU9IayAmCTJFBonGhnJJIIvIM+UtSmvSIzOqznTasjVnHOPLvu8AW+WOMovVPMZ+aiwZeeob4gIvaEd1CkTYgQq3bPRuYGOWUjFBdwjPHmeafBFA9YhFuiyfbFiEQVGTRq09yRuqdtrVDUTISonPK9R+ZvCNX1KDx2VvUNBgMCD4WhK9/qFa0lZC5gaiYi2A4KrG4RfI7nnq/xqt3pCX8cvH8AYx1rs3QYBIPqR7qT8UstfaM3QrhSOkJLcCpvz4gp0bEAUC",
					"token_type": "bearer",
					"refresh_token": "MCQHdmoGE1TF6Mto!qAphoeY2NlidiFH224yF8Q9biu7oXFxYYL4Ej7F9FYY8lkD6RIntOaqUSGlCQVwSYki*8HmJo16tyUh88kvcvlB7Zb8ZCxg2KNg8Rsm!a3816CBa8TQH55a8ctTNX3jKG6Ps7mPEt0lgUAhVzPYbZ*z5g*GSyZn8W5nxUbYb4llskQKhEiQaB0eEMGraFIgGsKV6X!DPdrjAeqck54Uqh7KTdihG3gH2BYCuBciT2QKmfyqzXh8iuzOflU7ZJ0zivkYWcrs2Py1Qyd3*cVF9hqwTVPOM*ntTq9wsGUnjG8B3xqwMFS5YWklfS6P8llyM!SO4Y6GpdyT!OpgkMghVEvHaSfnken5lJiIBu97LqzYAjcSbIgpfRpJg!r9OwVmmnxd!0DQ$",
					"expiry": 1485194954
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Gmail, this is a real account, but similar to slack, it is difficult to have control about the data to be crawl
	// Usually, this index should be empty, but can be some data as I have no control over mails received (And I won't clean it though SMTP calls)
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + GMAIL_URL + `",
				"token": {
					"access_token": "ya29.GlvcA-QVxeaFH7J2CZD_uKblYMymf7depUQtKBISIoX1n1QPoM404Uw1vmImU94jOtaUUqbhne91HwKNV475mgoGCCb3vf7le0OMwg9Tt5bxc3pZxQMXGbyu4Qvd",
					"token_type": "Bearer",
					"refresh_token": "1/9w8YJy35Hv9AXRqjHKkSVVIb9Gfz7BTdoc-lwXabcfZaC9R1orqySkFMKlGXOjc_",
					"expiry": 1485188853
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Slack, adding a real account would crawl too many files from kazoup team, so, the crawling is not really tested for slack...
	// Timeouts would be difficult to manage
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "` + SLACK_URL + `",
				"token": {
					"access_token": "",
					"token_type": "",
					"refresh_token": "",
					"expiry": 0
				}
			}
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Invalid data
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "DOES_NOT_EXISTS://` + USER_ID + `",
				"token": {
					"access_token": "access_token",
					"token_type": "token_type",
					"refresh_token": "refresh_token",
					"expiry": 0
				}
			}
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var search_fro_crawled_files = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Search",
		"request": {
			"index": "` + globals.GetMD5Hash(USER_ID) + `",
			"user_id": "` + USER_ID + `",
			"type": "file",
			"file_type": "files",
			"from": 0,
			"size": 1000
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var delete_datasources_test_data = testTable{
	// Box
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(BOX_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Dropbox
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(DROPBOX_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// GoogleDrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(GOOGLE_DRIVE_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Onedrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(ONE_DRIVE_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Gmail
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(GMAIL_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	// Slack
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Delete",
		"request": {
			"id": "` + globals.GetMD5Hash(SLACK_URL+USER_ID) + `"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

func TestDatasourceCreate(t *testing.T) {
	// Create datasources,
	rangeTestTable(datasources_test_data, t)

	// Crawlers are triggered for created datasources.
	// Wait half a minute to let crawlers do its job. (There are 4 files per test account)
	time.Sleep(time.Second * 30)

	// Check crawlers behavior. Does indexes exists? There is expected number of elements?
	// We could do this assertions with curl request to ES directly (no internal dependencies), on the other hand,
	// We can do it using kazoup platform. That way ensures a better level of integrity of the system.
	rangeTestTableWithChecker(search_fro_crawled_files, func(rsp *http.Response, t *testing.T) {
		type TestRsp struct {
			Result string `json:"result"`
			Info   string `json:"info"`
		}

		var tr TestRsp
		var tl map[string]int

		if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if err := json.Unmarshal([]byte(tr.Info), &tl); err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		if tl["total"] < 12 {
			t.Errorf("Expected 12 result, got %v", tl["total"])
		}
	}, t)

	// Clean datasources created for the test
	rangeTestTable(delete_datasources_test_data, t)
}

/*
func TestDatasourceSearch(t *testing.T) {
	b := []byte(`{
		"service":"com.kazoup.srv.datasource",
		"method":"DataSource.Search",
		"request":{
			"index":"datasources",
			"type":"datasource",
			"from":0,
			"size":9999,
			"user_id": "` + USER_ID + `"
		}
	}`)

	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Error create request %v", err)
	}

	req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Error performing request with body: %s %v", string(b), err)
	}
	defer res.Body.Close()

	type TestRsp struct {
		Result string `json:"result"`
		Info   string `json:"info"`
	}

	var tr TestRsp

	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if err := json.Unmarshal([]byte(tr.Result), &datasources); err != nil {
		t.Errorf("Error unmarshalling result: %v", err)
	}

	if len(datasources) != NUM_DS_CREATED {
		t.Errorf("Expected %v datasources, got %v", NUM_DS_CREATED, len(datasources))
	}
}

func TestDatasourceScan(t *testing.T) {
	for _, v := range datasources {
		b := []byte(`{
			"service": "com.kazoup.srv.datasource",
			"method": "DataSource.Scan",
			"request": {
				"id": "` + v.Id + `"
			}
		}`)

		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(b))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(b), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != STATUS_OK {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v %v", STATUS_OK, string(b), rsp.StatusCode, string(b))
		}
	}
	time.Sleep(time.Second * 25)
}

func TestDatasourceDelete(t *testing.T) {
	for _, v := range datasources {
		b := []byte(`{
			"service": "com.kazoup.srv.datasource",
			"method": "DataSource.Delete",
			"request": {
				"id": "` + v.Id + `"
			}
		}`)

		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(b))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(b), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != STATUS_OK {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v %v", STATUS_OK, string(b), rsp.StatusCode, string(b))
		}
	}
	time.Sleep(time.Second)
}
*/
