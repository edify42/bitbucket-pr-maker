package main

import (
	"fmt"
	"os"

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

	opt := &bitbucket.PullRequestsOptions{
		Owner:             *owner,
		RepoSlug:          *reposlug,
		SourceBranch:      *sourceBranch,
		DestinationBranch: *destinationBranch,
		Title:             *title,
		Description:       *description,
		CloseSourceBranch: *closeSource,
	}

	res, err := c.Repositories.PullRequests.Create(opt)
	if err != nil {
		fmt.Println(res)
		panic(err)
	}

}
