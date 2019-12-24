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

function _G.package(name)
    return {__name = name, __type == "package"}
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
        --定义构造函数
        __call = function(def, params)
            local ret = {}
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

            --设置内嵌类型的方法
            for embedName, embedStruct in pairs(embed) do
                for name, func in pairs(embedStruct.__mtable) do
                    ret[name] = function(self, ...)
                        return func(ret[embedName], ...)
                    end
                end
            end

            --设置成员方法(会隐藏内嵌结构方法)
            for name, func in pairs(def.__mtable) do
                ret[name] = func
            end
            --
            return ret
        end
    })
end

function _G.interface(define)
    define.__type = "interface"
    return define
end

function _G.method(...)
    return {...}
end

function _G.as(obj, inf)
    return obj
end

function _G.slice(t, iStart, iEnd)
    local newt = {}
    for i = iStart, iEnd - 1 do
        table.insert(newt, t[i])
    end
    return newt
end

function _G.len(value)
    if type(value) == "string" then
        return utf8len(value)
    else
        return select("#", value)
    end
end

function _G.make(defType, ...)
    return defType(...)
end

function _G.array(param1, param2)
    if type(param1) == "table" then
        return param1
    end
    return {}
end

function _G.map()
    return {}
end

----------------------------------
--内建类型
local function buildIn(name, value)
    _G[name] = setmetatable({__name = name}, {
        __call = function()
            return value
        end
    })
end
buildIn("number", 0)
buildIn("string", "")
buildIn("bool", false)
----------------------------------
--类型映射
--typedefine(_G, "float64", number)
