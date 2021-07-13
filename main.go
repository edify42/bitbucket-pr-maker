package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	bitbucket "github.com/ktrysmt/go-bitbucket"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// required args/flags
	owner        = kingpin.Flag("owner", "bitbucket owner").Short('o').Required().String()
	reposlug     = kingpin.Flag("repo", "git repo").Required().Short('r').String()
	sourceBranch = kingpin.Flag("sourceBranch", "source pr branch").Short('s').Required().String()
	title        = kingpin.Flag("title", "title of pr").Required().Short('t').String()

	// non-required args/flags
	closeSource       = kingpin.Flag("closeSource", "boolean switch to close source branch after merge").Default("true").Bool()
	destinationBranch = kingpin.Arg("destinationBranch", "destination pr branch (default is master)").Default("master").String()
	description       = kingpin.Flag("description", "description in the pr body").Short('d').String()
)

func main() {

	c := &bitbucket.Client{}
	kingpin.Parse()

	if os.Getenv("BB_AUTH") == "basic" {
		c = bitbucket.NewBasicAuth(os.Getenv("user"), os.Getenv("secret"))
	} else if os.Getenv("BB_AUTH") == "oauth" {
		c = bitbucket.NewOAuthClientCredentials(os.Getenv("user"), os.Getenv("secret"))
	} else {
		fmt.Println("ERROR: cannot continue without authentication option. Set environment variable BB_AUTH")
		os.Exit(1)
	}

	var reviewers []string
	r := os.Getenv("reviewers")
	if r != "" {
		reviewers = strings.Split(r, ",")
	}

	opt := &bitbucket.PullRequestsOptions{
		Owner:             *owner,
		RepoSlug:          *reposlug,
		SourceBranch:      *sourceBranch,
		DestinationBranch: *destinationBranch,
		Title:             *title,
		Description:       *description,
		CloseSourceBranch: *closeSource,
		Reviewers:         reviewers,
	}

	res, err := c.Repositories.PullRequests.Create(opt)
	if err != nil {
		fmt.Printf("%v", res)
		panic(err)
	}
	//	fmt.Println(reflect.TypeOf(res).PkgPath(), reflect.TypeOf(res).Name())
	// myMap := reflect.ValueOf(res)
	// for _, e := range myMap.MapKeys() {
	// 	thing := e["test"]
	// 	//fmt.Println(v)
	// }
	// m := make(map[string]interface{})
	// myMap := reflect.ValueOf(res)
	// for _, e := range myMap.MapKeys() {
	// 	m[e.String()] = myMap.Elem()
	// }
	// fmt.Println(m)
	byteData, _ := json.Marshal(res)
	//var t Message
	//json.Unmarshal(byteData, &t.Data)
	// fmt.Println(string(byteData))
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(byteData, &jsonMap)
	if err != nil {
		panic(err)
	}
	// fetchkey("html", jsonMap)

	a := getMap("links", jsonMap)
	b := getMap("html", a)
	d := b["href"]
	fmt.Println(d)
}

// dumpmap to stdout...is messy
func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

func getMap(key string, m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if k == key {
			thing, ok := v.(map[string]interface{})
			if !ok {
				// Can't assert, handle error.
				fmt.Println("cant assert onthing at all")
			}
			return thing
		}
	}
	return nil
}

// fetchkey to stdout...is messy
func fetchkey(space string, m map[string]interface{}) {
	for k, v := range m {
		if k != space {
			if mv, ok := v.(map[string]interface{}); ok {
				fetchkey(space, mv)
			} else {
				fmt.Printf("%v %v : %v\n", space, k, v)
			}
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}
