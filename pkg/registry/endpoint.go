package registry

import (
	"encoding/json"

	"github.com/douyu/jupiter/pkg/server"
)

// Endpoints ...
type Endpoints struct {
	// 服务节点列表
	Nodes map[string]server.ServiceInfo

	// 路由配置
	RouteConfigs map[string]RouteConfig
}

// RouteConfig ...
type RouteConfig struct {
	ID     string `json:"id" toml:"id"`
	Scheme string `json:"scheme" toml:"scheme"`
	Host   string `json:"host" toml:"host"`

	Deployment string   `json:"deployment"`
	URI        string   `json:"uri"`
	Upstream   Upstream `json:"upstream"`
}

// String ...
func (config RouteConfig) String() string {
	bs, _ := json.Marshal(config)
	return string(bs)
}

// Upstream ...
type Upstream struct {
	Nodes  map[string]int `json:"nodes"`
	Groups map[string]int `json:"groups"`
}
