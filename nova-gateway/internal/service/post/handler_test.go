package post

import (
	"context"
	"reflect"
	"testing"

	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"google.golang.org/grpc"
)

type fakeAuthorUserClient struct {
	gotIDs []int64
}

func (f *fakeAuthorUserClient) BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	f.gotIDs = append([]int64(nil), in.GetIds()...)
	return &userv1.BatchGetUsersResponse{
		Users: []*userv1.User{
			{Id: 1, Username: "alice", Nickname: "Alice Nick"},
			{Id: 2, Username: "bob"},
		},
	}, nil
}

func TestEnrichPostAuthorsBatchGetsUsers(t *testing.T) {
	userClient := &fakeAuthorUserClient{}
	module := &Module{userClient: userClient}
	posts := []*postv1.Post{
		{Id: 1, AuthorId: 1},
		{Id: 2, AuthorId: 2},
		{Id: 3, AuthorId: 1},
		{Id: 4},
		nil,
	}

	if err := module.enrichPostAuthors(context.Background(), posts...); err != nil {
		t.Fatal(err)
	}

	if want := []int64{1, 2}; !reflect.DeepEqual(userClient.gotIDs, want) {
		t.Fatalf("batch ids = %v, want %v", userClient.gotIDs, want)
	}
	if posts[0].AuthorName != "Alice Nick" {
		t.Fatalf("author name = %q, want Alice Nick", posts[0].AuthorName)
	}
	if posts[1].AuthorName != "bob" {
		t.Fatalf("author name = %q, want bob", posts[1].AuthorName)
	}
	if posts[2].AuthorName != "Alice Nick" {
		t.Fatalf("author name = %q, want Alice Nick", posts[2].AuthorName)
	}
}
