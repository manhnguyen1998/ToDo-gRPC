package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	todov1 "example.com/todo/gen/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRead(t *testing.T) {
	server := &ToDoServer{}

	// テストデータのセットアップ
	// ここでは、テスト用のToDoを作成し、それをsync.Mapに保存します。
	testID := "test-id"
	testTodo := &todov1.ToDo{
		Id:     testID,
		Name:   "Test Todo",
		Status: todov1.Status_STATUS_IMPORTANT,
	}
	m.Store(testID, testTodo)

	// 既存のToDoに対するテストケース
	t.Run("Existing ToDo", func(t *testing.T) {
		// テスト用のリクエストを作成
		req := &connect.Request[todov1.ReadRequest]{
			Msg: &todov1.ReadRequest{Id: testID},
		}
		// Read関数を呼び出し
		resp, err := server.Read(context.Background(), req)
		// エラーが返されないことを確認
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// 返されたToDoのIDが期待されるものであることを確認
		if resp.Msg.Todo.Id != testID {
			t.Errorf("Expected ID %v, got %v", testID, resp.Msg.Todo.Id)
		}
		// 返されたToDoのステータスが期待されるものであることを確認
		if resp.Msg.Todo.Status != todov1.Status_STATUS_IMPORTANT {
			t.Errorf("Expected Status %v, got %v", todov1.Status_STATUS_IMPORTANT, resp.Msg.Todo.Status)
		}
	})

	// 存在しないToDoに対するテストケース
	t.Run("Non-existing ToDo", func(t *testing.T) {
		nonExistingID := "non-existing-id"
		req := &connect.Request[todov1.ReadRequest]{
			Msg: &todov1.ReadRequest{Id: nonExistingID},
		}
		_, err := server.Read(context.Background(), req)
		// エラーが返されることを確認
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		// エラーがgRPCのステータスエラーであることを確認
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("Expected gRPC status error, got %T", err)
		}
		// エラーコードがNotFoundであることを確認
		if st.Code() != codes.NotFound {
			t.Errorf("Expected error code %v, got %v", codes.NotFound, st.Code())
		}
	})
}
