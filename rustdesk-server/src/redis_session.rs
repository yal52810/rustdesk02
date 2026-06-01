/// Attempt to register a session for `user_id`/`device_id` with concurrent `limit`.
/// This will add/update the device's score to current timestamp, then if the
/// cardinality exceeds `limit` it will pop the oldest entries (ZPOPMIN).
/// If the popped entries include the current device, the registration fails
/// (returns false). Returns true on success or when Redis is unavailable.
pub async fn check_and_register(user_id: i32, device_id: &str, limit: i32) -> bool {
    // Build redis URL from environment or default
    let redis_url = std::env::var("REDIS_URL").unwrap_or_else(|_| "redis://127.0.0.1/".to_string());

    // Try to create client and connection; if anything fails, allow the registration
    let client = match redis::Client::open(redis_url.as_str()) {
        Ok(c) => c,
        Err(_) => return true,
    };
    let mut conn = match client.get_async_connection().await {
        Ok(c) => c,
        Err(_) => return true,
    };

    // Use a Lua script to perform the operations atomically:
    // 1) ZADD key score member
    // 2) ZCARD key
    // 3) if card > limit then ZPOPMIN key (card-limit)
    // 4) return 1 if member still present, 0 if it was popped
    let key = format!("rustdesk:user:{}:sessions", user_id);
    let now = chrono::Utc::now().timestamp_millis();

    let lua = r#"
    local key = KEYS[1]
    local member = ARGV[1]
    local score = ARGV[2]
    local limit = tonumber(ARGV[3])
    redis.call('ZADD', key, score, member)
    local cnt = redis.call('ZCARD', key)
    if cnt <= limit then
        return 1
    end
    local n = cnt - limit
    local popped = redis.call('ZPOPMIN', key, n)
    -- popped is a flat array: member1, score1, member2, score2, ...
    for i=1,#popped,2 do
        if popped[i] == member then
            return 0
        end
    end
    return 1
    "#;

    let script = redis::Script::new(lua);
    let res: Result<i32, _> = script.key(&key).arg(device_id).arg(now).arg(limit).invoke_async(&mut conn).await;
    match res {
        Ok(v) => v == 1,
        Err(_) => true,
    }
}
