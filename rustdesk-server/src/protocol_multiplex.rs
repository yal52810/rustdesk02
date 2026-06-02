//! 协议嗅探 + 多路分发
//!
//! 单端口 443 同时处理两种连接：
//! - TLS ClientHello (0x16 0x03 ...) → WSS 路径
//! - 白噪声 (其他字节) → ChaCha20-Poly1305 路径
//!
//! 用法: `io_loop` 中 `listener.accept()` 后调用 `accept_multiplex()` 替代直接 TCP 处理。

use hbb_common::{
    log,
    sodiumoxide::crypto::box_,
    stream_encrypt::EncryptedStream,
    tcp::FramedStream,
    tokio::net::TcpStream,
    ResultType,
};

/// 多路复用后的流类型
pub enum MultiplexedStream {
    /// TLS → WebSocket Secure 路径（WSS 专业线路）
    Wss(FramedStream),
    /// ChaCha20-Poly1305 加密路径（公网白噪声）
    Encrypted(EncryptedStream<FramedStream>),
}

/// 检测首 3 字节是否为 TLS ClientHello
///
/// TLS 握手首字节固定为 0x16，第 2-3 字节为协议版本 0x03 0x00~0x03
#[inline]
fn looks_like_tls(peek: &[u8; 3]) -> bool {
    peek[0] == 0x16 && peek[1] == 0x03 && peek[2] <= 0x03
}

/// 协议嗅探 + 分发
///
/// peek 首 3 字节：
/// - TLS ClientHello (0x16 0x03 0x00~03) → 保持原始流，由上层做 TLS accept + WS 升级
/// - 其他 → 读取 2 字节 magic prefix [0x00, 0x00] 后启动 ChaCha20-Poly1305 解密
///
/// ChaCha20 首字节有 ~1/65536 概率恰好 = 0x16 0x03，导致误判。
/// 加密流在首包前加 2 字节 magic [0x00, 0x00] 消除此概率。
pub async fn accept_multiplex(
    mut stream: TcpStream,
    pk: Option<&box_::PrecomputedKey>,
) -> ResultType<MultiplexedStream> {
    use hbb_common::tokio::io::AsyncReadExt;

    let peer_addr = stream.peer_addr()?;
    let mut peek = [0u8; 3];
    stream.peek(&mut peek).await?;

    if looks_like_tls(&peek) {
        log::info!("protocol_multiplex: TLS ClientHello detected → WSS path");
        let framed = FramedStream::from(stream, peer_addr);
        Ok(MultiplexedStream::Wss(framed))
    } else {
        log::debug!("protocol_multiplex: white noise → ChaCha20-Poly1305 path");
        // Consume 2-byte magic prefix added by client to avoid false TLS detection
        let mut magic = [0u8; 2];
        stream.read_exact(&mut magic).await?;

        let framed = FramedStream::from(stream, peer_addr);
        if let Some(key) = pk {
            Ok(MultiplexedStream::Encrypted(EncryptedStream::new(
                framed, key,
            )))
        } else {
            // No precomputed key → fallback to plain FramedStream
            log::warn!("protocol_multiplex: no key provided, using unencrypted stream");
            Ok(MultiplexedStream::Wss(framed))
        }
    }
}
