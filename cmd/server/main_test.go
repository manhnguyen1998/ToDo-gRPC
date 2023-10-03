package main

import (
	"sync"
	"testing"

	"connectrpc.com/connect"
	todov1 "example.com/todo/gen/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRead(t *testing.T) {
	cases := []struct {
		name    string
		id      string
		want    *todov1.ToDo
		wantErr error
	}{
		{
			name: "exit id",
			id:   "exit-id",
			want: &todov1.ToDo{
				Id:     "exit-id",
				Name:   "exit-name",
				Status: todov1.Status_STATUS_DONE,
			},
			wantErr: nil,
		},
		{
			name:    "not found id",
			id:      "not-found-id",
			want:    nil,
			wantErr: status.Error(codes.NotFound, "not-found-id is not found"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			todoServer := &ToDoServer{
				m: sync.Map{},
			}
			todoServer.m.Store("exit-id", &todov1.ToDo{
				Id:     "exit-id",
				Name:   "exit-name",
				Status: todov1.Status_STATUS_DONE,
			})
			response, err := todoServer.Read(nil, &connect.Request[todov1.ReadRequest]{
				Msg: &todov1.ReadRequest{
					Id: c.id,
				},
			})
			if err != nil && err.Error() != c.wantErr.Error() {
				t.Errorf("got %v, want %v", err, c.wantErr)
			}
			if err == nil {
				if response.Msg.Todo.Id != c.want.Id {
					t.Errorf("got %v, want %v", response.Msg.Todo.Id, c.want.Id)
				}
				if response.Msg.Todo.Name != c.want.Name {
					t.Errorf("got %v, want %v", response.Msg.Todo.Name, c.want.Name)
				}
				if response.Msg.Todo.Status != c.want.Status {
					t.Errorf("got %v, want %v", response.Msg.Todo.Status, c.want.Status)
				}
			}
		})
	}
}
