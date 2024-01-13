package controllers

import (
	"context"
	"fmt"
	"net/http"
	"niexq-html2pdf/config"
	"niexq-html2pdf/reqvos"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	niexqext "github.com/nie312122330/niexq-gotools"
	"github.com/nie312122330/niexq-gotools/dateext"
	"github.com/nie312122330/niexq-gotools/fileext"
	"github.com/nie312122330/niexq-gotools/logext"
	"github.com/nie312122330/niexq-gowebapi/ginext"
	"github.com/nie312122330/niexq-gowebapi/voext"
	"go.uber.org/zap"
)

var logger *zap.Logger
var chromeCtx context.Context

func init() {
	logger = logext.DefaultLogger(config.AppConf.Server.AppName)
	//初始化chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.DisableGPU, chromedp.Headless, chromedp.NoSandbox)
	if config.AppConf.ChromeConf.ExecPath != "" {
		logger.Info("配置了Chrome的路径为:" + config.AppConf.ChromeConf.ExecPath)
		opts = append(opts, chromedp.ExecPath(config.AppConf.ChromeConf.ExecPath))
	} else {
		logger.Info("使用默认安装的Chrome路径")
	}
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	chromeCtx = ctx
}

// PubCtrRegisterRouter ...
func PubCtrRegisterRouter(engine *gin.Engine, appName string) {
	contextRouter := engine.Group("/" + appName + "/pub")
	contextRouter.Use(cors())
	contextRouter.POST("/url2Pdf.do", Url2Pdf)
	contextRouter.POST("/html2Pdf.do", Html2Pdf)
}

// Url2Pdf ...
func Url2Pdf(gCtx *gin.Context) {
	var reqVo reqvos.Url2PdfReq
	ginext.ValidReq(gCtx, &reqVo)

	fileName, err := printUrl2Pdf(reqVo.HttpUrl, &reqVo.BasePdfReq)
	if err != nil {
		//写出错误
		resp := voext.NewNoBaseResp(-1, "HttpUrl转PDF失败:"+err.Error())
		gCtx.JSON(http.StatusOK, &resp)
		return
	}

	resp := voext.NewOkBaseResp(fileName)
	gCtx.JSON(http.StatusOK, &resp)

}

// Html2Pdf ...
func Html2Pdf(gCtx *gin.Context) {
	var reqVo reqvos.HtmlPdfReq
	ginext.ValidReq(gCtx, &reqVo)
	//先把内容转为文件
	htmlFileName := createFilePath("html")
	htmlFileName, _ = filepath.Abs(htmlFileName)
	fileext.WriteFileContent(htmlFileName, reqVo.HtmlTxt, false)

	fileName, err := printUrl2Pdf(fmt.Sprintf("file://"+htmlFileName), &reqVo.BasePdfReq)
	if err != nil {
		//写出错误
		resp := voext.NewNoBaseResp(-1, "HttpUrl转PDF失败:"+err.Error())
		gCtx.JSON(http.StatusOK, &resp)
		return
	}

	resp := voext.NewOkBaseResp(fileName)
	gCtx.JSON(http.StatusOK, &resp)

}

func printUrl2Pdf(urlstr string, pdfReq *reqvos.BasePdfReq) (*string, error) {
	logger.Info(fmt.Sprintf("开始转换：%s", urlstr))
	ctx, cancel := chromedp.NewContext(chromeCtx)
	defer cancel()
	// 捕获PDF
	var buf []byte
	if err := chromedp.Run(ctx, printToPdf(urlstr, pdfReq, &buf)); err != nil {
		return nil, err
	}
	fileName := createFilePath("pdf")
	fileext.WriteFile(fileName, &buf, false)

	return &fileName, nil
}

// print a specific pdf page.
func printToPdf(urlstr string, pdfReq *reqvos.BasePdfReq, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			priParams := page.PrintToPDF()
			priParams = priParams.WithLandscape(pdfReq.Landscape)
			priParams = priParams.WithDisplayHeaderFooter(pdfReq.DisplayHeaderFooter)
			priParams = priParams.WithPrintBackground(pdfReq.PrintBackground)
			priParams = priParams.WithPaperWidth(pdfReq.PaperWidth)
			priParams = priParams.WithPaperHeight(pdfReq.PaperHeight)
			priParams = priParams.WithPaperHeight(pdfReq.PaperHeight)
			priParams = priParams.WithMarginTop(pdfReq.MarginTop)
			priParams = priParams.WithMarginBottom(pdfReq.MarginBottom)
			priParams = priParams.WithMarginLeft(pdfReq.MarginLeft)
			priParams = priParams.WithMarginRight(pdfReq.MarginRight)
			priParams = priParams.WithPageRanges(pdfReq.PageRanges)

			buf, _, err := priParams.Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func createFilePath(suffix string) string {
	ymd, _ := dateext.Now().Format("yyyy-MM-dd")
	items := []string{config.AppConf.ChromeConf.Pdfdir}
	items = append(items, strings.Split(ymd, "-")...)
	items = append(items, niexqext.UUIDUperStr()+"."+suffix)
	fileName := fileext.JoinPath(items...)
	return fileName
}
