package config

import (
	"os"
)

var DEBUG = (os.Getenv("DEBUG") == "true")
