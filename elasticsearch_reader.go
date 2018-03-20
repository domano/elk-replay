package main

import (
	"context"
	"io"

	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

type ESReadFunc func(ReadArgs) error

type ReadArgs struct {
	FilterPairs [][2]string
	Indices     []string
}

func NewReader(ctx context.Context, client *elastic.Client) (readerFunc ESReadFunc, out <-chan LogEntry, err error) {
	out = make(chan LogEntry)

	readerFunc = func(args ReadArgs) error {
		if len(args.Indices) == 0 {
			args.Indices, err = client.IndexNames()
			if err != nil {
				errors.Wrap(err, "Could not fetch indices: \n")
			}
		}
		query := elastic.NewBoolQuery()
		for _, fp := range args.FilterPairs {
			query.Must(elastic.NewMatchQuery(fp[0], fp[1]))
		}

		scrollServ := client.Scroll(args.Indices...)
		for {
			result, err := scrollServ.Do(ctx)
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return errors.Wrap(err, "Could not scroll on indices: \n")
			}
			result.Hits.Hits
		}

	}
	return readerFunc, out, nil
}
