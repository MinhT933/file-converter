package template 


func RunImportWorkFlow[T any](ctx context.Context, job ImportJob[T], path string, batchSize int) <-chan ImportResult[T] {
	// Implement the workflow logic here
	return nil
}