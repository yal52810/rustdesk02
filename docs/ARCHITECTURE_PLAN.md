# 全自动二维竞速 + 主控端独裁制 + 三种场景模式

## 用户侧：三种场景模式（不暴露协议术语）

```
设置 → 连接模式：

┌─────────────────────────────────────────────────────┐
│ ● 自动择优 (默认推荐)                                │
│   智能选择最优线路和协议，适合绝大多数用户             │
│                                                     │
│ ○ 专业版 (极致低延迟)                                │
│   强制直连线路，适合局域网或确定无限制的网络            │
│   点开可选择 TCP 中继节点                            │
│                                                     │
│ ○ 企业校园版 (深度穿透)                              │
│   强制 WebSocket 穿透，适合极端受限的公司/校园网络     │
│   点开可选择 WS 中继节点                             │
└─────────────────────────────────────────────────────┘
```

**三条铁律：**
- **UDP/IPv6 P2P 始终开启，不可禁用** — 永远优先直连
- **用户只选"场景模式"，不选协议** — 底层 TCP/WSS 对用户完全透明
- **手动指定节点隐藏在"高级选项"中** — 99% 用户不需要，仅排障兜底

---

## 连接策略总览

```
主控端点击 "连接"
      │
      ├─ ① P2P 直连 (UDP打洞 / IPv6)  ← 始终优先,不可禁用,零成本
      │    失败/超时(3s)
      │
      └─ ② 根据用户选的"场景模式"执行不同策略：
            │
            ├─ 自动择优 (Concurrent Racing):
            │     N个中继 × 2种协议 同时发起
            │     TCP+ChaCha20 (500ms) ∥ WSS+TLS (800ms)
            │     首个完成握手者胜出 → (节点X, 协议Y)
            │
            ├─ 专业版 (TCP Only, Quick Fallback):
            │     N个中继 × 仅 TCP+ChaCha20
            │     300ms 超时 → 失败立即尝试下一个节点
            │
            └─ 企业校园版 (WSS Only, Quick Fallback):
                  N个中继 × 仅 WSS
                  300ms 超时 → 失败立即尝试下一个节点
            │
            └─ ③ 主控端发号施令
                  PunchHoleRequest {
                    relay_server: "节点X:port",
                    relay_use_wss: true/false
                  }
                  │
                  └─ hbbs → PunchHole → 被控端
                        │
                        └─ ④ 被控端无条件跟随
                              连接到节点X 的 协议Y
```

---

## 架构选型：方案二 — 客户端调度（主控端独裁制）

- **hbbs (信令中心)**：单节点集中部署，负责 ID 注册、心跳、信令转发。保持轻量。
- **hbbr (中继集群)**：多节点分散部署（北京/上海/广州等），纯粹 TCP/WSS 流量转发。
- **Admin API**：充当"节点目录服务器"，从数据库查询可用中继列表返回给客户端。
- **客户端**：主控端拉取节点列表 → 二维竞速 → 发号施令给被控端。

### Ed25519 密钥同步（防 Key Mismatch）

所有 hbbs 和 hbbr 节点**必须使用同一对** `id_ed25519` / `id_ed25519.pub`：

```
中心节点生成密钥 → 通过 docker-compose 挂载或 CI/CD 分发 → 覆盖所有边缘 hbbr 节点
```

部署方式：在 docker-compose 中将密钥文件以 volume 形式挂载到所有 hbbs/hbbr 容器：

```yaml
services:
  hbbs:
    volumes:
      - ./keys:/data  # id_ed25519 + id_ed25519.pub
  relay-bj:
    volumes:
      - ./keys:/data  # 同一对密钥
  relay-sh:
    volumes:
      - ./keys:/data  # 同一对密钥
```

---

## Part A: 节点池下发（双通道）

### A1. Admin API 下发（冷启动）

**文件**: `rustdesk-api/http/controller/api/webClient.go`

`ServerConfig()` / `ServerConfigV2()` 响应新增 `relay_servers` 数组：

```json
{
  "id_server": "...",
  "key": "...",
  "relay_server": "...",
  "relay_servers": [
    {
      "name": "北京-联通",
      "region": "CN-North",
      "relay_server": "bj.example.com:21117",
      "ws_host": "wss://bj.example.com:21119",
      "priority": 10,
      "cost_weight": 1,
      "support_wss": true
    }
  ],
  "api_server": "...",
  "ws_host": "..."
}
```

**文件**: `rustdesk-api/service/serverScheduler.go`
新增方法：`GetAllAvailableServers() []*RelayServerEntry` — 查询所有 `is_active=true AND is_online=true` 的中继服务器。

### A2. hbbs ConfigUpdate 推送（热更新）

**文件**: `libs/hbb_common/protos/rendezvous.proto`

```protobuf
message RelayServerEntry {
  string name = 1;
  string region = 2;
  string relay_server = 3;  // host:21117
  int32 priority = 4;
  int32 cost_weight = 5;
  bool support_wss = 6;
  string ws_host = 7;       // wss://host:21119
}

message ConfigUpdate {
  int32 serial = 1;
  repeated string rendezvous_servers = 2;
  repeated RelayServerEntry relay_servers = 3;   // 新增
}
```

hbbs 在 `fetch_server_topology()` 周期中从 Redis 拉取并推送给所有在线客户端。中继上下线实时生效，无需客户端轮询。

---

## Part B: 主控端独裁制（信令链路改造）

### B1. Protobuf 改动

```protobuf
message PunchHoleRequest { 
  string id = 1; 
  NatType nat_type = 2;
  string licence_key = 3;
  ConnType conn_type = 4;
  string token = 5;
  string version = 6;
  int32 udp_port = 7;
  bool force_relay = 8;
  int32 upnp_port = 9;
  bytes socket_addr_v6 = 10;
  string relay_server = 11;    // 新增：主控端指定的中继 (host:21117 或 host:21119)
  bool relay_use_wss = 12;     // 新增：是否走 WSS 协议
}
```

`relay_use_wss` 告诉被控端应该用 TCP+ChaCha20 还是 WSS+TLS 去连中继。

### B2. hbbs 改造：转发而非决策

**文件**: `rustdesk-server/src/rendezvous_server.rs` (`handle_punch_hole_request`, line 1120)

```rust
// 主控端指定了中继 → 直接使用，不走服务端评分
let mut relay_server = if !ph.relay_server.is_empty() {
    ph.relay_server.clone()
} else {
    self.get_relay_server(addr.ip(), peer_addr.ip()).await  // 兼容旧客户端
};
```

`PunchHole` 消息本身已有 `relay_server` 字段（line 35），不需要改。`relay_use_wss` 信息可以通过新增 `PunchHole` 的 bool 字段传递，或者编码在 relay_server 中（末尾端口 21119 = WSS）。

### B3. 被控端改造：无条件跟随 + 协议自适应

**文件**: `rustdesk-client/src/rendezvous_mediator.rs`

`get_relay_server()` (line 801) 优先级反转：

```rust
fn get_relay_server(&self, provided_by_rendezvous_server: String) -> String {
    // 主控端(通过hbbs)指定的优先 — 无条件跟随
    let mut relay_server = provided_by_rendezvous_server;
    if relay_server.is_empty() {
        relay_server = Config::get_option("relay-server");  // 本地配置兜底
    }
    if relay_server.is_empty() {
        relay_server = crate::increase_port(&self.host, 1);
    }
    relay_server
}
```

收到 `PunchHole` 时检查 `relay_use_wss` 决定用 TCP 还是 WSS 连接中继。

### B4. 主控端改造：三种模式 + 竞速策略

**文件**: `rustdesk-client/src/client.rs`

**新增 `ConnectionMode` 枚举**：

```rust
#[derive(Clone, Copy, PartialEq)]
pub enum ConnectionMode {
    Auto,       // 自动择优：并发竞速 TCP ∥ WSS
    TcpOnly,    // 专业版：仅 TCP+ChaCha20
    WssOnly,    // 企业校园版：仅 WSS
}
```

**`_start()` 连接入口改造**（原 line ~460 之前插入）：

```rust
// ① P2P 直连优先（始终启用，不可禁用）
let p2p_result = tokio::time::timeout(
    Duration::from_secs(3),
    Self::try_p2p_direct(peer_addr, udp_port, ipv6_addr)
).await;

if let Ok(Ok(stream)) = p2p_result {
    return Ok(stream);  // P2P 成功 → 零成本，直接使用
}

// ② P2P 失败 → 根据用户选择的场景模式竞速中继
let relay_list = self.fetch_relay_server_list().await;
let mode = self.connection_mode();  // 来自用户设置

let (best_node, best_protocol, stream) = match mode {
    ConnectionMode::Auto    => race_concurrent(&relay_list).await?,
    ConnectionMode::TcpOnly => race_fallback_tcp(&relay_list).await?,
    ConnectionMode::WssOnly => race_fallback_wss(&relay_list).await?,
};

// ③ 发号施令
ph.relay_server = best_node.relay_addr_for(&best_protocol);
ph.relay_use_wss = best_protocol.is_wss();
```

**策略 A — 自动择优：并发竞速 `race_concurrent()`**（推荐，最快）：

```rust
/// N 个中继 × 2 种协议，同时发探针，首个完成握手者胜出
async fn race_concurrent(nodes: &[RelayNode]) -> Result<(Node, Protocol, Stream)> {
    let mut tasks = FuturesUnordered::new();

    for node in nodes {
        // 探针 A: TCP + ChaCha20 (500ms 超时 — 少一次 TLS 握手)
        tasks.push(probe_tcp_encrypted(node.clone(), 500));
        // 探针 B: WSS + TLS (800ms 超时 — 多 TLS 握手)
        if node.support_wss {
            tasks.push(probe_wss(node.clone(), 800));
        }
    }

    while let Some(result) = tasks.next().await {
        if let Ok(ok) = result {
            cancel_all_pending(tasks);  // 销毁慢的连接
            return Ok(ok);
        }
    }
    Err("所有线路和协议全部不可达")
}
```

**策略 B — 专业版 / 企业校园版：快速降级 `race_fallback_*()`**（最省资源）：

```rust
/// 300ms 快速降级：TCP 不通立刻切下一个节点
async fn race_fallback_tcp(nodes: &[RelayNode]) -> Result<(Node, Protocol, Stream)> {
    for node in nodes {
        let result = tokio::time::timeout(
            Duration::from_millis(300),
            connect_tcp_encrypted(node)
        ).await;
        if let Ok(Ok(stream)) = result {
            return Ok((node.clone(), Protocol::TcpEncrypted, stream));
        }
    }
    Err("所有 TCP 节点不可达")
}

/// WSS 版本同理，300ms 超时逐个尝试
async fn race_fallback_wss(nodes: &[RelayNode]) -> Result<(Node, Protocol, Stream)> {
    for node in nodes.iter().filter(|n| n.support_wss) {
        let result = tokio::time::timeout(
            Duration::from_millis(300),
            connect_wss(node)
        ).await;
        if let Ok(Ok(stream)) = result {
            return Ok((node.clone(), Protocol::Wss, stream));
        }
    }
    Err("所有 WSS 节点不可达")
}
```

**竞速策略对比**：

| 策略 | 适用模式 | 探针数 | 特点 |
|------|---------|--------|------|
| 并发竞速 | 自动择优 | 2N | 最快，同时发所有探针，先通先用 |
| 快速降级 | 专业版/企业版 | N | 最省资源，300ms 超时逐个尝试 |

---

## Part C: 传输层加密（防 DPI 指纹识别）

单端口 443 嗅探：

```
客户端连 :443
   │
   ├─ TCP+ChaCha20 探针 → 直接加密流
   └─ WSS+TLS 探针     → TLS 握手 → WS 升级
```

服务端 `protocol_multiplex.rs` 嗅探首 3 字节分发。

**文件**: `hbb_common/src/stream_encrypt.rs` 新建  
**文件**: `rustdesk-server/src/protocol_multiplex.rs` 新建

---

## Part D: 手动指定中继节点（隐藏高级选项）

仅当用户主动进入"高级设置 → 中继节点"时展示。

根据当前选中的"场景模式"，展示不同的可选节点：

```
场景模式 = 自动择优 → 显示全部节点（TCP + WSS 混合列表）
场景模式 = 专业版   → 仅显示 TCP 节点
场景模式 = 企业校园 → 仅显示 WS 节点
```

- 用户手动选择节点后，**跳过竞速**，直接用指定节点连接
- 携带当前场景模式自动确定的协议（TCP 或 WSS）
- 底层仍走 `PunchHoleRequest.relay_server` + `relay_use_wss`
- 注意：手动选节点后，场景模式对该次连接不再生效

---

## 实施顺序

```
Phase 1: Protobuf 扩展
  - PunchHoleRequest 加 relay_server(11) + relay_use_wss(12)
  - ConfigUpdate 加 RelayServerEntry + relay_servers(3)
  - 重新生成 Rust + Go 代码

Phase 2: Admin API 节点池下发（Go 端）
  - serverScheduler.go: GetAllAvailableServers()
  - server.go: GetRelayServerList()
  - webClient.go: ServerConfig/ServerConfigV2 返回 relay_servers

Phase 3: hbbs 转发策略
  - rendezvous_server.rs: handle_punch_hole_request() 跳过服务端决策
  - ConfigUpdate 推送 relay_servers

Phase 4: 被控端无条件跟随
  - rendezvous_mediator.rs: get_relay_server() 优先级反转
  - handle_punch_hole() 根据 relay_use_wss 选择协议

Phase 5: 主控端二维竞速（核心）
  - client.rs: race_2d_matrix() 新建
  - client.rs: _start() 中 P2P优先 → 竞速 → 发号施令
  - 从 admin API / ConfigUpdate 获取 relay_list

Phase 6: 传输层加密
  - hbb_common/src/stream_encrypt.rs 新建
  - rustdesk-server/src/protocol_multiplex.rs 新建

Phase 7: 编译 + 部署 + 验证
```

## 验证

1. **P2P 优先**: 同局域网两台机器 → P2P 直连成功 → 不触发中继
2. **二维竞速**: 公网环境 → TCP+ChaCha20 探针先胜出（~30ms 握手）
3. **防火墙降级**: 公司内网 → TCP 探针超时 → WSS 探针胜出 → 自动接管
4. **被控端跟随**: 主控端选"北京-WSS" → 被控端日志 "Following viewer: relay=bj:21119 protocol=wss"
5. **管理端下发**: `curl /api/server-config-v2` → `relay_servers` 数组非空
6. **流量白噪声**: tcpdump → TCP+ChaCha20 无任何协议特征
7. **旧客户端兼容**: 旧客户端不填 relay_server → hbbs 走原有评分逻辑
