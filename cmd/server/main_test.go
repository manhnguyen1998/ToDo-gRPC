package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"connectrpc.com/connect"
	todov1 "example.com/todo/gen/todo/v1"
	"example.com/todo/gen/todo/v1/todov1connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToDoServer_Read(t *testing.T) {
	t.Parallel()
	todoServer := &ToDoServer{m: sync.Map{},}
	mux := http.NewServeMux()
	mux.Handle(todov1connect.NewToDoServiceHandler(todoServer))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := todov1connect.NewToDoServiceClient(
		server.Client(),
		server.URL,
	)
	grpcClient := todov1connect.NewToDoServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)
	clients := []todov1connect.ToDoServiceClient{connectClient, grpcClient}

	type args struct {
		ctx context.Context
		req *connect.Request[todov1.ReadRequest]
	}
	tests := []struct {
		name    string
		args    args
		want    *connect.Response[todov1.ReadResponse]
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Read with non-exist id",
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&todov1.ReadRequest{
					Id: "Hello",
				}),
			},
			want: nil,
			wantErr: status.Error(codes.Unknown, "Hello is not found"),
		},
		{
			name: "Read with exist id",
			args: args{
				ctx: context.Background(),
				req: connect.NewRequest(&todov1.ReadRequest{
					Id: "exist-id",
				}),
			},
			want: connect.NewResponse(&todov1.ReadResponse{
				Todo: &todov1.ToDo{
					Id: "exist-id",
					Name: "exist-name",
					Status: todov1.Status_STATUS_DONE,
				},

			}),
			wantErr: nil,
		},
	}

	todoServer.m.Store("exist-id", &todov1.ToDo{
		Id:     "exist-id",
		Name:   "exist-name",
		Status: todov1.Status_STATUS_DONE,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clients[0].Read(tt.args.ctx, tt.args.req)
			if err != nil {
				if status.Code(err) != status.Code(tt.wantErr) {
					t.Errorf("ToDoServer.Read() error = %v, wantErr %v", status.Code(err), status.Code(tt.wantErr))
				}
				return
			} else {
				if !reflect.DeepEqual(got.Msg.Todo, tt.want.Msg.Todo) {
				t.Errorf("ToDoServer.Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func NewToDoServer(i int) {
	panic("unimplemented")
}
