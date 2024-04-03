package goappbase

import (
	"log"

	"github.com/mitoteam/mttools"
)

func Do() {
	log.Println("It works! " + mttools.RandomString(5))
}
