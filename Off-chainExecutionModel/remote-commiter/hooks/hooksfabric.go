package hooks

import (
	"github.com/vadiminshakov/committer/hooks/src"
	pb "github.com/vadiminshakov/committer/proto"
	"github.com/vadiminshakov/committer/server"
)

type ProposeHook func(req *pb.ProposeRequest) bool
type CommitHook func(req *pb.CommitRequest) bool

func Get() ([]server.Option, error) {
	proposeHook := func(f ProposeHook) func(*server.Server) error {
		return func(server *server.Server) error {
			server.ProposeHook = f
			return nil
		}
	}
	commitHook := func(f CommitHook) func(*server.Server) error {
		return func(server *server.Server) error {
			server.CommitHook = f
			return nil
		}
	}
	//把src.Propose src.Commit 两个函数提取，方便绑定到server的ProposeHook 和 CommitHook
	//把传入的src.Propose作为
	return []server.Option{proposeHook(src.Propose), commitHook(src.Commit)}, nil
}
