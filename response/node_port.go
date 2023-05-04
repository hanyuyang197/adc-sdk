package response

type Node_port struct {
	Node_id          int    `json:"node_id" form:"node_id" gorm:"column:node_id;comment:;type:bigint;"`
	Node_name        string `json:"node_name" form:"node_name" gorm:"column:node_name;comment:"`
	Port_number      int    `json:"port_number" form:"port_number" gorm:"column:port_number;comment:"`
	Protocol         int    `json:"protocol" form:"protocol" gorm:"column:protocol;comment:"`
	Status           int    `json:"status" form:"status" gorm:"column:status;comment:"`
	Weight           int    `json:"weight" form:"weight" gorm:"column:weight;comment:"`
	Conn_limit       int    `json:"conn_limit" form:"conn_limit" gorm:"column:conn_limit;comment:"`
	Graceful_delete  int    `json:"graceful_delete" form:"graceful_delete" gorm:"column:graceful_delete;comment:"`
	Graceful_disable int    `json:"graceful_disable" form:"graceful_disable" gorm:"column:graceful_disable;comment:"`
	Graceful_persist int    `json:"graceful_persist" form:"graceful_persist" gorm:"column:graceful_persist;comment:"`
	Graceful_time    int    `json:"graceful_time" form:"graceful_time" gorm:"column:graceful_time;comment:"`
	Phm_profile      string `json:"phm_profile" form:"phm_profile" gorm:"column:phm_profile;comment:"`
	Healthcheck      string `json:"healthcheck" form:"healthcheck" gorm:"column:healthcheck;comment:"`
	GroupID          uint   `json:"group_id" form:"group_id" gorm:"column:group_id;comment:设备组"`
	DevID            int    `json:"dev_id" form:"dev_id" gorm:"column:dev_id;comment:设备主键id"`
	Class            string `json:"class" form:"class" gorm:"column:class;default:'NodePort';comment:"` // 兼容模板化下发
}

func (Node_port) TableName() string {
	return "node_port"
}

type NodePortRequest struct {
	Port Node_port `json:"port"`
	Name string    `json:"name" form:"name" gorm:"column:name;comment:节点名称"`
}

type NodePortIdsReq struct {
	Ids       []int  `json:"ids" form:"ids"`
	Node_name string `json:"node_name" form:"node_name"`
}
