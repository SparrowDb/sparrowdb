package script

var supportedTypes = []func(buf []byte) bool{
	func(buf []byte) bool {
		return len(buf) > 3 && buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47
	},
	func(buf []byte) bool {
		return len(buf) > 2 && buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF
	},
}

// IsSupportedFileType check if file is supported by sparrowdb
// image manipulation library
func IsSupportedFileType(buf []byte) bool {
	for _, f := range supportedTypes {
		if f(buf) == true {
			return true
		}
	}
	return false
}
