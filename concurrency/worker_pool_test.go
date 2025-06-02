package concurrency_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/charmingruby/pipo/lib/concurrency"
	"github.com/stretchr/testify/require"
)

func Test_WorkerPool(t *testing.T) {
	t.Run("should process messages in concurrently", func(t *testing.T) {
		processFunc := func(_ context.Context, msg dummyInput) (dummyOutput, error) {
			return dummyOutput{
				ID:     msg.ID,
				Text:   msg.RawText,
				Status: "processed",
			}, nil
		}

		wp := concurrency.NewWorkerPool(processFunc, 10)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		wp.Run(ctx)

		amountOfMessages := 10000
		messages := make([]dummyInput, amountOfMessages)

		for i := range amountOfMessages {
			messages[i] = dummyInput{
				ID:      i,
				RawText: fmt.Sprintf("message-%d", i),
			}
		}

		var wg sync.WaitGroup
		wg.Add(2)

		processedResults := make(chan dummyOutput, amountOfMessages)
		errorResults := make(chan error, amountOfMessages)

		go func() {
			defer wg.Done()
			for msg := range wp.Output() {
				processedResults <- msg
			}
			close(processedResults)
		}()

		go func() {
			defer wg.Done()
			for err := range wp.Error() {
				errorResults <- err
			}
			close(errorResults)
		}()

		err := wp.SendBatch(ctx, messages)
		require.NoError(t, err)

		err = wp.Close()
		require.NoError(t, err)

		wg.Wait()

		processedCount := 0
		for msg := range processedResults {
			require.Equal(t, "processed", msg.Status)
			processedCount++
		}

		for err := range errorResults {
			require.NoError(t, err)
		}

		require.True(t, wp.IsClosed())
		require.Equal(t, amountOfMessages, processedCount, "all messages should be processed")
	})
}

type dummyInput struct {
	RawText string
	ID      int
}

type dummyOutput struct {
	Text   string
	Status string
	ID     int
}
