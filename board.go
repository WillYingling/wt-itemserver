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

type personOptions struct {
	Name      string
	HatType   string
	CapColor  string
	BrimColor string
	PomColor  string
}

type extra struct {
	Type  string
	Name  string
	Notes string
}

type snowmanBoard struct {
	Title         string
	Size          int
	PeopleOptions []personOptions
	Extras        []extra
}

func (sb *snowmanBoard) validate() error {
	if sb.Size < minBoardSize ||
		sb.Size > maxBoardSize ||
		len(sb.PeopleOptions) != sb.Size {
		return errSizeInvalid
	}

	maxTitleLength := 30
	if sb.Size == 2 {
		maxTitleLength = 15
	}

	if len(sb.Title) > maxTitleLength {
		return errTitleTooLong
	}

	/*
		for _, person := range sb.people {
			err := person.validate()
			if err != nil {
				return err
			}
		}
	*/
	return nil
}
