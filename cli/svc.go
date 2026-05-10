package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type Version struct {
	ProjectId int
	Ref       string
	Type      string
	Version   string
	Stop      func()
}

func GetVersions(projectsIds []int) <-chan Version {
	var cversion <-chan Version
	{
		cv := make(chan Version)
		cversion = cv
		wg := &sync.WaitGroup{}
		go func() {
			defer close(cv)
			for _, id := range projectsIds {
				wg.Go(func() {
					perpage := 100
					ctx, ctxCanc := context.WithCancel(context.Background())
					defer ctxCanc()
					for page := 1; ; page++ {
						select {
						case <-ctx.Done():
							return
						default:
						}
						items, err := glclient.Packages(id, page, perpage)
						if err != nil {
							fmt.Fprintf(os.Stderr, "unable to retrieve projectId=%d", id)
							return
						}
						if len(items) == 0 {
							return
						}
						for _, it := range items {
							cv <- Version{
								ProjectId: id,
								Ref:       it.Pipeline.Ref,
								Type:      it.PackageType,
								Version:   it.Version,
								Stop:      ctxCanc,
							}
						}
					}
				})
			}
			wg.Wait()
		}()
	}
	return cversion
}
