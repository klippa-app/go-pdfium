package requests

type FPDFPage_FlattenUsage int

const (
	FPDFPage_FlattenUsageNormalDisplay FPDFPage_FlattenUsage = 0
	FPDFPage_FlattenUsagePrint         FPDFPage_FlattenUsage = 1
)

type FPDFPage_Flatten struct {
	Page  Page
	Usage FPDFPage_FlattenUsage // The usage flag for the flattening.
}
