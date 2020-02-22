-- 系统函数
-- Created by IntelliJ IDEA.
-- User: mercu
-- Date: 2019/12/8
-- Time: 14:38
-- To change this template use File | Settings | File Templates.
--
--每个文件最开始定义了以下变量
--[[
local _PACKAGE = {} --包内容
 ]]

local strings = string
local meta = getmetatable("")
meta.__add = function(a, b)
    return a .. b
end

local slice_meta = {
    __type = "slice",
    __index = function(t, key)
        return t.data[key + t.start]
    end,
    __newindex = function(t, key, value)
        t.data[key + t.start] = value
    end
}
local function is_slice(t)
    local meta = getmetatable(t)
    return meta == slice_meta
end

local function utf8len(byte)
    local bit = require"bit"
    if bit.band(byte, 0x80) == 0 then
        return 1
    end

    local count = 0
    for i = 1, 8 do
        count = count + 1
        if bit.band(bit.lshift(byte, i), 0x80) == 0 then
            break
        end
    end
    return count
end

local function copy_table(dst, src)
    for index, value in pairs(src) do
        dst[index] = value
    end
end

function _G.import(path)
    return require(path)
end

function _G.checkType(value, defType)
    return value
end

function _G.new(define)
    return define()
end

function _G.rangestr(value)
    local charSeq = {}
    local index = 1
    while index <= strings.len(value) do
        local byte = strings.byte(value, index)
        local len = utf8len(byte)
        local utf8Char = strings.sub(value, index, index + len - 1)
        table.insert(charSeq, utf8Char)
        index = index + len
    end
    return ipairs(charSeq)
end

function _G.range(value)
    if type(value) == "string" then
        return rangestr(value)
    elseif is_slice(value) then
        return ipairs(value.data)
    else
        return pairs(value)
    end
    return nil
end

local metakey = {
    __name = "string",
    __type = "string",
    __mtable = "table",
}
function _G.struct(define)
    define.__mtable = {}
    define.__type = "struct"
    return setmetatable(define, {
        __tostring = function() return  define.__type .. ":".. define.__name end,
        --定义构造函数
        __call = function(def, params)
            local ret = {}
            local objName = tostring(ret) .. "#" ..define.__name
            local embed = {}
            --设置成员初始值
            for key, value in pairs(def) do
                if type(key) == "number" and value.__name then
                    key = value.__name
                    embed[key] = value
                end

                if not metakey[key] then
                    if params and params[key] then
                        ret[key] = params[key]
                    elseif value.__type == "interface" then
                        ret[key] = nil
                     else
                        ret[key] = value()
                    end
                end
            end

            --设置成员方法(会隐藏内嵌结构方法)
            for name, func in pairs(def.__mtable) do
                ret[name] = func
            end
            --内嵌结构
            setmetatable(ret, {
                __tostring = function() return objName end,
                __newindex = function(t, k, v)
                    local foundEs
                    for name, es in pairs(embed) do
                        local esObj = rawget(t, name)
                        local value = esObj[k]
                        if value or (es[k] and es[k].__name == "interface") then
                            if foundEs then
                                error(string.format("repeat member name[%s] in embed struct", k))
                            end
                            esObj[k] = v
                            foundEs = true
                        end
                    end
                    if not foundEs and def[k].__type == "interface" then
                        rawset(t, k, v)
                    end
                end,
                __index = function(t, k)
                    local foundValue
                    for name,_ in pairs(embed) do
                        local value = rawget(t, name)[k]
                        if value then
                            if foundValue then
                                error(string.format("repeat member name[%s] in embed struct", k))
                            end
                            foundValue = value
                        end
                    end
                    return foundValue
                end
            })
            return ret
        end
    })
end

function _G.interface(define)
    define.__type = "interface"
    return setmetatable(define, {
        __tostring = function() return  define.__type .. ":".. define.__name end
    })
end

function _G.method(...)
    return {...}
end

function _G.as(obj, inf)
    return obj
end

--will change old silce, for performance
function _G.append(t, ...)
    if is_slice(t) then
        for _, v in ipairs({...}) do
            table.insert(t.data, v, t.start + t.len)
            t.len = t.len + 1
        end
        return t
    else
        error("not a slice")
    end
end

function _G.copy(dst, src)
    if is_slice(dst) and is_slice(src) then
        dst.start = 0
        dst.len = src.len
        dst.data = {}
        copy_table(dst.data, src.data)
    else
        error("not a slice")
    end
end

function _G.slice(t, iStart, iEnd)
    iStart = iStart or 0
    iEnd = iEnd or #t
    local newt = {
        start = iStart + 1,
        len = iEnd - iStart,
        data = t,
    }

    return setmetatable(newt, slice_meta)
end

function _G.unpack_slice(t)
    if is_slice(t) then
        return unpack(t.data)
    else
        error("not a slice")
    end
end

function _G.len(value)
    if type(value) == "string" then
        return utf8len(value)
    elseif is_slice(value) then
        return value.len
    else
        return select("#", value)
    end
end

local array_meta = {
    __call = function(t, param)
        for index = len(param) + 1, t.length do
            table.insert(param, t.eleType())
        end
        return param
    end
}
function _G.array_type(eleType, num)
    local tArray = {}
    tArray.eleType = eleType
    tArray.length = num
    return setmetatable(tArray, array_meta)
end

local slice_meta = {
    __call = function(t, param)
        for index = len(param) + 1, t.length do
            table.insert(param, t.eleType())
        end
        return _G.slice(param)
    end
}

function _G.slice_type(eleType)
    local tArray = {}
    tArray.eleType = eleType
    tArray.length = 0
    return setmetatable(tArray, slice_meta)
end

function _G.make(defType, ...)
    if getmetatable(defType) == slice_meta then
        local length = ...
        local slice_struct = defType()
        slice_struct.length = length
        return slice_struct({})
    end
    return {}
end

function _G.map(param1)
    return param1
end

----------------------------------
--内建类型
local function buildIn(name, value)
    local nameType = _G[name]
    if not nameType then
        nameType = {__name = name}
    end
    _G[name] = setmetatable(nameType, {
        __call = function()
            return value
        end,
        __tostring = function() return "buildin:" .. name end
    })
end
buildIn("number", 0)
buildIn("string", "")
buildIn("bool", false)

----------------------------------
--类型映射
--typedefine(_G, "float64", number)
