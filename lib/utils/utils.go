package utils

// ToCmdLine convert strings to [][]byte
func ToCmdLine(cmds ...string) [][]byte {
	args := make([][]byte, len(cmds))
	for i, s := range cmds {
		args[i] = []byte(s)
	}
	return args
}

// ToCmdLineWithCmdName cmdName and []byte-type argument to CmdLine
func ToCmdLineWithCmdName(cmdName string, args ...[]byte) [][]byte {
	res := make([][]byte, len(args)+1)
	res[0] = []byte(cmdName)
	for i, s := range args {
		res[i+1] = []byte(s)
	}
	return res
}
