package responses

type FPDFPage_FlattenResult int

const (
	FPDFPage_FlattenResultFail        FPDFPage_FlattenResult = 0 // Flatten operation failed.
	FPDFPage_FlattenResultSuccess     FPDFPage_FlattenResult = 1 // Flatten operation succeed.
	FPDFPage_FlattenResultNothingToDo FPDFPage_FlattenResult = 2 // There is nothing can be flatten.
)

type FPDFPage_Flatten struct {
	Page   int                    // The page number (0-index based).
	Result FPDFPage_FlattenResult // The result of the flatten.
}
