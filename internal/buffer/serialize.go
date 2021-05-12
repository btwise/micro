package buffer

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/encoding"

	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/util"
)

// The SerializedBuffer holds the types that get serialized when a buffer is saved
// These are used for the savecursor and saveundo options
type SerializedBuffer struct {
	EventHandler *EventHandler
	Cursor       Loc
	ModTime      time.Time
}

// Serialize serializes the buffer to config.ConfigDir/buffers
func (b *Buffer) Serialize() error {
	if !b.Settings["savecursor"].(bool) && !b.Settings["saveundo"].(bool) {
		return nil
	}
	if b.Path == "" {
		return nil
	}

	name := filepath.Join(config.ConfigDir, "buffers", util.EscapePath(b.AbsPath))

	return overwriteFile(name, encoding.Nop, func(file io.Writer) error {
		err := gob.NewEncoder(file).Encode(SerializedBuffer{
			b.EventHandler,
			b.GetActiveCursor().Loc,
			b.ModTime,
		})
		return err
	}, false)
}

// Unserialize loads the buffer info from config.ConfigDir/buffers
func (b *Buffer) Unserialize() error {
	// If either savecursor or saveundo is turned on, we need to load the serialized information
	// from ~/.config/micro/buffers
	if b.Path == "" {
		return nil
	}
	file, err := os.Open(filepath.Join(config.ConfigDir, "buffers", util.EscapePath(b.AbsPath)))
	if err == nil {
		defer file.Close()
		var buffer SerializedBuffer
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&buffer)
		if err != nil {
			return errors.New(err.Error() + "\n您可能要删除~/.config/micro/buffers中的文件 (这些文件\n存储'saveundo'和'savecursor'选项的信息) 如果\n这个问题仍然存在.\n这可能是由于升级到2.0版并删除了'buffers'引起的\n目录将重置光标并撤消历史记录并解决问题.")
		}
		if b.Settings["savecursor"].(bool) {
			b.StartCursor = buffer.Cursor
		}

		if b.Settings["saveundo"].(bool) {
			// We should only use last time's eventhandler if the file wasn't modified by someone else in the meantime
			if b.ModTime == buffer.ModTime {
				b.EventHandler = buffer.EventHandler
				b.EventHandler.cursors = b.cursors
				b.EventHandler.buf = b.SharedBuffer
			}
		}
	}
	return nil
}
