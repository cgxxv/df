package dag

import (
	"context"
	"os"
)

type NodeMeta struct {
	Name        string   `json:"name"`
	Parents     []string `json:"parent"`
	Children    []string `json:"children"`
	Application string   `json:"application"`
	Service     string   `json:"service"`
}

func getSchedulerNodeMeta(ctx context.Context, service string) []*NodeMeta {
	/*
		1 -- 2 -- 5
		|         |
		3 --------4
	*/
	// var nms = []*NodeMeta{
	// 	{"1", []string{}, []string{"2", "3"}, zae.Service()},
	// 	{"2", []string{"1"}, []string{"5"}, zae.Service()},
	// 	{"3", []string{"1"}, []string{"4"}, zae.Service()},
	// 	{"4", []string{"3", "5"}, []string{}, zae.Service()},
	// 	{"5", []string{"2"}, []string{"4"}, zae.Service()},
	// }

	/*
		1 -- 2 -- 5
		|		 |
		3 -- 4
	*/
	var nms = []*NodeMeta{
		{"1", []string{}, []string{"2", "3"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"2", []string{"1"}, []string{"5"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"3", []string{"1"}, []string{"4"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"4", []string{"3", "5"}, []string{}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"5", []string{"2"}, []string{"4"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
	}
	return nms
}

//debug
func GetSchedulerNodeMeta(ctx context.Context, service string) []*NodeMeta {
	return getSchedulerNodeMeta(ctx, service)
}
