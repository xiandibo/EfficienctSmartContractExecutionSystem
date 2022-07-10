package addrtoip

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

var ConsensusLeaderNode = NodeList

// LeaderToVoter 通过leader节点地址找到共识组的voter
// 用于共识组内部共识
// 端口分配：300x leader 400x 500x voter
var LeaderToVoter = map[string][]string{
	// Leader : Voters
	// 10.20.3.46
	"10.20.3.46:3001": {"10.20.6.199:4111", "10.20.36.229:4112"},
	"10.20.3.46:3002": {"10.20.0.57:4121", "10.20.6.199:4122"},
	"10.20.3.46:3003": {"10.20.36.229:4131", "10.20.0.57:4132"},
	//"10.20.3.46:3004": {"10.20.3.46:4004", "10.20.3.46:5004"},

	// 10.20.6.199
	"10.20.6.199:3001": {"10.20.3.46:4211", "10.20.36.229:4212"},
	"10.20.6.199:3002": {"10.20.0.57:4221", "10.20.3.46:4222"},
	"10.20.6.199:3003": {"10.20.36.229:4231", "10.20.0.57:4232"},
	//"10.20.6.199:3004": {"10.20.6.199:4004", "10.20.6.199:5004"},

	// 10.20.36.229
	"10.20.36.229:3001": {"10.20.3.46:4311", "10.20.6.199:4312"},
	"10.20.36.229:3002": {"10.20.0.57:4321", "10.20.3.46:4322"},
	"10.20.36.229:3003": {"10.20.6.199:4331", "10.20.0.57:4332"},
	//"10.20.36.229:3004": {"10.20.36.229:4004", "10.20.36.229:5004"},

	// 10.20.0.57
	"10.20.0.57:3001": {"10.20.3.46:4411", "10.20.6.199:4412"},
	"10.20.0.57:3002": {"10.20.36.229:4421", "10.20.3.46:4422"},
	"10.20.0.57:3003": {"10.20.6.199:4431", "10.20.36.229:4432"},
	//"10.20.0.57:3004": {"10.20.0.57:4004", "10.20.0.57:5004"},

}

// 对于每个leader节点，需要初始化其与组内voter的grpc链接
func EstablishConnWithVoter(leader string) error {
	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  5 * time.Second,
		},
		MinConnectTimeout: 1000 * time.Millisecond,
	}
	// 只有leader能够连接其Voter
	s := LeaderToVoter[leader]
	for _, addr := range s {
		conn, err := grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure())
		if err != nil {
			return errors.Wrap(err, "failed to connect")
		} else {
			NodeAddrConnMap[addr] = conn
		}
	}
	return nil
}