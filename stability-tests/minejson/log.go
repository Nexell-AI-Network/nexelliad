package main

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/logger"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/panics"
)

var (
	backendLog = logger.NewBackend()
	log        = backendLog.Logger("MNJS")
	spawn      = panics.GoroutineWrapperFunc(log)
)
