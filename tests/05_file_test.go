package tests

/*
Basic end-to-end tests for file-srv
Those test DOES NOT ensure a good behavior of the service, but provides some basic checks.
Authorization (oath2), internal interaction of kazoup services, designed asynchronicity,
plus real tests data usage adds complexity to automate those tests.

There is no implementation for slack or gmail, as does not make sense.
box auth is expired and won't refresh itself as after refresh first time we do not store test data. It won't work.

Some basic checks for GoogleDrive, OneDrive and Dropbox.
Gdrive creates new file per request, even if the type and name is the same.
Dropbox and onedrive will overwrite the file.
*/

/*import (
	"encoding/json"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"net/http"
	"testing"
	"time"
)

var ds_file_tests_data = testTable{
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
}

var create_file_tests_data = testTable{
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Create",
		"request":{
			"datasource_id":"` + globals.GetMD5Hash(DROPBOX_URL+USER_ID) + `",
			"mime_type":"document",
			"file_name":"test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Create",
		"request":{
			"datasource_id":"` + globals.GetMD5Hash(GOOGLE_DRIVE_URL+USER_ID) + `",
			"mime_type":"document",
			"file_name":"test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Create",
		"request":{
			"datasource_id":"` + globals.GetMD5Hash(ONE_DRIVE_URL+USER_ID) + `",
			"mime_type":"document",
			"file_name":"test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var ds_delete_tests_data = testTable{
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
}

func TestFileCreate(t *testing.T) {
	// Create datasources to perfomr file operations over them
	rangeTestTable(ds_file_tests_data, t)

	time.Sleep(time.Second * 30)

	// This loop is due to unmarshalling to known types
	for k, v := range create_file_tests_data {
		// Create a file per datasource
		rangeTestTableWithChecker(testTable{v}, func(rsp *http.Response, t *testing.T) {
			type TestRsp struct {
				Data   string `json:"data"`
				DocUrl string `json:"doc_url"`
			}

			var tr TestRsp
			var f file.File

			switch k {
			case 0:
				//Dropbox
				f = &file.KazoupDropboxFile{}
			case 1:
				// Google
				f = &file.KazoupGoogleFile{}
			case 2:
				//Onedrive
				f = &file.KazoupOneDriveFile{}
			}

			if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
				t.Fatalf("Error decoding response: %v", err)
			}

			if err := json.Unmarshal([]byte(tr.Data), f); err != nil {
				t.Fatalf("Error unmarshalling response data: %v", err)
			}

			b := []byte(`{
				"service":"com.kazoup.srv.file",
				"method":"File.Delete",
				"request":{
					"datasource_id": "` + f.GetDatasourceID() + `",
					"index": "` + f.GetIndex() + `",
					"file_id": "` + f.GetID() + `",
					"original_id": "` + f.GetIDFromOriginal() + `",
					"original_file_path": "` + f.GetPathDisplay() + `"
				}
			}`)

			// Now we test file deletion
			makeRequest(b, &http.Response{StatusCode: 200}, t)
		}, t)
	}

	// Remove all datrasources created for the test
	rangeTestTable(ds_delete_tests_data, t)
}*/

// Tear down of TestFileCreate
/*
func TestFileDelete(t *testing.T) {

}
*/
