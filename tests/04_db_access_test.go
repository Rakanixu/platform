package tests

import (
	"github.com/kazoup/platform/lib/globals"
	"net/http"
	"testing"
	"time"
)

var db_create_forbidden_access = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_create_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_create_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_create_forbidden_access_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_create_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_create_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var db_read_forbidden_access = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Read",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_read_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_read_forbidden_access_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Read",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_read_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var db_update_forbidden_access = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Update",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"updated string\",\"bool\": false,\"int\": 0}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_update_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_update_forbidden_access_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Update",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"updated string\",\"bool\": false,\"int\": 0}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_update_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var db_delete_forbidden_access = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Delete",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_delete_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_delete_forbidden_access_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Delete",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_delete_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var db_delete_forbidden_accessbyquery = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_deletebyquery_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.DeleteByQuery",
		"request": {
			"indexes": ["db_srv_deletebyquery_test"],
			"types": ["test_document"],
			"last_seen": 2
		}
	}`), &http.Response{StatusCode: 403}, time.Second},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_deletebyquery_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_delete_forbidden_accessbyquery_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_deletebyquery_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.DeleteByQuery",
		"request": {
			"indexes": ["db_srv_deletebyquery_test"],
			"types": ["test_document"],
			"last_seen": 2
		}
	}`), &http.Response{StatusCode: 401}, time.Second},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_deletebyquery_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var db_create_forbidden_access_for_search = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_db_srv_search_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"user_id\": \"` + USER_ID + `\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_db_srv_search_test",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"name\": \"orange\", \"user_id\": \"` + USER_ID + `\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
}

var db_search_forbidden_access = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Search",
		"request": {
			"index": "index_db_srv_search_test",
			"type": "test_document",
			"term": "tree",
			"from": 0,
			"size": 2
		}
	}`), &http.Response{StatusCode: 403}, time.Second},
}

var db_search_forbidden_access_systemToken = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Search",
		"request": {
			"index": "db_srv_search_test",
			"type": "test_document",
			"term": "tree"
		}
	}`), &http.Response{StatusCode: 403}, time.Second},
}

var db_delete_forbidden_accessindex_for_search = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "index_db_srv_search_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

var db_create_forbidden_access_for_searchbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"id\": \"test_id_1\", \"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"id\": \"test_id_2\", \"name\": \"orange\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 403}, noDuration},
}

var db_search_forbidden_accessbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.SearchById",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_1",
			"user_id": "` + USER_NAME + `",
			"name": "tree"
		}
	}`), &http.Response{StatusCode: 403}, time.Second},
}

var db_search_forbidden_accessbyid_systemToken = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.SearchById",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_1",
			"name": "tree"
		}
	}`), &http.Response{StatusCode: 403}, time.Second},
}

var db_delete_forbidden_accessindex_for_searchbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_searchbyid_test"
		}
	}`), &http.Response{StatusCode: 500}, noDuration},
}

func TestDBCreateForbiddenAccess(t *testing.T) {
	// Create document, delete index
	rangeTestTable(db_create_forbidden_access, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBCreateInvalidJWTForbiddenAccess(t *testing.T) {
	// Create document, delete index
	rangeTestTable(db_create_forbidden_access_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBReadForbiddenAccess(t *testing.T) {
	// Create document, read document, delete index
	rangeTestTable(db_read_forbidden_access, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBReadInvalidJWTForbiddenAccess(t *testing.T) {
	// Create document, read document, delete index
	rangeTestTable(db_read_forbidden_access_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBUpdateForbiddenAccess(t *testing.T) {
	// Create document, update document, delete index
	rangeTestTable(db_update_forbidden_access, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBUpdateInvalidJWTForbiddenAccess(t *testing.T) {
	// Create document, update document, delete index
	rangeTestTable(db_update_forbidden_access_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBDeleteForbiddenAccess(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_delete_forbidden_access, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBDeleteInvalidJWTForbiddenAccess(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_delete_forbidden_access_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBDeleteByQueryForbiddenAccess(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_delete_forbidden_accessbyquery, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBDeleteByQueryInvalidJWTForbiddenAccess(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_delete_forbidden_accessbyquery_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBSearchForbiddenAccess(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_forbidden_access_for_search, JWT_TOKEN_USER_1, emptyHeader, t)

	// Search by term for 1 document
	rangeTestTable(db_search_forbidden_access, JWT_TOKEN_USER_1, emptyHeader, t)

	// Delete index used for test
	rangeTestTable(db_delete_forbidden_accessindex_for_search, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBSearchSystemTokenForbiddenAccess(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_forbidden_access_for_search, globals.SYSTEM_TOKEN, emptyHeader, t)

	rangeTestTable(db_search_forbidden_access_systemToken, globals.SYSTEM_TOKEN, emptyHeader, t)

	// Delete index used for test
	rangeTestTable(db_delete_forbidden_accessindex_for_search, globals.SYSTEM_TOKEN, emptyHeader, t)
}

func TestDBSearchByIdForbiddenAccess(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_forbidden_access_for_searchbyid, JWT_TOKEN_USER_1, emptyHeader, t)

	// Search by id for 1 document
	rangeTestTable(db_search_forbidden_accessbyid, JWT_TOKEN_USER_1, emptyHeader, t)

	// Delete index used for test
	rangeTestTable(db_delete_forbidden_accessindex_for_searchbyid, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBSearchByIdSystemTokenForbiddenAccess(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_forbidden_access_for_searchbyid, globals.SYSTEM_TOKEN, emptyHeader, t)

	rangeTestTable(db_search_forbidden_accessbyid_systemToken, globals.SYSTEM_TOKEN, emptyHeader, t)

	// Delete index used for test
	rangeTestTable(db_delete_forbidden_accessindex_for_searchbyid, globals.SYSTEM_TOKEN, emptyHeader, t)
}
