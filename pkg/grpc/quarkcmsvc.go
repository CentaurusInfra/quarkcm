/*
Copyright 2022 quarkcm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grpc

import (
	context "context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/klog"

	"github.com/CentaurusInfra/quarkcm/pkg/datastore"
)

var (
	port = flag.Int("port", 51051, "The server port")
)

type server struct {
	UnimplementedQuarkCMServiceServer
}

func StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterQuarkCMServiceServer(s, &server{})
	klog.Infof("grpc server listening at %v", lis.Addr())
	go s.Serve(lis)
}

// To quickly test this function, execute:
// grpcurl -plaintext -import-path pkg/grpc/ -proto quarkcmsvc.proto -d '{"client_name": "a client name"}' [::]:51051 quarkcmsvc.QuarkCMService/TestPing
func (s *server) TestPing(ctx context.Context, in *TestRequestMessage) (*TestResponseMessage, error) {
	inStr, _ := json.Marshal(in)
	klog.Infof("grpc Service called TestPing %s", inStr)
	hostname, _ := os.Hostname()
	return &TestResponseMessage{ServerName: hostname}, nil
}

func (s *server) ListNode(ctx context.Context, in *emptypb.Empty) (*NodeListMessage, error) {
	klog.Info("grpc Service called ListNode")

	nodeObjects := datastore.ListNode()
	length := len(nodeObjects)
	nodeMessages := make([]*NodeMessage, 0, length)
	for i := 0; i < length; i++ {
		nodeObject := nodeObjects[i]
		nodeMessages = append(nodeMessages, &NodeMessage{
			Name:              nodeObject.Name,
			Hostname:          nodeObject.Hostname,
			Ip:                nodeObject.IP,
			CreationTimestamp: nodeObject.CreationTimestamp,
			ResourceVersion:   int32(nodeObject.ResourceVersion),
		})
	}

	return &NodeListMessage{Nodes: nodeMessages}, nil
}

func (s *server) ListPod(ctx context.Context, in *emptypb.Empty) (*PodListMessage, error) {
	klog.Info("grpc Service called ListPod")

	podObjects := datastore.ListPod()
	length := len(podObjects)
	podMessages := make([]*PodMessage, 0, length)
	for i := 0; i < length; i++ {
		podObject := podObjects[i]
		podMessages = append(podMessages, &PodMessage{
			Key:             podObject.Key,
			Ip:              podObject.IP,
			NodeName:        podObject.NodeName,
			ResourceVersion: int32(podObject.ResourceVersion),
		})
	}

	return &PodListMessage{Pods: podMessages}, nil
}
