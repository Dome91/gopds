package configuration

import "flag"

var Address string
var Port string

func ParseFlags() {
	flag.StringVar(&Address, "address", "", "-address=192.168.0.1")
	flag.StringVar(&Port, "port", "3000", "-port=4000")
	flag.Parse()
}
