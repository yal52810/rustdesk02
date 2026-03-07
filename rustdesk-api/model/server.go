package model

import "time"

// Server RustDesk 服务器线路配置
type Server struct {
	IdModel
	Name        string     `json:"name" gorm:"size:100;not null;comment:线路名称"`
	Region      string     `json:"region" gorm:"size:50;comment:地区标签(CN-East,HK,US等)"`
	IdServer    string     `json:"id_server" gorm:"size:255;not null;comment:ID服务器地址"`
	RelayServer string     `json:"relay_server" gorm:"size:255;not null;comment:中继服务器地址"`
	Key         string     `json:"key" gorm:"size:255;comment:服务器公钥"`
	ApiServer   string     `json:"api_server" gorm:"size:255;comment:API服务器地址"`
	WsHost      string     `json:"ws_host" gorm:"size:255;comment:WebSocket地址(wss://domain:21118)"`
	SupportTCP  bool       `json:"support_tcp" gorm:"default:true;comment:是否支持TCP协议"`
	SupportWSS  bool       `json:"support_wss" gorm:"default:false;comment:是否支持WSS协议"`
	CostWeight  int        `json:"cost_weight" gorm:"default:1;comment:线路成本权重(1-10,越大越贵)"`
	IsDefault   bool       `json:"is_default" gorm:"default:false;comment:是否默认线路"`
	IsActive    bool       `json:"is_active" gorm:"default:true;comment:是否启用"`
	Priority    int        `json:"priority" gorm:"default:0;comment:优先级(越大越优先)"`
	IsOnline    bool       `json:"is_online" gorm:"default:true;comment:是否在线"`
	LastCheckAt *time.Time `json:"last_check_at" gorm:"comment:最后检查时间"`
	Description string     `json:"description" gorm:"type:text;comment:描述"`
	TimeModel
}

// ServerList 服务器列表响应
type ServerList struct {
	Servers []*Server `json:"list"`
	Pagination
}
