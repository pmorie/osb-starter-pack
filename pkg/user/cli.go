package user

import (
	"flag"
)

// Options holds the options specified by the user's code on the command
// line. Users should add their own options here and add flags for them in
// AddUserFlags.
type Options struct {
	CatalogPath string
	Async       bool
}

// AddUserFlags is a hook called to initialize the CLI flags for user options it
// is called after the flags are added for the skeleton and before flag.Parse is
// called.
func AddUserFlags(o *Options) {
	flag.StringVar(&o.CatalogPath, "catalogPath", "", "The path to the catalog")
	flag.BoolVar(&o.Async, "async", false, "Indicates whether the broker is handling the requests asynchronously.")
}
