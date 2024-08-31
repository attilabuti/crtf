// Based on Rich Text Format (RTF) Compression Algorithm
// https://msdn.microsoft.com/en-us/library/cc463890(v=exchg.80).aspx
package crtf

const (
	initDictSize = 207
	maxDictSize  = 4096
)

var (
	compressed   = [4]byte{'L', 'Z', 'F', 'u'}
	uncompressed = [4]byte{'M', 'E', 'L', 'A'}
)

var initDict = []byte(
	"{\\rtf1\\ansi\\mac\\deff0\\deftab720{\\fonttbl;}{\\f0\\fnil \\froman " +
		"\\fswiss \\fmodern \\fscript \\fdecor MS Sans SerifSymbolArialTimes New " +
		"RomanCourier{\\colortbl\\red0\\green0\\blue0\r\n\\par \\pard\\plain\\" +
		"f0\\fs20\\b\\i\\u\\tab\\tx",
)
