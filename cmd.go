package main

import (
	"bytes"
	"fmt"
	flag "github.com/spf13/pflag"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"strings"
)

type CurlContext struct {
	help                       bool
	version                    bool
	verbose                    bool
	method                     string
	silentFail                 bool
	output                     string
	headerOutput               string
	userAgent                  string
	theUrl                     string
	ignoreBadCerts             bool
	userAuth                   string
	isSilent                   bool
	headOnly                   bool
	includeHeadersInMainOutput bool
	showErrorEvenIfSilent      bool
	referer                    string
	errorOutput                string
	cookies                    []string
	cookieJar                  string
	uploadFile                 string
	formEncoded                []string
	formMultipart              []string
	bodyContenttype            string
	_body                      io.Reader
}

func (ctx *CurlContext) SetBody(body io.Reader, mimeType string, httpMethod string) {
	ctx._body = body
	ctx.bodyContenttype = mimeType
	ctx.SetMethodIfNotSet(httpMethod)
}
func (ctx *CurlContext) SetMethodIfNotSet(httpMethod string) {
	if ctx.method == "" {
		ctx.method = httpMethod
	}
}

func main() {
	ctx := &CurlContext{}

	parseArgs(ctx)

	if ctx.theUrl == "" {
		logError(ctx, "URL was not found in command line.")
		os.Exit(8)
	}

	//request := buildRequest(ctx) - todo: requires http.Request implemented
	//client := buildClient(ctx) - todo: requires http.Client implemented
	//resp, err := client.Do(request)
	//processResponse(ctx, resp, err, request) - todo: requires http.Request & http.Response implemented
}

func parseArgs(ctx *CurlContext) {
	var empty []string
	flag.BoolVarP(&ctx.help, "help", "h", false, "Show information of all the available flags")
	flag.BoolVarP(&ctx.version, "version", "V", false, "Return version and exit")
	flag.BoolVarP(&ctx.verbose, "verbose", "v", false, "Logs all headers, and body to output")
	flag.StringVarP(&ctx.method, "method", "X", "GET", "HTTP method to use (usually GET unless otherwise modified by other parameters)")
	flag.StringVarP(&ctx.output, "output", "o", "stdout", "Where to output results")
	flag.StringVarP(&ctx.userAuth, "user", "u", "", "User:password for HTTP authentication")
	flag.StringSliceVarP(&ctx.formEncoded, "data", "d", empty, "HTML form data, set mime type to 'application/x-www-form-urlencoded'")
	flag.StringSliceVarP(&ctx.formMultipart, "form", "F", empty, "HTML form data, set mime type to 'multipart/form-data'")
	flag.StringVar(&ctx.errorOutput, "stderr", "stderr", "Log errors to this replacement for stderr")
	flag.StringVarP(&ctx.headerOutput, "dump-header", "D", "", "Where to output headers (not on by default)")
	flag.StringVarP(&ctx.userAgent, "user-agent", "A", "go-curling/##DEV##", "User-agent to use")
	flag.StringVarP(&ctx.referer, "referer", "e", "", "Referer URL to use with HTTP request")
	flag.StringVar(&ctx.theUrl, "url", "", "Requesting URL")
	flag.BoolVarP(&ctx.silentFail, "fail", "f", false, "If fail do not emit contents just return fail exit code (-6)")
	flag.BoolVarP(&ctx.ignoreBadCerts, "insecure", "k", false, "Ignore invalid SSL certificates")
	flag.BoolVarP(&ctx.isSilent, "silent", "s", false, "Silence all program console output")
	flag.BoolVarP(&ctx.showErrorEvenIfSilent, "show-error", "S", false, "Show error info even if silent mode on")
	flag.BoolVarP(&ctx.headOnly, "head", "I", false, "Only return headers (ignoring body content)")
	flag.BoolVarP(&ctx.includeHeadersInMainOutput, "include", "i", false, "Include headers (prepended to body content)")
	flag.StringSliceVarP(&ctx.cookies, "cookie", "b", empty, "HTTP cookie, raw HTTP cookie only (use -c for cookie jar files)")
	flag.StringVarP(&ctx.cookieJar, "cookie-jar", "c", "", "File for storing (read and write) cookies")
	flag.StringVarP(&ctx.uploadFile, "upload-file", "T", "", "Raw file to PUT (default) to the url given, not encoded")
	flag.Parse()
	initCLI(ctx)
}

func initCLI(ctx *CurlContext) {
	//if ctx.help || flag.NFlag() == 0 {
	//	printUsage()
	//	os.Exit(0)
	//}

	if ctx.verbose {
		if ctx.headerOutput == "" {
			ctx.headerOutput = ctx.output
		}
	}

	// do sanity checks and "fix" some parts left remaining from flag parsing
	tempUrl := strings.Join(flag.Args(), " ")
	if ctx.theUrl == "" && tempUrl != "" {
		ctx.theUrl = tempUrl
	}
	ctx.userAgent = strings.ReplaceAll(ctx.userAgent, "##DE"+"V##", "dev-branch")

	if ctx.silentFail || ctx.isSilent {
		ctx.isSilent = true   // implied
		ctx.silentFail = true // both are the same thing right now, we only emit errors (or content)
	}
	if ctx.headOnly {
		if ctx.headerOutput == "/dev/null" {
			ctx.headerOutput = "-"
		}
		ctx.SetMethodIfNotSet("HEAD")
	}

	if ctx.theUrl != "" {
		u, err := url.Parse(ctx.theUrl)
		changed := false
		if err != nil {
			panic(err)
		}
		if u.Scheme == "" {
			u.Scheme = "http"
			changed = true
		}
		if u.Host == "" {
			u.Host = "localhost"
			changed = true
		}
		if changed {
			ctx.theUrl = u.String()
		}
	}

	handleFormsAndFiles(ctx)

	// this should be LAST!
	ctx.SetMethodIfNotSet("GET")

	ctx.headerOutput = standardizeFileRef(ctx.headerOutput)
	ctx.output = standardizeFileRef(ctx.output)
	ctx.errorOutput = standardizeFileRef(ctx.errorOutput)
}

func standardizeFileRef(file string) string {
	if file == "/dev/null" || file == "null" || file == "" {
		return "/dev/null"
	}
	if file == "/dev/stderr" || file == "stderr" {
		return "/dev/stderr"
	}
	if file == "/dev/stdout" || file == "stdout" || file == "-" {
		return "/dev/stdout"
	}
	return file // no change
}

func handleFormsAndFiles(ctx *CurlContext) {
	if ctx.uploadFile != "" {
		f, err := os.ReadFile(ctx.uploadFile)
		if err != nil {
			logErrorF(ctx, "Failed to read file %s", ctx.uploadFile)
			os.Exit(9)
		}
		mime := mime.TypeByExtension(path.Ext(ctx.uploadFile))
		if mime == "" {
			mime = "application/octet-stream"
		}
		body := &bytes.Buffer{}
		body.Write(f)
		ctx.SetBody(body, mime, "POST")

	} else if len(ctx.formEncoded) > 0 {
		formBody := url.Values{}
		for _, item := range ctx.formEncoded {
			if strings.HasPrefix(item, "@") {
				filename := strings.TrimPrefix(item, "@")
				fullForm, err := os.ReadFile(filename)
				if err != nil {
					logErrorF(ctx, "Failed to read file %s", filename)
					os.Exit(9)
				}
				formLines := strings.Split(string(fullForm), "\n")
				for _, line := range formLines {
					splits := strings.SplitN(line, "=", 2)
					name := splits[0]
					value := splits[1]
					formBody.Set(name, value)
				}
			} else {
				splits := strings.SplitN(item, "=", 2)
				os.Stdout.WriteString(item)
				name := splits[0]
				value := splits[1]

				if strings.HasPrefix(value, "@") {
					filename := strings.TrimPrefix(value, "@")
					valueRaw, err := os.ReadFile(filename)
					if err != nil {
						logErrorF(ctx, "Failed to read file %s", filename)
						os.Exit(9)
					}
					//formBody.Set(name, base64.StdEncoding.EncodeToString(valueRaw))
					formBody.Set(name, string(valueRaw))
				} else {
					formBody.Set(name, value)
				}
			}
		}
		body := strings.NewReader(formBody.Encode())
		ctx.SetBody(body, "application/x-www-form-urlencoded", "POST")

	} else if len(ctx.formMultipart) > 0 {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for _, item := range ctx.formMultipart {
			if strings.HasPrefix(item, "@") {
				filename := strings.TrimPrefix(item, "@")
				fullForm, err := os.ReadFile(filename)
				if err != nil {
					logErrorF(ctx, "Failed to read file %s", filename)
					os.Exit(9)
				}
				formLines := strings.Split(string(fullForm), "\n")
				for _, line := range formLines {
					splits := strings.SplitN(line, "=", 2)
					name := splits[0]
					value := splits[1]
					part, _ := writer.CreateFormField(name)
					part.Write([]byte(value))
				}
			} else {
				splits := strings.SplitN(item, "=", 2)
				name := splits[0]
				value := splits[1]

				if strings.HasPrefix(value, "@") {
					filename := strings.TrimPrefix(value, "@")
					valueRaw, err := os.ReadFile(filename)
					if err != nil {
						logErrorF(ctx, "Failed to read file %s", filename)
						os.Exit(9)
					}
					part, _ := writer.CreateFormFile(name, path.Base(filename))
					part.Write(valueRaw)
				} else {
					part, _ := writer.CreateFormField(name)
					part.Write([]byte(value))
				}
			}
		}
		writer.Close()

		ctx.SetBody(body, "multipart/form-data; boundary="+writer.Boundary(), "POST")
	}
}

func logErrorF(ctx *CurlContext, entry string, value interface{}) {
	logError(ctx, fmt.Sprintf(entry, value))
}
func logError(ctx *CurlContext, entry string) {
	if (!ctx.isSilent && !ctx.silentFail) || !ctx.showErrorEvenIfSilent {
		writeToFileBytes(ctx.errorOutput, []byte(entry+"\n"))
	}
}
func writeToFileBytes(file string, body []byte) {
	if file == "/dev/null" {
		// do nothing
	} else if file == "/dev/stderr" {
		os.Stderr.Write(body)

	} else if file == "/dev/stdout" {
		os.Stdout.Write(body)
	} else {
		os.WriteFile(file, body, 0644)
	}
}

func printUsage() {
	fmt.Print("Usage: curl-go <command> [flags]\n\n")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
