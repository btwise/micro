package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	isatty "github.com/mattn/go-isatty"
	lua "github.com/yuin/gopher-lua"
	"github.com/zyedidia/micro/v2/internal/action"
	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/clipboard"
	"github.com/zyedidia/micro/v2/internal/config"
	ulua "github.com/zyedidia/micro/v2/internal/lua"
	"github.com/zyedidia/micro/v2/internal/screen"
	"github.com/zyedidia/micro/v2/internal/shell"
	"github.com/zyedidia/micro/v2/internal/util"
	"github.com/zyedidia/tcell/v2"
)

var (
	// Event channel
	autosave chan bool

	// Command line flags
	flagVersion   = flag.Bool("version", false, "显示版本号和信息")
	flagConfigDir = flag.String("config-dir", "", "指定配置目录的自定义位置")
	flagOptions   = flag.Bool("options", false, "显示所有选项帮助")
	flagDebug     = flag.Bool("debug", false, "启用调试模式 (将调试信息打印到 ./log.txt)")
	flagPlugin    = flag.String("plugin", "", "插件命令")
	flagClean     = flag.Bool("clean", false, "清理配置目录")
	optionFlags   map[string]*string

	sigterm chan os.Signal
	sighup  chan os.Signal
)

func InitFlags() {
	flag.Usage = func() {
		fmt.Println("用法: micro [选项] [文件]...")
		fmt.Println("-clean")
		fmt.Println("    \t清理配置目录")
		fmt.Println("-config-dir dir")
		fmt.Println("    \t指定配置目录的自定义位置")
		fmt.Println("[FILE]:LINE:COL (如果启用了`parsecursor`选项)")
		fmt.Println("+LINE:COL")
		fmt.Println("    \t指定一行和一列以在打开缓冲区时开始游标")
		fmt.Println("-options")
		fmt.Println("    \t显示所有选项帮助")
		fmt.Println("-debug")
		fmt.Println("    \t启用调试模式 (将调试信息记录到 ./log.txt)")
		fmt.Println("-version")
		fmt.Println("    \t显示版本号和信息")

		fmt.Print("\nMicro的插件可以使用以下命令在命令行中进行管理.\n")
		fmt.Println("-plugin install [PLUGIN]...")
		fmt.Println("    \t安装插件")
		fmt.Println("-plugin remove [PLUGIN]...")
		fmt.Println("    \t删除插件")
		fmt.Println("-plugin update [PLUGIN]...")
		fmt.Println("    \t更新插件 (如果未提供任何参数，则更新所有插件)")
		fmt.Println("-plugin search [PLUGIN]...")
		fmt.Println("    \t搜索插件")
		fmt.Println("-plugin list")
		fmt.Println("    \t列出已安装的插件")
		fmt.Println("-plugin available")
		fmt.Println("    \t列出可用的插件")

		fmt.Print("\nMicro的选项也可以通过命令行参数进行设置，以实现快速调整\n. 实际配置, 请使用settings.json文件\n(参看 'help options').\n\n")
		fmt.Println("-option value")
		fmt.Println("    \t在此会话中将`option`设置为`value`")
		fmt.Println("    \t例如: `micro -syntax off file.c`")
		fmt.Println("\n使用 `micro -options` 查看配置选项的完整列表")
	}

	optionFlags = make(map[string]*string)

	for k, v := range config.DefaultAllSettings() {
		optionFlags[k] = flag.String(k, "", fmt.Sprintf("%s 选项. 默认值: '%v'.", k, v))
	}

	flag.Parse()

	if *flagVersion {
		// If -version was passed
		fmt.Println("版本:", util.Version)
		fmt.Println("提交哈希:", util.CommitHash)
		fmt.Println("编译于", util.CompileDate)
		os.Exit(0)
	}

	if *flagOptions {
		// If -options was passed
		var keys []string
		m := config.DefaultAllSettings()
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := m[k]
			fmt.Printf("-%s 值\n", k)
			fmt.Printf("    \t默认值: '%v'\n", v)
		}
		os.Exit(0)
	}

	if util.Debug == "关" && *flagDebug {
		util.Debug = "开"
	}
}

// DoPluginFlags parses and executes any flags that require LoadAllPlugins (-plugin and -clean)
func DoPluginFlags() {
	if *flagClean || *flagPlugin != "" {
		config.LoadAllPlugins()

		if *flagPlugin != "" {
			args := flag.Args()

			config.PluginCommand(os.Stdout, *flagPlugin, args)
		} else if *flagClean {
			CleanConfig()
		}

		os.Exit(0)
	}
}

// LoadInput determines which files should be loaded into buffers
// based on the input stored in flag.Args()
func LoadInput(args []string) []*buffer.Buffer {
	// There are a number of ways micro should start given its input

	// 1. If it is given a files in flag.Args(), it should open those

	// 2. If there is no input file and the input is not a terminal, that means
	// something is being piped in and the stdin should be opened in an
	// empty buffer

	// 3. If there is no input file and the input is a terminal, an empty buffer
	// should be opened

	var filename string
	var input []byte
	var err error
	buffers := make([]*buffer.Buffer, 0, len(args))

	btype := buffer.BTDefault
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		btype = buffer.BTStdout
	}

	files := make([]string, 0, len(args))
	flagStartPos := buffer.Loc{-1, -1}
	flagr := regexp.MustCompile(`^\+(\d+)(?::(\d+))?$`)
	for _, a := range args {
		match := flagr.FindStringSubmatch(a)
		if len(match) == 3 && match[2] != "" {
			line, err := strconv.Atoi(match[1])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			col, err := strconv.Atoi(match[2])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			flagStartPos = buffer.Loc{col - 1, line - 1}
		} else if len(match) == 3 && match[2] == "" {
			line, err := strconv.Atoi(match[1])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			flagStartPos = buffer.Loc{0, line - 1}
		} else {
			files = append(files, a)
		}
	}

	if len(files) > 0 {
		// Option 1
		// We go through each file and load it
		for i := 0; i < len(files); i++ {
			buf, err := buffer.NewBufferFromFileAtLoc(files[i], btype, flagStartPos)
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			// If the file didn't exist, input will be empty, and we'll open an empty buffer
			buffers = append(buffers, buf)
		}
	} else if !isatty.IsTerminal(os.Stdin.Fd()) {
		// Option 2
		// The input is not a terminal, so something is being piped in
		// and we should read from stdin
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			screen.TermMessage("从标准输入读取错误: ", err)
			input = []byte{}
		}
		buffers = append(buffers, buffer.NewBufferFromStringAtLoc(string(input), filename, btype, flagStartPos))
	} else {
		// Option 3, just open an empty buffer
		buffers = append(buffers, buffer.NewBufferFromStringAtLoc(string(input), filename, btype, flagStartPos))
	}

	return buffers
}

func main() {
	defer func() {
		if util.Stdout.Len() > 0 {
			fmt.Fprint(os.Stdout, util.Stdout.String())
		}
		os.Exit(0)
	}()

	// runtime.SetCPUProfileRate(400)
	// f, _ := os.Create("micro.prof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	var err error

	InitFlags()

	InitLog()

	err = config.InitConfigDir(*flagConfigDir)
	if err != nil {
		screen.TermMessage(err)
	}

	config.InitRuntimeFiles()
	err = config.ReadSettings()
	if err != nil {
		screen.TermMessage(err)
	}
	err = config.InitGlobalSettings()
	if err != nil {
		screen.TermMessage(err)
	}

	// flag options
	for k, v := range optionFlags {
		if *v != "" {
			nativeValue, err := config.GetNativeValue(k, config.DefaultAllSettings()[k], *v)
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			config.GlobalSettings[k] = nativeValue
		}
	}

	DoPluginFlags()

	err = screen.Init()
	if err != nil {
		fmt.Println(err)
		fmt.Println("致命：Micro无法初始化屏幕.")
		os.Exit(1)
	}

	sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sighup, syscall.SIGHUP)

	m := clipboard.SetMethod(config.GetGlobalOption("clipboard").(string))
	clipErr := clipboard.Initialize(m)

	defer func() {
		if err := recover(); err != nil {
			if screen.Screen != nil {
				screen.Screen.Fini()
			}
			if e, ok := err.(*lua.ApiError); ok {
				fmt.Println("Lua API错误:", e)
			} else {
				fmt.Println("Micro 遇到错误:", errors.Wrap(err, 2).ErrorStack(), "\n如果您可以重现此错误, 请在https://github.com/zyedidia/micro/issues报告")
			}
			// backup all open buffers
			for _, b := range buffer.OpenBuffers {
				b.Backup()
			}
			os.Exit(1)
		}
	}()

	err = config.LoadAllPlugins()
	if err != nil {
		screen.TermMessage(err)
	}

	action.InitBindings()
	action.InitCommands()

	err = config.InitColorscheme()
	if err != nil {
		screen.TermMessage(err)
	}

	err = config.RunPluginFn("preinit")
	if err != nil {
		screen.TermMessage(err)
	}

	args := flag.Args()
	b := LoadInput(args)

	if len(b) == 0 {
		// No buffers to open
		screen.Screen.Fini()
		runtime.Goexit()
	}

	action.InitTabs(b)
	action.InitGlobals()

	err = config.RunPluginFn("init")
	if err != nil {
		screen.TermMessage(err)
	}

	err = config.RunPluginFn("postinit")
	if err != nil {
		screen.TermMessage(err)
	}

	if clipErr != nil {
		log.Println(clipErr, " 或更改 'clipboard' option")
	}

	if a := config.GetGlobalOption("autosave").(float64); a > 0 {
		config.SetAutoTime(int(a))
		config.StartAutoSave()
	}

	screen.Events = make(chan tcell.Event)

	// Here is the event loop which runs in a separate thread
	go func() {
		for {
			screen.Lock()
			e := screen.Screen.PollEvent()
			screen.Unlock()
			if e != nil {
				screen.Events <- e
			}
		}
	}()

	// clear the drawchan so we don't redraw excessively
	// if someone requested a redraw before we started displaying
	for len(screen.DrawChan()) > 0 {
		<-screen.DrawChan()
	}

	// wait for initial resize event
	select {
	case event := <-screen.Events:
		action.Tabs.HandleEvent(event)
	case <-time.After(10 * time.Millisecond):
		// time out after 10ms
	}

	for {
		DoEvent()
	}
}

// DoEvent runs the main action loop of the editor
func DoEvent() {
	var event tcell.Event

	// Display everything
	screen.Screen.Fill(' ', config.DefStyle)
	screen.Screen.HideCursor()
	action.Tabs.Display()
	for _, ep := range action.MainTab().Panes {
		ep.Display()
	}
	action.MainTab().Display()
	action.InfoBar.Display()
	screen.Screen.Show()

	// Check for new events
	select {
	case f := <-shell.Jobs:
		// If a new job has finished while running in the background we should execute the callback
		ulua.Lock.Lock()
		f.Function(f.Output, f.Args)
		ulua.Lock.Unlock()
	case <-config.Autosave:
		ulua.Lock.Lock()
		for _, b := range buffer.OpenBuffers {
			b.Save()
		}
		ulua.Lock.Unlock()
	case <-shell.CloseTerms:
	case event = <-screen.Events:
	case <-screen.DrawChan():
		for len(screen.DrawChan()) > 0 {
			<-screen.DrawChan()
		}
	case <-sighup:
		for _, b := range buffer.OpenBuffers {
			if !b.Modified() {
				b.Fini()
			}
		}
		os.Exit(0)
	case <-sigterm:
		for _, b := range buffer.OpenBuffers {
			if !b.Modified() {
				b.Fini()
			}
		}

		if screen.Screen != nil {
			screen.Screen.Fini()
		}
		os.Exit(0)
	}

	ulua.Lock.Lock()
	// if event != nil {
	if action.InfoBar.HasPrompt {
		action.InfoBar.HandleEvent(event)
	} else {
		action.Tabs.HandleEvent(event)
	}
	// }
	ulua.Lock.Unlock()
}
