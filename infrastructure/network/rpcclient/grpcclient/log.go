package grpcclient

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/logger"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/panics"
)

var log = logger.RegisterSubSystem("RPCC")
var spawn = panics.GoroutineWrapperFunc(log)
