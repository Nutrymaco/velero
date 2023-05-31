/*
Copyright the Velero contributors.

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

package framework

import (
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework/common"
	proto "github.com/vmware-tanzu/velero/pkg/plugin/generated"
)

// PreRestoreActionGRPCClient implements the PreRestoreAction interface and uses a
// gRPC client to make calls to the plugin server.
type PreRestoreActionGRPCClient struct {
	*common.ClientBase
	grpcClient proto.PreRestoreActionClient
}

func newPreRestoreActionGRPCClient(base *common.ClientBase, clientConn *grpc.ClientConn) interface{} {
	return &PreRestoreActionGRPCClient{
		ClientBase: base,
		grpcClient: proto.NewPreRestoreActionClient(clientConn),
	}
}

func (c *PreRestoreActionGRPCClient) Execute(restore *api.Restore) error {
	restoreJSON, err := json.Marshal(restore)
	if err != nil {
		return errors.WithStack(err)
	}
	req := &proto.PreRestoreActionExecuteRequest{
		Plugin:  c.Plugin,
		Restore: restoreJSON,
	}
	if _, err = c.grpcClient.Execute(context.Background(), req); err != nil {
		return common.FromGRPCError(err)
	}
	return nil
}
