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
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_create_index_test"
		}
	}`), &http.Response{StatusCode: 200}},
}

var config_status = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.Status",
		"request": {}
	}`), &http.Response{StatusCode: 200}},
}

var config_addalias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_add_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_add_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
}

var config_deleteindex = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_index_test"
		}
	}`), &http.Response{StatusCode: 200}},
}

var config_deletealias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteAlias",
		"request": {
			"index": "db_config_srv_delete_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_delete_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
}

var config_renamealias = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.RenameAlias",
		"request": {
			"index": "db_config_srv_rename_alias_test",
			"old_alias": "test_alias",
			"new_alias": "test_alias_renamed"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_rename_alias_test"
		}
	}`), &http.Response{StatusCode: 200}},
}

func TestDBCreateIndex(t *testing.T) {
	// Create index, then delete index
	rangeTestTable(config_createindex, t)
}

func TestDBStatus(t *testing.T) {
	rangeTestTable(config_status, t)
}

func TestDBAddAlias(t *testing.T) {
	// Create index, add alias, delete index
	rangeTestTable(config_addalias, t)
}

func TestDBDeleteIndex(t *testing.T) {
	// Create index, delete index
	rangeTestTable(config_deleteindex, t)
}

func TestDBDeleteAlias(t *testing.T) {
	// Create index, add alias, delete alias, delete index
	rangeTestTable(config_deletealias, t)
}

func TestDBRenameAlias(t *testing.T) {
	// Create index, add alias, rename alias, delete index
	rangeTestTable(config_renamealias, t)
}

/* THIS was used with local files,
func TestDBAggregate(t *testing.T) {

}*/
