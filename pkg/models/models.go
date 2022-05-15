package models

type Muscle struct {
	Id         int
	Name       string
	SimpleName string
	IsFront    bool
}

type Equipment struct {
	Id   int
	Name string
}
