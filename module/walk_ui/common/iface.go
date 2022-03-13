package common

import "github.com/lxn/walk"

type Page interface {
	// Provided by Walk
	walk.Container
	Parent() walk.Container
	SetParent(parent walk.Container) error
}
