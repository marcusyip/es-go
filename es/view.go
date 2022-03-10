package es

type View interface {
	GetTableName() string
}

type BaseView struct{}
