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

import (
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

var delete_file_tests_data = testTable{
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Delete",
		"request":{
			"datasource_id":"f8dd3f4cdecfcd6a8103570a38f9c723",
			"index":"index86e5e9616454-4779-84cd-7e308167f0c2",
			"file_id":"57ed5a61f513b361b18df7dfa42af6de",
			"original_id":"B8F82FA8A78FA9C8!115",
			"original_file_path":""
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Delete",
		"request":{
			"datasource_id":"f8dd3f4cdecfcd6a8103570a38f9c723",
			"index":"index86e5e9616454-4779-84cd-7e308167f0c2",
			"file_id":"57ed5a61f513b361b18df7dfa42af6de",
			"original_id":"B8F82FA8A78FA9C8!115",
			"original_file_path":""
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service":"com.kazoup.srv.file",
		"method":"File.Delete",
		"request":{
			"datasource_id":"f8dd3f4cdecfcd6a8103570a38f9c723",
			"index":"index86e5e9616454-4779-84cd-7e308167f0c2",
			"file_id":"57ed5a61f513b361b18df7dfa42af6de",
			"original_id":"B8F82FA8A78FA9C8!115",
			"original_file_path":""
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

func TestFileCreate(t *testing.T) {
	// Create datasources to perfomr file operations over them
	rangeTestTable(ds_file_tests_data, t)

	time.Sleep(time.Second * 30)

	// Create a file per datasource
	rangeTestTable(create_file_tests_data, t)

	/*

		data
		:
		"{"id":"05ca5f37f3d1fb751a19bcc358058384","user_id":"google-apps|pablo.aguirre@kazoup.com","name":"test","url":"https://docs.google.com/a/kazoup.com/document/d/1ure_-YanCUKHh_AgqdYCup5hVCY-ZUMZxwnsuhOMXHU/edit?usp=drivesdk","modified":"2017-01-27T16:55:11.306Z","file_size":0,"is_dir":false,"category":"None","mime_type":"application/vnd.google-apps.document","depth":0,"file_type":"googledrive","last_seen":1485536115,"access":"private","datasource_id":"5f7393781db95c51ad03cdf23a42dd1f","index":"index22e2f85a0bbc-4b44-a4b0-6a043f86c8e6","original":{"capabilities":{"canComment":true,"canCopy":true,"canEdit":true,"canReadRevisions":true,"canShare":true},"createdTime":"2017-01-27T16:55:11.306Z","iconLink":"https://ssl.gstatic.com/docs/doclist/images/icon_11_document_list.png","id":"1ure_-YanCUKHh_AgqdYCup5hVCY-ZUMZxwnsuhOMXHU","isAppAuthorized":true,"kind":"drive#file","lastModifyingUser":{"displayName":"Pablo Aguirre","emailAddress":"pablo.aguirre@kazoup.com","kind":"drive#user","me":true,"permissionId":"09634826227332287579"},"mimeType":"application/vnd.google-apps.document","modifiedByMeTime":"2017-01-27T16:55:11.306Z","modifiedTime":"2017-01-27T16:55:11.306Z","name":"test","ownedByMe":true,"owners":[{"displayName":"Pablo Aguirre","emailAddress":"pablo.aguirre@kazoup.com","kind":"drive#user","me":true,"permissionId":"09634826227332287579"}],"parents":["0AI9SNXL1FJ8pUk9PVA"],"permissions":[{"displayName":"Pablo Aguirre","emailAddress":"pablo.aguirre@kazoup.com","id":"09634826227332287579","kind":"drive#permission","role":"owner","type":"user"}],"spaces":["drive"],"version":"5366","viewedByMe":true,"viewedByMeTime":"2017-01-27T16:55:11.461Z","viewersCanCopyContent":true,"webViewLink":"https://docs.google.com/a/kazoup.com/document/d/1ure_-YanCUKHh_AgqdYCup5hVCY-ZUMZxwnsuhOMXHU/edit?usp=drivesdk","writersCanShare":true}}"
		doc_url
		:
		"https://docs.google.com/a/kazoup.com/document/d/1ure_-YanCUKHh_AgqdYCup5hVCY-ZUMZxwnsuhOMXHU/edit?usp=drivesdk"
	*/

}

// Tear down of TestFileCreate
/*
func TestFileDelete(t *testing.T) {

}
*/

/*Gdrive
{
"service":"com.kazoup.srv.file",
"method":"File.Delete",
"request":{
"datasource_id":"5f7393781db95c51ad03cdf23a42dd1f",
"index":"indexc1374dc3ffcb-4755-b265-5d69d786dc33",
"file_id":"328859ee95ab34fca6a17dbbb3d1c4e5",
"original_id":"1rJYyqzjaWo5pFpcZlzfuRgec_IMisQHz7vdy7eETdoc",
"original_file_path":""
}
}*/

/*
//Dropbox
{
   "service":"com.kazoup.srv.file",
   "method":"File.Delete",
   "request":{
      "datasource_id":"e80d54ad29d18cb62cf9bb2bb54fcfd5",
      "index":"indexcad4e1f4c1fd-42d3-ad92-b07ecbcdf5a8",
      "file_id":"f35de774285d862cfcc0c0053955c3ef",
      "original_id":"id:lXWZMx78s2AAAAAAAAAABg",
      "original_file_path":"/test.docx"
   }
}
*/
