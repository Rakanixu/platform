package elastic

import (
	"encoding/json"
	"errors"
	"fmt"

	proto "github.com/kazoup/platform/elastic/srv/proto/elastic"
	lib "github.com/mattbaird/elastigo/lib"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("not found")
	// Hosts elasticsearch
	Hosts []string
	conn  *lib.Conn
)

// Init ES connection
func Init() {
	conn = lib.NewConn()
	conn.SetHosts(Hosts)
}

// Create record
func Create(cr *proto.CreateRequest) error {
	_, err := conn.Index(cr.Index, cr.Type, cr.Id, nil, cr.Data)

	return err
}

// Read record
func Read(rr *proto.ReadRequest) (string, error) {
	r, err := conn.Get(rr.Index, rr.Type, rr.Id, nil)
	if err != nil {
		return "", err
	}

	data, _ := r.Source.MarshalJSON()

	return string(data), nil
}

// Update record
func Update(ur *proto.UpdateRequest) error {
	_, err := conn.Index(ur.Index, ur.Type, ur.Id, nil, ur.Data)

	return err
}

// Delete record
func Delete(dr *proto.DeleteRequest) error {
	_, err := conn.Delete(dr.Index, dr.Type, dr.Id, nil)

	return err
}

// Search ES index
func Search(sr *proto.SearchRequest) (string, error) {
	if len(sr.Query) <= 0 {
		sr.Query = "*"
	}
	size := fmt.Sprintf("%d", sr.Limit)
	from := fmt.Sprintf("%d", sr.Offset)

	out, err := lib.Search(sr.Index).Type(sr.Type).Size(size).From(from).Search(sr.Query).Result(conn)
	if err != nil {
		return "", err
	}

	return string(out.RawJSON), nil
}

// Query DSL ES
func Query(sr *proto.QueryRequest) (string, error) {
	result, err := conn.Search(sr.Index, sr.Type, nil, sr.Query)

	if err != nil {
		return "", err
	}

	return string(result.RawJSON), nil
}

// CreateIndexWithSettings creates an ES index with settings
func CreateIndexWithSettings(r *proto.CreateIndexWithSettingsRequest) error {
	var settingsMap map[string]interface{}

	if err := json.Unmarshal([]byte(r.Settings), &settingsMap); err != nil {
		return err
	}
	fmt.Println(r.Settings)
	fmt.Println(settingsMap)

	_, err := conn.CreateIndexWithSettings(r.Index, settingsMap)

	if err != nil {
		return err
	}

	return nil
}

// PutMappingFromJSON puts a mapping into ES
func PutMappingFromJSON(r *proto.PutMappingFromJSONRequest) error {
	if _, err := conn.CloseIndex(r.Index); err != nil {
		return err
	}

	if err := conn.PutMappingFromJSON(r.Index, r.Type, []byte(r.Mapping)); err != nil {
		return err
	}

	if _, err := conn.OpenIndex(r.Index); err != nil {
		return err
	}

	return nil
}
