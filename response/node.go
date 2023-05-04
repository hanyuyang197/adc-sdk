package response

// 节点(node)
type Node struct {
	Tc_name          string      `json:"tc_name" form:"tc_name" gorm:"column:tc_name;comment:流量控制模版名称"`
	Graceful_delete  int         `json:"graceful_delete" form:"graceful_delete" gorm:"column:graceful_delete;comment:删除节点触发软关机"`
	Graceful_disable int         `json:"graceful_disable" form:"graceful_disable" gorm:"column:graceful_disable;comment:禁用节点触发软关机"`
	Graceful_persist int         `json:"graceful_persist" form:"graceful_persist" gorm:"column:graceful_persist;comment:禁用节点触发软关机后会话保持表有效"`
	Graceful_time    int         `json:"graceful_time" form:"graceful_time" gorm:"column:graceful_time;comment:软关机超时，单位秒"`
	Name             string      `json:"name" form:"name" gorm:"column:name;comment:节点名称"`
	Host             string      `json:"host" form:"host" gorm:"column:host;comment:节点主机名或IP地址"`
	Weight           int         `json:"weight" form:"weight" gorm:"column:weight;comment:节点权重"`
	Healthcheck      string      `json:"healthcheck" form:"healthcheck" gorm:"column:healthcheck;comment:节点关联的健康检查名"`
	Status           int         `json:"status" form:"status" gorm:"column:status;comment:节点使能状态"`
	Conn_limit       int         `json:"conn_limit" form:"conn_limit" gorm:"column:conn_limit;comment:节点连接限制"`
	Desc_rserver     string      `json:"desc_rserver" form:"desc_rserver" gorm:"column:desc_rserver;comment:描述"`
	Desc             string      `json:"desc" form:"desc" gorm:"column:desc;comment:描述"`
	GroupID          uint        `json:"group_id" form:"group_id" gorm:"column:group_id;comment:设备组"`
	DevID            int         `json:"dev_id" form:"dev_id" gorm:"column:dev_id;comment:设备主键id"`
	Ports            []Node_port `json:"ports" form:"ports" gorm:"-"`                                                        // 兼容模板化下发
	Class            string      `json:"class" form:"class" gorm:"column:class;default:'Node';comment:"`                     // 兼容模板化下发
	DeployStatus     int         `json:"deploy_status" form:"deploy_status" gorm:"column:deploy_status;default:15;comment:"` // 部署状态详细见Manager.go
	RespResult
}

func (Node) TableName() string {
	return "node"
}

type NodeRequest struct {
	Node interface{} `json:"node"`
	Name string      `json:"name" form:"name" gorm:"column:name;comment:节点名称"`
}

type NodeIdsReq struct {
	Ids   []int    `json:"ids" form:"ids"`
	Names []string `json:"names" form:"names"`
}
