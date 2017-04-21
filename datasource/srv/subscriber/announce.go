package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	deletebucket_msg "github.com/kazoup/platform/lib/protomsg/deletebucket"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnDatasourceCreate
func (a *AnnounceHandler) OnDatasourceCreate(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Trigger scan on datasource creation
	if globals.HANDLER_DATASOURCE_CREATE == msg.Handler {
		var r proto_datasource.CreateRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
			globals.DiscoverTopic,
			r.Endpoint,
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceDelete
func (a *AnnounceHandler) OnDatasourceDelete(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Trigger bucket deletion for datasource
	if globals.HANDLER_DATASOURCE_DELETE == msg.Handler {
		var r proto_datasource.DeleteRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
			globals.DeleteBucketTopic,
			&deletebucket_msg.DeleteBucketMsg{
				Index: r.Index,
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceScan, trigger scan
func (a *AnnounceHandler) OnDatasourceScan(ctx context.Context, msg *announce.AnnounceMessage) error {
	var e *proto_datasource.Endpoint

	if globals.HANDLER_DATASOURCE_SCAN == msg.Handler {
		var r proto_datasource.ScanRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		rr, err := operations.Read(ctx, &proto_operations.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    r.Id,
		})
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(rr.Result), &e); err != nil {
			return err
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
			globals.DiscoverTopic,
			e,
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceScanAll trigger scans
func (a *AnnounceHandler) OnDatasourceScanAll(ctx context.Context, msg *announce.AnnounceMessage) error {
	var es []*proto_datasource.Endpoint
	var e *proto_datasource.Endpoint

	if globals.HANDLER_DATASOURCE_SCANALL == msg.Handler {
		var r proto_datasource.ScanAllRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if len(r.DatasourcesId) > 0 {
			for _, v := range r.DatasourcesId {
				rr, err := operations.Read(ctx, &proto_operations.ReadRequest{
					Index: globals.IndexDatasources,
					Type:  globals.TypeDatasource,
					Id:    v,
				})
				if err != nil {
					return err
				}

				if err := json.Unmarshal([]byte(rr.Result), &e); err != nil {
					return err
				}

				es = append(es, e)
			}

		} else {
			srvRes, err := operations.Search(ctx, &proto_operations.SearchRequest{
				Index: globals.IndexDatasources,
				Type:  globals.TypeDatasource,
				From:  0,
				Size:  9999,
			})
			if err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(srvRes.Result), &es); err != nil {
				return err
			}
		}

		for _, v := range es {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
				globals.DiscoverTopic,
				v,
			)); err != nil {
				return err
			}
		}
	}

	return nil
}
