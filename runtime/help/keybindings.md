# 键绑定

Micro 有大量的热键,使它容易和强大的使用和所有
热键完全可根据您的喜好定制.

如果使用' > bind '进行更改,则自定义键绑定将存储在micro内部
命令或也可以添加在文件`~/.config/micro/bindings.json`
下面讨论.所使用的json格式的默认键绑定列表
micro,请参见文件末尾.以获得更用户友好的列表
说明什么是默认热键和他们做什么,请参阅
' >help defaultkeys '(包括一个json格式的默认键列表
在本文档的末尾).

如果~ /.config/micro/bindings.json '不存在,你可以简单地创建它.
Micro会知道该怎么做的.

你可以使用Ctrl +箭头逐字移动(Mac使用Alt +箭头).Alt +左右
移动光标到行的开始和结束(Ctrl +左/右Mac),和Ctrl +上下移动
光标指向缓冲区的开始和结束.

你可以按住shift键,在移动时选择所有这些移动动作.

## 重新绑定键

绑定可以使用`~/.config/micro/bindings.json`文件重新绑定.
每个键都绑定到一个动作.

例如,要绑定' Ctrl-y '来撤销,绑定' Ctrl-z '来重做,可以将
下面是`bindings.json`文件.

```json
{
	"Ctrl-y": "Undo",
	"Ctrl-z": "Redo"
}
```

**注意:**语法`<Modifier><key>`等价于`<Modifier>-<key>`.在
此外,终端不支持Ctrl-Shift绑定,它们是一样的
简单的Ctrl绑定.这意味着`CtrlG`, `Ctrl-G`和`Ctrl-g`所有
意思相同.然而,对于Alt,情况并非如此:`AltG`和`Alt-G`
意思是`Alt-Shift-g`,而`Alt-g`不需要Shift修饰符.

除了编辑你的`~/.config/micro/bindings.json`,你可以运行
`>bind <keycombo> <action>`有关可绑定动作的列表,请参见下面.

您还可以在重新绑定时链接命令.例如,如果你想让Alt-s
保存并退出你可以像这样绑定它:

```json
{
    "Alt-s": "Save,Quit"
}
```

每个操作将返回一个成功标志.动作可以这样链接
只有成功或失败,或者两者皆有,这个链条才会继续.
`,`分隔符将始终链接到下一个操作.
如果之前的操作成功,将`|`分隔符中止链,
如果之前的操作失败,则`&`将中止链.例如,在默认情况下
绑定,TAB被绑定为

```
"Tab": "Autocomplete|IndentSelection|InsertTab"
```

这意味着如果`Autocomplete`操作成功,链将中止.否则,它将尝试`IndentSelection`,如果这也失败了,它
将执行`InsertTab`.

## 绑定命令

还可以绑定密钥以在命令模式下执行命令(请参阅(`help commands`).只需在绑定前加上`command:`.例如:

```json
{
    "Alt-p": "command:pwd"
}
```

** macOS注意事项**:默认情况下,macOS终端不转发alt事件和
而是插入unicode字符.要解决这个问题,请执行以下操作:

* iTerm2:在`首选项->配置文件->键`中为“左选项键”选择`Esc+`.
* 终端.app:在`首选项->配置文件->键盘`中启用`使用选项键作为元键`.

现在,当你按下`Alt-p` `Alt-p`命令将被执行,将显示信息栏中的工作目录.

你也可以用`command-edit:`绑定一个“可编辑”命令.这意味着
当您按下绑定时,Micro不会立即执行命令,但是相反,只需将字符串放在命令模式下的信息栏中.例如,
你可以重新绑定 `Ctrl-g` 到 `> help`:

```json
{
    "Ctrl-g": "command-edit:help "
}
```

现在,当你按下`Ctrl-g`,`help`将出现在命令栏和你的光标将放在它的后面(注意json中的空格,该空格控制
光标位置).

## 绑定原始转义序列

对象上的绑定键感兴趣时,请阅读本节
支持的绑定键列表.

使用基于终端的编辑器的缺点之一是编辑器必须
通过终端获取关于关键事件的所有信息.终端
经常(但不总是)以转义序列的形式发送这些事件
从`0x1b`开始.

例如,如果micro读取`\x1b[1;5D`,在大多数终端上,这将意味着
用户按ctrl + left.

但是对于许多关键和弦,终端不会发送任何转义代码
发送一个已经在使用的转义码.例如,对于`CtrlBackspace`, 我的
终端发送`\u007f`(注意这不是以`0x1b`开头),它也是
发送`Backspace`意味着微不能绑定`CtrlBackspace`.

但是,有些终端确实允许绑定键来发送特定的转义
您定义的序列.然后从微你可以直接绑定那些转义
动作的序列.例如,要绑定' CtrlBackspace ',您可以指示
您的终端发送`\x1bctrlback`,然后将其绑定到`bindings.json`中:

```json
{
    "\u001bctrlback": "DeleteWordLeft"
}
```

下面是在不同终端发送原始转义的一些说明

### iTerm2

在iTerm2中,您可以在`Preferences->Profiles->Keys`中执行此操作,然后单击
`+`,输入你的键绑定,并为`Action`选择`Send Escape Sequence`.
对于上面的例子,你可以在框中输入`ctrlback` (`\x1b`)由iTerm2自动发送.

### 使用loadkey的Linux

你可以在linux中使用loadkeys程序来做到这一点.


## 解除绑定键

方法也可以禁用任何默认键绑定
用户的`bindings.json`文件中的`None`操作.

## 可绑定动作和可绑定键

默认键绑定列表包含大多数可能的操作和键
你可以用,但不是所有的.以下是两者的完整列表.

可能采取的行动的完整列表:

```
CursorUp
CursorDown
CursorPageUp
CursorPageDown
CursorLeft
CursorRight
CursorStart
CursorEnd
SelectToStart
SelectToEnd
SelectUp
SelectDown
SelectLeft
SelectRight
SelectToStartOfText
SelectToStartOfTextToggle
WordRight
WordLeft
SelectWordRight
SelectWordLeft
MoveLinesUp
MoveLinesDown
DeleteWordRight
DeleteWordLeft
SelectLine
SelectToStartOfLine
SelectToEndOfLine
InsertNewline
InsertSpace
Backspace
Delete
Center
InsertTab
Save
SaveAll
SaveAs
Find
FindLiteral
FindNext
FindPrevious
Undo
Redo
Copy
CopyLine
Cut
CutLine
DuplicateLine
DeleteLine
IndentSelection
OutdentSelection
OutdentLine
IndentLine
Paste
SelectAll
OpenFile
Start
End
PageUp
PageDown
SelectPageUp
SelectPageDown
HalfPageUp
HalfPageDown
StartOfLine
EndOfLine
StartOfText
StartOfTextToggle
ParagraphPrevious
ParagraphNext
ToggleHelp
ToggleDiffGutter
ToggleRuler
JumpLine
ClearStatus
ShellMode
CommandMode
Quit
QuitAll
AddTab
PreviousTab
NextTab
NextSplit
Unsplit
VSplit
HSplit
PreviousSplit
ToggleMacro
PlayMacro
Suspend (Unix only)
ScrollUp
ScrollDown
SpawnMultiCursor
SpawnMultiCursorUp
SpawnMultiCursorDown
SpawnMultiCursorSelect
RemoveMultiCursor
RemoveAllMultiCursors
SkipMultiCursor
None
JumpToMatchingBrace
Autocomplete
```

`StartOfTextToggle` 和 `SelectToStartOfTextToggle` 动作之间切换
跳转到文本的开头(第一个)和行的开头.

您还可以绑定一些鼠标操作(这些操作必须绑定到鼠标按钮)

```
MousePress
MouseMultiCursor
```

下面是你可以绑定的所有可能键的列表:

```
Up
Down
Right
Left
UpLeft
UpRight
DownLeft
DownRight
Center
PageUp
PageDown
Home
End
Insert
Delete
Help
Exit
Clear
Cancel
Print
Pause
Backtab
F1
F2
F3
F4
F5
F6
F7
F8
F9
F10
F11
F12
F13
F14
F15
F16
F17
F18
F19
F20
F21
F22
F23
F24
F25
F26
F27
F28
F29
F30
F31
F32
F33
F34
F35
F36
F37
F38
F39
F40
F41
F42
F43
F44
F45
F46
F47
F48
F49
F50
F51
F52
F53
F54
F55
F56
F57
F58
F59
F60
F61
F62
F63
F64
CtrlSpace
Ctrl-a
Ctrl-b
Ctrl-c
Ctrl-d
Ctrl-e
Ctrl-f
Ctrl-g
Ctrl-h
Ctrl-i
Ctrl-j
Ctrl-k
Ctrl-l
Ctrl-m
Ctrl-n
Ctrl-o
Ctrl-p
Ctrl-q
Ctrl-r
Ctrl-s
Ctrl-t
Ctrl-u
Ctrl-v
Ctrl-w
Ctrl-x
Ctrl-y
Ctrl-z
CtrlLeftSq
CtrlBackslash
CtrlRightSq
CtrlCarat
CtrlUnderscore
Backspace
OldBackspace
Tab
Esc
Escape
Enter
```

您还可以绑定一些鼠标按钮(它们可以绑定到正常的操作或鼠标操作)

```
MouseLeft
MouseMiddle
MouseRight
MouseWheelUp
MouseWheelDown
MouseWheelLeft
MouseWheelRight
```

## 键序列

可以通过在括号中一个接一个地指定有效键来绑定键序列,例如
 `<Ctrl-x><Ctrl-c>`.

# 默认按键绑定配置.

MacOS上有少数按键绑定与其他按键绑定不同
操作系统.这是因为不同的操作系统有不同的
文本编辑默认的约定.

```json
{
    "Up":             "CursorUp",
    "Down":           "CursorDown",
    "Right":          "CursorRight",
    "Left":           "CursorLeft",
    "ShiftUp":        "SelectUp",
    "ShiftDown":      "SelectDown",
    "ShiftLeft":      "SelectLeft",
    "ShiftRight":     "SelectRight",
    "AltLeft":        "WordLeft", (Mac)
    "AltRight":       "WordRight", (Mac)
    "AltUp":          "MoveLinesUp",
    "AltDown":        "MoveLinesDown",
    "CtrlShiftRight": "SelectWordRight",
    "CtrlShiftLeft":  "SelectWordLeft",
    "AltLeft":        "StartOfTextToggle",
    "AltRight":       "EndOfLine",
    "AltShiftRight":  "SelectWordRight", (Mac)
    "AltShiftLeft":   "SelectWordLeft", (Mac)
    "CtrlLeft":       "StartOfText", (Mac)
    "CtrlRight":      "EndOfLine", (Mac)
    "AltShiftLeft":   "SelectToStartOfTextToggle",
    "CtrlShiftLeft":  "SelectToStartOfTextToggle", (Mac)
    "ShiftHome":      "SelectToStartOfTextToggle",
    "AltShiftRight":  "SelectToEndOfLine",
    "CtrlShiftRight": "SelectToEndOfLine", (Mac)
    "ShiftEnd":       "SelectToEndOfLine",
    "CtrlUp":         "CursorStart",
    "CtrlDown":       "CursorEnd",
    "CtrlShiftUp":    "SelectToStart",
    "CtrlShiftDown":  "SelectToEnd",
    "Alt-{":          "ParagraphPrevious",
    "Alt-}":          "ParagraphNext",
    "Enter":          "InsertNewline",
    "Ctrl-h":          "Backspace",
    "Backspace":      "Backspace",
    "Alt-CtrlH":      "DeleteWordLeft",
    "Alt-Backspace":  "DeleteWordLeft",
    "Tab":            "Autocomplete|IndentSelection|InsertTab",
    "Backtab":        "OutdentSelection|OutdentLine",
    "Ctrl-o":          "OpenFile",
    "Ctrl-s":          "Save",
    "Ctrl-f":          "Find",
    "Alt-F":           "FindLiteral",
    "Ctrl-n":          "FindNext",
    "Ctrl-p":          "FindPrevious",
    "Ctrl-z":          "Undo",
    "Ctrl-y":          "Redo",
    "Ctrl-c":          "CopyLine|Copy",
    "Ctrl-x":          "Cut",
    "Ctrl-k":          "CutLine",
    "Ctrl-d":          "DuplicateLine",
    "Ctrl-v":          "Paste",
    "Ctrl-a":          "SelectAll",
    "Ctrl-t":          "AddTab",
    "Alt-,":           "PreviousTab",
    "Alt-.":           "NextTab",
    "Home":           "StartOfText",
    "End":            "EndOfLine",
    "CtrlHome":       "CursorStart",
    "CtrlEnd":        "CursorEnd",
    "PageUp":         "CursorPageUp",
    "PageDown":       "CursorPageDown",
    "CtrlPageUp":     "PreviousTab",
    "CtrlPageDown":   "NextTab",
    "Ctrl-g":          "ToggleHelp",
    "Alt-g":          "ToggleKeyMenu",
    "Ctrl-r":          "ToggleRuler",
    "Ctrl-l":          "command-edit:goto ",
    "Delete":         "Delete",
    "Ctrl-b":          "ShellMode",
    "Ctrl-q":          "Quit",
    "Ctrl-e":          "CommandMode",
    "Ctrl-w":          "NextSplit",
    "Ctrl-u":          "ToggleMacro",
    "Ctrl-j":          "PlayMacro",
    "Insert":         "ToggleOverwriteMode",

    // Emacs-style keybindings
    "Alt-f": "WordRight",
    "Alt-b": "WordLeft",
    "Alt-a": "StartOfLine",
    "Alt-e": "EndOfLine",

    // Integration with file managers
    "F2":  "Save",
    "F3":  "Find",
    "F4":  "Quit",
    "F7":  "Find",
    "F10": "Quit",
    "Esc": "Escape",

    // Mouse bindings
    "MouseWheelUp":   "ScrollUp",
    "MouseWheelDown": "ScrollDown",
    "MouseLeft":      "MousePress",
    "MouseMiddle":    "PastePrimary",
    "Ctrl-MouseLeft": "MouseMultiCursor",

    "Alt-n":        "SpawnMultiCursor",
    "AltShiftUp":   "SpawnMultiCursorUp",
    "AltShiftDown": "SpawnMultiCursorDown",
    "Alt-m":        "SpawnMultiCursorSelect",
    "Alt-p":        "RemoveMultiCursor",
    "Alt-c":        "RemoveAllMultiCursors",
    "Alt-x":        "SkipMultiCursor",
}
```

## 窗口类型绑定

也可以为不同的窗口类型指定键绑定.例如,
创建一个只影响命令栏的绑定,使用 `command` 子组:

```
{
    "command": {
        "Ctrl-w": "WordLeft"
    }
}
```

可能的窗口类型是 `buffer` (普通缓冲区),`command` (命令栏),
以及 `terminal` (终端窗格).命令和终端窗口的默认值
如下所示:

```
{
    "terminal": {
        "<Ctrl-q><Ctrl-q>": "Exit",
        "<Ctrl-e><Ctrl-e>": "CommandMode",
        "<Ctrl-w><Ctrl-w>": "NextSplit"
    },

    "command": {
        "Up":             "HistoryUp",
        "Down":           "HistoryDown",
        "Right":          "CursorRight",
        "Left":           "CursorLeft",
        "ShiftUp":        "SelectUp",
        "ShiftDown":      "SelectDown",
        "ShiftLeft":      "SelectLeft",
        "ShiftRight":     "SelectRight",
        "AltLeft":        "StartOfTextToggle",
        "AltRight":       "EndOfLine",
        "AltUp":          "CursorStart",
        "AltDown":        "CursorEnd",
        "AltShiftRight":  "SelectWordRight",
        "AltShiftLeft":   "SelectWordLeft",
        "CtrlLeft":       "WordLeft",
        "CtrlRight":      "WordRight",
        "CtrlShiftLeft":  "SelectToStartOfTextToggle",
        "ShiftHome":      "SelectToStartOfTextToggle",
        "CtrlShiftRight": "SelectToEndOfLine",
        "ShiftEnd":       "SelectToEndOfLine",
        "CtrlUp":         "CursorStart",
        "CtrlDown":       "CursorEnd",
        "CtrlShiftUp":    "SelectToStart",
        "CtrlShiftDown":  "SelectToEnd",
        "Enter":          "ExecuteCommand",
        "CtrlH":          "Backspace",
        "Backspace":      "Backspace",
        "OldBackspace":   "Backspace",
        "Alt-CtrlH":      "DeleteWordLeft",
        "Alt-Backspace":  "DeleteWordLeft",
        "Tab":            "CommandComplete",
        "Backtab":        "CycleAutocompleteBack",
        "Ctrl-z":         "Undo",
        "Ctrl-y":         "Redo",
        "Ctrl-c":         "CopyLine|Copy",
        "Ctrl-x":         "Cut",
        "Ctrl-k":         "CutLine",
        "Ctrl-v":         "Paste",
        "Home":           "StartOfTextToggle",
        "End":            "EndOfLine",
        "CtrlHome":       "CursorStart",
        "CtrlEnd":        "CursorEnd",
        "Delete":         "Delete",
        "Ctrl-q":         "AbortCommand",
        "Ctrl-e":         "EndOfLine",
        "Ctrl-a":         "StartOfLine",
        "Ctrl-w":         "DeleteWordLeft",
        "Insert":         "ToggleOverwriteMode",
        "Ctrl-b":         "WordLeft",
        "Ctrl-f":         "WordRight",
        "Ctrl-d":         "DeleteWordLeft",
        "Ctrl-m":         "ExecuteCommand",
        "Ctrl-n":         "HistoryDown",
        "Ctrl-p":         "HistoryUp",
        "Ctrl-u":         "SelectToStart",

        // Emacs-style keybindings
        "Alt-f": "WordRight",
        "Alt-b": "WordLeft",
        "Alt-a": "StartOfText",
        "Alt-e": "EndOfLine",

        // Integration with file managers
        "F10": "AbortCommand",
        "Esc": "AbortCommand",

        // Mouse bindings
        "MouseWheelUp":   "HistoryUp",
        "MouseWheelDown": "HistoryDown",
        "MouseLeft":      "MousePress",
        "MouseMiddle":    "PastePrimary"
    }
}
```

## 最后指出

注意:在一些旧的终端模拟器和Windows机器上,`Ctrl-h` 应该是
用于退格.

此外,alt键可以使用 `Alt-key` 绑定.例如 `Alt-a`或 `Alt-Up`.
Micro支持在 `Alt` 和 `-` 等修饰符之间使用可选的比如`Ctrl` 所以
 `Alt-a` 可以重写为 `Alta` (大小写对alt绑定很重要).这就是为什么
 在默认的键绑定,你可以看到 `AltShiftLeft` 而不是`Alt-ShiftLeft` 
 (它们是等效的).

请注意,终端模拟器是奇怪的应用程序和微型的
接收终端决定发送的关键事件.一些终端模拟器
可能不会发送某些事件,即使本文档说微可以收到
事件.当你按下a键时,可以看到micro从终端接收到什么