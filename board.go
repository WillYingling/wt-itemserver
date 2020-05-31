package main

import (
	"errors"
)

const (
	minBoardSize  = 2
	maxBoardSize  = 7
	maxNameLength = 15
)

var (
	errSizeInvalid  = errors.New("Board size is invalid")
	errTitleTooLong = errors.New("Board title is too long")
	errNameTooLong  = errors.New("Name is too long")
	errHatInvalid   = errors.New("Invalid hat options")
)

type hatType int

const (
	beanie hatType = iota
	tophat
	beret
	none
)

type extraType int

const (
	dog extraType = iota
	cat
	individual
	snowMom
	snowDad
)

type snowman struct {
	name       string
	hat        hatType
	hatOptions map[string]string
}

func (s *snowman) validate() error {
	if len(s.name) > maxNameLength {
		return errNameTooLong
	}

	switch s.hat {
	case beanie:
		capColor := s.hatOptions["cap"]
		brimColor := s.hatOptions["brim"]
		//pomColor := s.hatOptions["pom"]

		if capColor == "" || brimColor == "" {
			return errHatInvalid
		}
	case tophat:
		color := s.hatOptions["color"]
		if color != "black" {
			return errHatInvalid
		}
	}
	return nil
}

type extra struct {
	name string
}

type snowmanBoard struct {
	title  string
	size   int
	people []snowman
	extras []extra
}

func (sb *snowmanBoard) validate() error {
	if sb.size < minBoardSize ||
		sb.size > maxBoardSize ||
		len(sb.people) != sb.size {
		return errSizeInvalid
	}

	maxTitleLength := 30
	if sb.size == 2 {
		maxTitleLength = 15
	}

	if len(sb.title) > maxTitleLength {
		return errTitleTooLong
	}

	for _, person := range sb.people {
		err := person.validate()
		if err != nil {
			return err
		}
	}
	return nil
}
