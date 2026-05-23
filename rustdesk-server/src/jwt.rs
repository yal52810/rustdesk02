use jsonwebtoken::{decode, encode, Algorithm, DecodingKey, EncodingKey, Header, Validation};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use std::env;

pub static SECRET: Lazy<String> =
    Lazy::new(|| env::var("RUSTDESK_API_JWT_KEY").unwrap_or_else(|_| "".to_string()));

// 定义一个结构体来表示 JWT 的 payload
// 1. 将结构体设为 'pub'，以便其他模块 (如 rendezvous_server) 可以访问
#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    // 2. 将 'user_id' 设为 'pub'
    // 3. 将 'user_id' 类型改为 i32，以匹配您 database.rs 中 'get_subscription_status(user_id: i32)' 的函数签名
    pub user_id: i32,
    // 4. 添加 'pub uuid' 字段，以修复 'no field 'uuid'' 错误
    pub uuid: String,
    pub exp: usize,
}

// 5. 更新 generate_token 函数签名，以接收 'uuid' 和 'i32' 类型的 'user_id'
pub fn generate_token(user_id: i32, uuid: String, exp: i64) -> Result<String, String> {
    println!("secret: {:}", SECRET.to_string());
    let claims = Claims {
        user_id,
        uuid, // <-- 添加 uuid 到 claims
        exp: (chrono::Utc::now() + chrono::Duration::seconds(exp)).timestamp() as usize,
    };

    let token = encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(SECRET.as_ref()),
    );

    match token {
        Ok(t) => Ok(t),
        Err(e) => Err(e.to_string()),
    }
}

// 验证 JWT 的函数
// (此函数无需修改，它会自动使用更新后的 Claims 结构体)
pub fn verify_token(token: &str) -> Result<Claims, String> {
    // 解码 JWT
    let validation = Validation::new(Algorithm::HS256);

    let decoded = decode::<Claims>(
        &token,
        &DecodingKey::from_secret(SECRET.as_ref()),
        &validation,
    );
    match decoded {
        Ok(token_data) => {
            let now = chrono::Utc::now().timestamp() as usize;
            if token_data.claims.exp > now {
                Ok(token_data.claims)
            } else {
                Err("Token status invalid or expired".to_string())
            }
        }
        Err(e) => Err(format!("Invalid token: {}", e)), // 最好将错误 'e' 打印出来
    }
}