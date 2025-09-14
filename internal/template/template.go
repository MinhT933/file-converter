package template

import "context"

func RunImportWorkflow[T any](ctx context.Context, job ImportJob[T], path string, batchSize int) <-chan ImportResult[T] {
	// Implement the workflow logic here
	results := make(chan ImportResult[T])
	go func() {
		defer close(results)
		// Simulate processing

		//errs là gì ?
        rows, errs := job.Parse(ctx, path)
		buffer := make([]T, 0, batchSize)

		flush := func() {
			if len(buffer) == 0 {
				return
			}
			if err := job.InsertBatch(ctx, buffer); err != nil {
				for _, d := range buffer {
					results <- ImportResult[T]{Data: &d, Errs: []error{err}, Stage: "insert"}
				}
			} else {
				for _, d := range buffer {
					results <- ImportResult[T]{Data: &d, Stage: "success"}
				}
			}
			buffer = buffer[:0]
		}

		for row := range rows {
			data, err := job.Transform(row)
			if err != nil {
				results <- ImportResult[T]{Row: row, Errs: []error{err}, Stage: "transform"}
				job.ReportError(row, []error{err})
				continue
			}
			// Validate
			if vErrs := job.Validate(data); len(vErrs) > 0 {
				results <- ImportResult[T]{Row: row, Data: &data, Errs: vErrs, Stage: "validate"}
				job.ReportError(row, vErrs)
				continue
			}
			buffer = append(buffer, data)
			if len(buffer) >= batchSize {
				flush()
			}
			

			// nối batch cuối
			flush()

			if err := <- errs; err != nil {
				results <- ImportResult[T]{Errs: []error{err}, Stage: "parse"}
			}
		}
	}()
	return results
}