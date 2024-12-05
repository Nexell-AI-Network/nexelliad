package blockrelay

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/logger"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/panics"
)

var log = logger.RegisterSubSystem("PROT")
var spawn = panics.GoroutineWrapperFunc(log)
