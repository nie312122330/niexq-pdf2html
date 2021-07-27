package reqvos

// BasePdfReq ...
type BasePdfReq struct {
	Landscape           bool    `json:"landscape,omitempty"`           // Paper orientation. Defaults to false.
	DisplayHeaderFooter bool    `json:"displayHeaderFooter,omitempty"` // Display header and footer. Defaults to false.
	PrintBackground     bool    `json:"printBackground,omitempty"`     // Print background graphics. Defaults to false.
	PaperWidth          float64 `json:"paperWidth,omitempty"`          // Paper width in inches. Defaults to 8.5 inches.
	PaperHeight         float64 `json:"paperHeight,omitempty"`         // Paper height in inches. Defaults to 11 inches.
	MarginTop           float64 `json:"marginTop"`                     // Top margin in inches. Defaults to 1cm (~0.4 inches).
	MarginBottom        float64 `json:"marginBottom"`                  // Bottom margin in inches. Defaults to 1cm (~0.4 inches).
	MarginLeft          float64 `json:"marginLeft"`                    // Left margin in inches. Defaults to 1cm (~0.4 inches).
	MarginRight         float64 `json:"marginRight"`                   // Right margin in inches. Defaults to 1cm (~0.4 inches).
	PageRanges          string  `json:"pageRanges,omitempty"`          // Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all
}

// Url2PdfReq ...
type Url2PdfReq struct {
	BasePdfReq
	HttpUrl string `json:"httpUrl" binding:"required" description:"URL路径"`
}

// HtmlPdfReq ...
type HtmlPdfReq struct {
	BasePdfReq
	HtmlTxt string `json:"htmlTxt" binding:"required" description:"Html内容"`
}
