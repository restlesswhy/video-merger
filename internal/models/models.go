package models

import "github.com/pkg/errors"

type Mod string

const (
	SideBySide Mod = "sbs"
	PicInPic Mod = "pip"
)

func (m Mod) Validate() error {
	if m != SideBySide && m != PicInPic {
		return errors.New("wrong mod parameter")
	}
	return nil
}
