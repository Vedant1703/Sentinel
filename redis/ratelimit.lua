-- KEYS[1] = redis key (user/ip)
-- ARGV[1] = window in seconds
-- ARGV[2] = limit

local current = redis.call("GET", KEYS[1])

if not current then
  redis.call("SET", KEYS[1], 1)
  redis.call("EXPIRE", KEYS[1], ARGV[1])
  return 1
end

if tonumber(current) < tonumber(ARGV[2]) then
  redis.call("INCR", KEYS[1])
  return 1
end

return 0
