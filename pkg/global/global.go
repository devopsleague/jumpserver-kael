package global

import (
	"github.com/jumpserver/kael/pkg/config"
	"go.uber.org/zap"
)

var (
	Config *config.Config
	Logger *zap.Logger
)
