package admin

// PackageForm 套餐表单
type PackageForm struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name" validate:"required"`
	ValidDays   int      `json:"valid_days" validate:"required,gte=1"`
	DeviceLimit int      `json:"device_limit" validate:"required,gte=1"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	IsActive    bool     `json:"is_active"`
	Priority    int      `json:"priority"`
	ServerIds   []uint   `json:"server_ids"`
}

// ServerForm 服务器表单
type ServerForm struct {
	Id          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Region      string `json:"region"`
	IdServer    string `json:"id_server" validate:"required"`
	RelayServer string `json:"relay_server" validate:"required"`
	Key         string `json:"key"`
	ApiServer   string `json:"api_server"`
	WsHost      string `json:"ws_host"`
	SupportTCP  bool   `json:"support_tcp"`
	SupportWSS  bool   `json:"support_wss"`
	CostWeight  int    `json:"cost_weight"`
	IsDefault   bool   `json:"is_default"`
	IsActive    bool   `json:"is_active"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
}
