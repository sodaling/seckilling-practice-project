local buyNum = ARGV[1]
local goodsKey = KEYS[1]
local goodsNum = redis.call('get',goodsKey)
if goodsNum >= buyNum
then redis.call('decrby',goodsKey,buyNum)
    return goodsNum
else
    return '0'
end