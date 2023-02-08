package buffer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/screen"
	"github.com/zyedidia/micro/v2/internal/util"
	"golang.org/x/text/encoding"
)

const backupMsg = `检测到该文件的备份. 这可能意味着Micro
编辑此文件时崩溃, 或当前有另一个micro实例编辑此文件.

备份是在%s上创建的, 文件是

%s

* 'recover' 将备份作为未保存的更改应用于当前缓冲区.
  关闭缓冲区后,备份将被删除.
* 'ignore' 将忽略备份,放弃其更改. 备份文件
  将被删除.
* 'abort' 将中止打开的操作, 而是打开一个空缓冲区.

选项: [r]ecover, [i]gnore, [a]bort: `

var backupRequestChan chan *Buffer

func backupThread() {
	for {
		time.Sleep(time.Second * 8)

		for len(backupRequestChan) > 0 {
			b := <-backupRequestChan
			bfini := atomic.LoadInt32(&(b.fini)) != 0
			if !bfini {
				b.Backup()
			}
		}
	}
}

func init() {
	backupRequestChan = make(chan *Buffer, 10)

	go backupThread()
}

func (b *Buffer) RequestBackup() {
	if !b.requestedBackup {
		select {
		case backupRequestChan <- b:
		default:
			// channel is full
		}
		b.requestedBackup = true
	}
}

// Backup saves the current buffer to ConfigDir/backups
func (b *Buffer) Backup() error {
	if !b.Settings["backup"].(bool) || b.Path == "" || b.Type != BTDefault {
		return nil
	}

	backupdir, err := util.ReplaceHome(b.Settings["backupdir"].(string))
	if backupdir == "" || err != nil {
		backupdir = filepath.Join(config.ConfigDir, "backups")
	}
	if _, err := os.Stat(backupdir); os.IsNotExist(err) {
		os.Mkdir(backupdir, os.ModePerm)
	}

	name := filepath.Join(backupdir, util.EscapePath(b.AbsPath))

	err = overwriteFile(name, encoding.Nop, func(file io.Writer) (e error) {
		if len(b.lines) == 0 {
			return
		}

		// end of line
		eol := []byte{'\n'}

		// write lines
		if _, e = file.Write(b.lines[0].data); e != nil {
			return
		}

		for _, l := range b.lines[1:] {
			if _, e = file.Write(eol); e != nil {
				return
			}
			if _, e = file.Write(l.data); e != nil {
				return
			}
		}
		return
	}, false)

	b.requestedBackup = false

	return err
}

// RemoveBackup removes any backup file associated with this buffer
func (b *Buffer) RemoveBackup() {
	if !b.Settings["backup"].(bool) || b.Settings["permbackup"].(bool) || b.Path == "" || b.Type != BTDefault {
		return
	}
	f := filepath.Join(config.ConfigDir, "backups", util.EscapePath(b.AbsPath))
	os.Remove(f)
}

// ApplyBackup applies the corresponding backup file to this buffer (if one exists)
// Returns true if a backup was applied
func (b *Buffer) ApplyBackup(fsize int64) (bool, bool) {
	if b.Settings["backup"].(bool) && !b.Settings["permbackup"].(bool) && len(b.Path) > 0 && b.Type == BTDefault {
		backupfile := filepath.Join(config.ConfigDir, "backups", util.EscapePath(b.AbsPath))
		if info, err := os.Stat(backupfile); err == nil {
			backup, err := os.Open(backupfile)
			if err == nil {
				defer backup.Close()
				t := info.ModTime()
				msg := fmt.Sprintf(backupMsg, t.Format("Mon Jan _2 at 15:04, 2006"), util.EscapePath(b.AbsPath))
				choice := screen.TermPrompt(msg, []string{"r", "i", "a", "recover", "ignore", "abort"}, true)

				if choice%3 == 0 {
					// recover
					b.LineArray = NewLineArray(uint64(fsize), FFAuto, backup)
					b.isModified = true
					return true, true
				} else if choice%3 == 1 {
					// delete
					os.Remove(backupfile)
				} else if choice%3 == 2 {
					return false, false
				}
			}
		}
	}

	return false, true
}
