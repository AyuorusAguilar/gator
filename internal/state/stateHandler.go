package state

import (
	"github.com/AyuorusAguilar/gator/internal/config"
	"github.com/AyuorusAguilar/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db *database.Queries
}