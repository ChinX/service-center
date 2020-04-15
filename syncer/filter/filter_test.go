package filter

import (
	"encoding/json"
	"testing"

	pb "github.com/apache/servicecomb-service-center/syncer/proto"
)

func TestBlackFilter_Filter(t *testing.T) {
	data := &pb.SyncData{
		Services: []*pb.SyncService{
			{
				App:     "black_app",
				Name:    "black_name",
				Version: "black_version",
			}, {
				App:     "black_app1",
				Name:    "black_name",
				Version: "black_version",
			}, {
				App:     "black_app",
				Name:    "black_name1",
				Version: "black_version",
			}, {
				App:     "black_app2",
				Name:    "black_name2",
				Version: "black_version2",
			},
		},
	}
	//bl := NewWhiteList(
	//	NewMatcher(
	//		WithApp("black_app"),
	//	),
	//	NewMatcher(
	//		WithApp("black_app1"),
	//	),
	//)
	//nd := bl.Filter(data)


	bl := NewBlackList(
		MapToMatcher(map[string]string{"app":"black_app"}),
		MapToMatcher(map[string]string{"app":"black_app1"}),
	)
	nd := bl.Filter(data)

	byts, _ := json.MarshalIndent(nd, "", "  ")
	t.Log(string(byts))
}
