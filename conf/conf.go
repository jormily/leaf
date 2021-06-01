package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	ServerId 		int32
	ServerType 		string

	// console
	ConsolePort   	int
	ConsolePrompt 	string = "Leaf# "
	ProfilePath   	string
	IsMaster	  	bool
	MasterAddr 		string
	RpcStypes		[]string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)
