package constants

import "os"

// plugin related bits.
const (
	GinoPluginContractVersion             = "v2"
	GinoPluginBinaryPath                  = "/work/plugin"
	GinoPluginSockDir                     = "/tmp"
	GinoPluginSocketPath                  = "/tmp/gino.sock"
	GinoPluginSocketEnv                   = "GINO_PLUGIN_SOCKET"
	GinoPluginSocketPerms     os.FileMode = 0o600
	GinoPluginLogLevelEnv                 = "GINO_PLUGIN_LOG_LEVEL"
)
