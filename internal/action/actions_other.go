// +build plan9 nacl windows

package action

func (*BufPane) Suspend() bool {
	InfoBar.Error("挂起仅在BSD/Linux上受支持")
	return false
}
