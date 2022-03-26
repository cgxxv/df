package dag

import (
	"errors"
)

var (
	nodeMap = make(map[string]Node, 1024)
)

type Node interface {
	GetName() string
	GetParents() []string
	GetChildren() []string
	GetApplication() string
	GetService() string
	Task
}

type NodeMeta struct {
	Name        string   `json:"name"`
	Parent      []string `json:"parent"`
	Children    []string `json:"children"`
	Application string   `json:"application"`
	Service     string   `json:"service"`
}

func RegisterNode(node Node) error {
	_, ok := nodeMap[node.GetName()]
	if ok {
		return errors.New("Duplicated node name")
	}

	nodeMap[node.GetName()] = node
	return nil
}
