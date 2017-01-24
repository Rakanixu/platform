package tests

import (
	"bytes"
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	RPC_ENPOINT    = "https://web.kazoup.io:8082/rpc"
	USER_ID        = "test@kazoup.com"
	USER_PWD       = "ksu4awemtest"
	STATUS_OK      = 200
	NUM_DS_CREATED = 6
)

var c *http.Client
var datasources []proto.Endpoint

func init() {
	c = &http.Client{}
}

var createtests = []struct {
	in  []byte
	out *http.Response
}{
	// Box
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "box://` + USER_ID + `",
				"token": {
					"access_token": "fQO1ykhzJU7ig2KQP9wW6NMnzFYAN4Ox",
					"token_type": "bearer",
					"refresh_token": "T0mJ1ywOW1q5CzNZ9gXkNG9iaEiCJItpXBhlRScONPyUk2O7kjfIf5CvSMnCvM9P",
					"expiry": 1485186804
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
	// Dropbox
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "dropbox://` + USER_ID + `",
				"token": {
					"access_token": "jEG_xTrcB7AAAAAAAAAACu1VNyeRFSo0IbRWK-OmhOrivvwuXG8fyOVLyOD2SKoz",
					"token_type": "bearer",
					"expiry": -62135596800
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
	// GoogleDrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "googledrive://` + USER_ID + `",
				"token": {
					"access_token": "ya29.GlvcA4v06BpeFIgNQz7SDRXNa_97fww5mrpuDHEwPojo78dJuRUK4G9tiBPpzcIeF1yXQkHxAE_vsvuCifJwkJTtlXb71OKecqwMWF5lteTK14tMoK1WZXpm9fUV",
					"token_type": "Bearer",
					"refresh_token": "1/lumkR2KI6aoCkiBz7M4TVcA-2HvQkgYlo5q-lcXQbkQ",
					"expiry": 1485188549
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
	// Onedrive
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "onedrive://` + USER_ID + `",
				"token": {
					"access_token": "EwAAA61DBAAUGCCXc8wU/zFu9QnLdZXy+YnElFkAAfVOBh/VfpTPt4LIErXT2GEsp03M5ySRNIqX/DprHf3VBtrizDMwiRGqwskP06YVPRiN0xefENrceIOrDr2CnD/NXaD72JUErBz0AM6JtakX0OW/Ida4bztVWreX8zx7A0u057ipipajz39aoBzEQymxD1iVi6HRCqAZ6/53iOpx1k1e/1dhwtfwCEL9UfeHXXSAkBl02MLK/vXqKo+F+4AVrfrvmxIjsLESWpLa7ufS+RbAQOMDbuWOTiwgJtDjJQ9ZMi8c1JJFttLHMZR1WfmVXvqAqGOGTjA9jMUnW74enaS1MnxtBMxCXq4Yjlk/dKe+4lWQTam4rUgssAxddCEDZgAACMF/essAiY2A0AFoyyTfE247XNtP0RGC3GS3dYa9AutciTXwedhqx8lG3ImgrEQFcSlIuB4HJr766En2HNvBqyQyUr3N+wfqrUI5abdvKi4+ZDC+0Vfwkjo4q+JHG4QCa9Zkgtl7t5MFHgZ1sJyFLKCyaBAj/5rvd51rEnMteAoH4inZIvt4McbjsFlqwR6ZuLNvDyc2EbSJ88gWq3/r2gRBbUFh0Oz1gVzWnrl8xLzJ7cjh1vXpCfdQG0Ov7OIZyOBEycqGjM2UnhodojXJ3NJw97g/xWqJDCYRlsSwtGtdqagLI9KLbnpP+n0fUM8U9ejBkAA1OOccsB/g4sfAWzQa8ERLjS9NfStqgDUvKlnEsIlnPJHZyU9IayAmCTJFBonGhnJJIIvIM+UtSmvSIzOqznTasjVnHOPLvu8AW+WOMovVPMZ+aiwZeeob4gIvaEd1CkTYgQq3bPRuYGOWUjFBdwjPHmeafBFA9YhFuiyfbFiEQVGTRq09yRuqdtrVDUTISonPK9R+ZvCNX1KDx2VvUNBgMCD4WhK9/qFa0lZC5gaiYi2A4KrG4RfI7nnq/xqt3pCX8cvH8AYx1rs3QYBIPqR7qT8UstfaM3QrhSOkJLcCpvz4gp0bEAUC",
					"token_type": "bearer",
					"refresh_token": "MCQHdmoGE1TF6Mto!qAphoeY2NlidiFH224yF8Q9biu7oXFxYYL4Ej7F9FYY8lkD6RIntOaqUSGlCQVwSYki*8HmJo16tyUh88kvcvlB7Zb8ZCxg2KNg8Rsm!a3816CBa8TQH55a8ctTNX3jKG6Ps7mPEt0lgUAhVzPYbZ*z5g*GSyZn8W5nxUbYb4llskQKhEiQaB0eEMGraFIgGsKV6X!DPdrjAeqck54Uqh7KTdihG3gH2BYCuBciT2QKmfyqzXh8iuzOflU7ZJ0zivkYWcrs2Py1Qyd3*cVF9hqwTVPOM*ntTq9wsGUnjG8B3xqwMFS5YWklfS6P8llyM!SO4Y6GpdyT!OpgkMghVEvHaSfnken5lJiIBu97LqzYAjcSbIgpfRpJg!r9OwVmmnxd!0DQ$",
					"expiry": 1485194954
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
	// Gmail
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "gmail://` + USER_ID + `",
				"token": {
					"access_token": "ya29.GlvcA-QVxeaFH7J2CZD_uKblYMymf7depUQtKBISIoX1n1QPoM404Uw1vmImU94jOtaUUqbhne91HwKNV475mgoGCCb3vf7le0OMwg9Tt5bxc3pZxQMXGbyu4Qvd",
					"token_type": "Bearer",
					"refresh_token": "1/9w8YJy35Hv9AXRqjHKkSVVIb9Gfz7BTdoc-lwXabcfZaC9R1orqySkFMKlGXOjc_",
					"expiry": 1485188853
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
	// Slack
	{[]byte(`{
		"service": "com.kazoup.srv.datasource",
		"method": "DataSource.Create",
		"request": {
			"endpoint": {
				"user_id": "` + USER_ID + `",
				"url": "slack://` + USER_ID + `",
				"token": {
					"access_token": "",
					"token_type": "",
					"refresh_token": "",
					"expiry": 0
				}
			}
		}
	}`), &http.Response{StatusCode: 200}},
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
	}`), &http.Response{StatusCode: 500}},
}

func TestDatasourceCreate(t *testing.T) {
	for _, v := range createtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}

		if rsp.StatusCode != v.out.StatusCode {
			defer rsp.Body.Close()
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}

	time.Sleep(time.Second)
}

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
