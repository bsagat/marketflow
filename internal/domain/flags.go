package domain

import "flag"

var Port = flag.String("port", "8080", "Default server port number")
var HelpFlag = flag.Bool("help", false, "Show help message")
