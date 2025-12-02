package macho

const LcCodeSignature LoadCommandType = 0x1d

type LoadCommandType uint32

// CodeSigningCommand is Mach-O LcCodeSignature load command.
type CodeSigningCommand struct {
	Cmd        LoadCommandType // LcCodeSignature
	Size       uint32          // sizeof this command (16)
	DataOffset uint32          // file offset of data in __LINKEDIT segment
	DataSize   uint32          // file size of data in __LINKEDIT segment
}
