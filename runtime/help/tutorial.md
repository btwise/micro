# 教程

这是一个简短的介绍,微的配置系统,将给一些
简单示例显示如何配置设置、重新绑定键和使用
的 `init.lua` 配置微到您的喜好.

希望这对你有用.

参见 `> help defaultkeys` 以获得默认键绑定的列表和解释.

### 设置

在micro中,您的设置存储在`~/.config/micro/settings.json`,文件
这是在您第一次运行micro时创建的.它是一个json文件,包含所有
设置及其值.要更改选项,可以在 `settings.json` 里更改值,或者你可以在使用micro时直接输入它.

按Ctrl-e进入命令模式,并键入 `set option value` (在
未来,我将使用 `> set option value` 来指示按Ctrl-e).这个更改
将立即生效,并将被保存到 `settings.json` 文件
这样即使你关闭micro,设置也会保持不变.

您也可以在本地设置选项,这意味着设置将只有
你在缓冲区中设置的值.例如,如果你有两个
打开的拆分窗口,然后输入 `> setlocal tabsize 2`,字符缩进在当前缓冲区将只有2,
同时micro也不会保存这个局部的改动 `settings.json` 文件.方法中本地设置选项
`settings.json` 文件.例如,如果您希望在Ruby文件和4个其他文件 `tabsize` 为2,你可以在
`settings.json`中放入以下内容:

```json
{
    "*.rb": {
        "tabsize": 2
    },
    "tabsize": 4
}
```

Micro将只在匹配 `*.rb`的文件中将 `tabsize` 设置为2.

如果您想了解所有可用选项的更多信息,请参阅
`options`主题(`> help options`).

### 键绑定

键绑定的工作方式与选项非常相似.可以用`~/.config/micro/bindings.json`属性来配置它们.

例如,如果你想绑定 `Ctrl-r` 来重做,你可以在 `bindings.json` 中这样配置:

```json
{
    "Ctrl-r": "Redo"
}
```

非常简单.

你也可以在micro中使用 `> bind key action` 命令绑定键,
但使用该命令所做的绑定不会保存到 `bindings.json` 文件.

有关键绑定的更多信息,如可以绑定哪些键以及可以绑定什么动作是可用的,
请参阅 `keybindings` 帮助主题(`> help keybindings`).

### 使用Lua进行配置

如果你需要比json文件更强大的功能,你可以使用 `init.lua` 
文件.在 `~/.config/micro` 中创建它.这个文件是一个lua文件,当
Micro启动,本质上是一个单文件插件.插件名是 `initlua`.

这个例子将告诉你如何使用 `init.lua` 的文件,创建一个绑定
到`Ctrl-r`,假设当前文件是一个Go文件,它将对当前文件执行bash命令 `go run`.

你可以通过在 `init.lua` 中放入以下内容来实现:

```lua
local config = import("micro/config")
local shell = import("micro/shell")

function init()
    -- true means overwrite any existing binding to Ctrl-r
    -- this will modify the bindings.json file
    config.TryBindKey("Ctrl-r", "lua:initlua.gorun", true)
end

function gorun(bp)
    local buf = bp.Buf
    if buf:FileType() == "go" then
        -- the true means run in the foreground
        -- the false means send output to stdout (instead of returning it)
        shell.RunInteractiveShell("go run " .. buf.Path, true, false)
    end
end
```

或者,你可以摆脱 `TryBindKey`行,并把这一行放在`bindings.json` 文件中:

```json
{
    "Ctrl-r": "lua:initlua.gorun"
}
```

有关插件和micro使用的lua系统的更多信息,请参阅
`plugins` 帮助主题(`> help plugins`).