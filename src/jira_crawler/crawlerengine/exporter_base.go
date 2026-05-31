package crawlerengine

type ExportType string

const (
	ExportCSV   ExportType = "csv"
	ExportExcel ExportType = "excel"
	ExportJSON  ExportType = "json"
)

type IssueCrawlerExporter interface {
	Export(exportType ExportType, outputPath string, output any) error
}
