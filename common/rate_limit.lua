-- rate_limit.lua
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])  -- 单位：秒

local current = redis.call("INCR", key)
if current == 1 then
    -- 首次访问，设置过期时间
    redis.call("EXPIRE", key, window)
end

if current > limit then
    return 0  -- 超限
else
    return current  -- 当前计数
end