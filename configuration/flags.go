package configuration

import "flag"

var Address string
var Port string

func ParseFlags() {
	flag.StringVar(&Address, "address", "127.0.0.1", "-address=192.18.0.1")
	flag.StringVar(&Port, "port", "3000", "-port=4000")
	flag.Parse()
}
