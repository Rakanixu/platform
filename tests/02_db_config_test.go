//+build ignore
package tests

import (
	"net/http"
	"testing"
)

var config_createindex = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_create_index_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_create_index_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_createindex_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_create_index_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_create_index_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var config_status = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.Status",
		"request": {}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_status_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.Status",
		"request": {}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var config_addalias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_add_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_addalias_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_add_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var config_deleteindex = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_deleteindex_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var config_deletealias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_deletealias_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

var config_renamealias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.RenameAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"old_alias": "test_alias",
			"new_alias": "test_alias_renamed"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var config_renamealias_invalidJWT = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.RenameAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"old_alias": "test_alias",
			"new_alias": "test_alias_renamed"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 401}, noDuration},
}

func TestDBCreateIndex(t *testing.T) {
	// Create index, then delete index
	rangeTestTable(config_createindex, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBCreateIndexInvalidJWT(t *testing.T) {
	rangeTestTable(config_createindex_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBStatus(t *testing.T) {
	rangeTestTable(config_status, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBStatusInvalidJWT(t *testing.T) {
	rangeTestTable(config_status_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBAddAlias(t *testing.T) {
	// Create index, add alias, delete index
	rangeTestTable(config_addalias, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBAddAliasInvalidJWT(t *testing.T) {
	// Create index, add alias, delete index
	rangeTestTable(config_addalias_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBDeleteIndex(t *testing.T) {
	// Create index, delete index
	rangeTestTable(config_deleteindex, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBDeleteIndexInvalidJWT(t *testing.T) {
	// Create index, delete index
	rangeTestTable(config_deleteindex_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBDeleteAlias(t *testing.T) {
	// Create index, add alias, delete alias, delete index
	rangeTestTable(config_deletealias, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBDeleteAliasInvalidJWT(t *testing.T) {
	// Create index, add alias, delete alias, delete index
	rangeTestTable(config_deletealias_invalidJWT, JWT_INVALID, emptyHeader, t)
}

func TestDBRenameAlias(t *testing.T) {
	// Create index, add alias, rename alias, delete index
	rangeTestTable(config_renamealias, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestDBRenameAliasInvalidJWT(t *testing.T) {
	// Create index, add alias, rename alias, delete index
	rangeTestTable(config_renamealias_invalidJWT, JWT_INVALID, emptyHeader, t)
}
