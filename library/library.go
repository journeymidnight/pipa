package library

import "github.com/journeymidnight/pipa/imagick"

type Library interface {

}

func NewLibrary() Library {
	//TODO: support other libraries
	return imagick.Initialize()
}