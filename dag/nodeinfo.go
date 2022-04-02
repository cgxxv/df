package dag

import (
	"context"
	"os"
)

type Nodeinfo struct {
	UUID        string   `json:"uuid"`
	TaskName    string   `json:"task_name"`
	Parents     []string `json:"parent"`
	Children    []string `json:"children"`
	Application string   `json:"application"`
	Service     string   `json:"service"`
}

func getSchedulerNodeinfo(ctx context.Context, service string) []*Nodeinfo {
	/*
		1 -- 2 -- 5
		|         |
		3 --------4
	*/
	// var nms = []*Nodeinfo{
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
	var nms = []*Nodeinfo{
		{"1", "1", []string{}, []string{"2", "3"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"2", "2", []string{"1"}, []string{"5"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"3", "3", []string{"1"}, []string{"4"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"4", "4", []string{"3", "5"}, []string{}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
		{"5", "5", []string{"2"}, []string{"4"}, os.Getenv("APPLICATION"), os.Getenv("SERVICE")},
	}
	return nms
}

//debug
func GetSchedulerNodeinfo(ctx context.Context, service string) []*Nodeinfo {
	return getSchedulerNodeinfo(ctx, service)
}
