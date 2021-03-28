package cmd

import (
	"fmt"
	"github.com/smf8/gurl/internal/app/gurl/request"
	"github.com/spf13/cobra"
	"io"
	"math"
	"strings"
)

const (
	methodHelp  = `Specify HTTP method, either POST, GET (default), PATCH, DELETE, PUT`
	headersHelp = `HTTP headers in key:value format. use comma (,) as separator. Default Content-Type header is set to "application/x-www-form-urlencoded"`
	bodyHelp    = `HTTP request body, can't use with --json or --file`
	fileHelp    = `used to send a file as request body. This option also sets "Content-Type" header to "application/octet-stream". You can override this by setting Content-Type header manually with -H option`
	jsonHelp    = `used to specify json Body. This option also sets "Content-Type" header to "application/json". You can override this by setting Content-Type header manually with -H option`
	timeoutHelp = `Set client request Timeout. this timeout is only considered in client side.`
)

func main(gurl *request.GURL) error {
	httpRequest, err := gurl.ToHTTPRequest()
	if err != nil {
		return fmt.Errorf("failed to launch gurl command: %w", err)
	}

	fmt.Println(httpRequest)
	fmt.Println("method:", httpRequest.Method)

	body, err := io.ReadAll(httpRequest.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Data: %s", string(body))

	fmt.Println("headers:", httpRequest.Header)
	fmt.Println("URL:", httpRequest.URL.String())
	return nil
}

// NewCommand creates an instance of gURL cli cobra request
func NewCommand() *cobra.Command {

	gurl := new(request.GURL)
	gurl.Headers = make(map[string]string)
	gurl.QueryParams = make(map[string][]string)

	var headers *[]string

	gurlCommand := &cobra.Command{
		Use:   "gurl URL [-M Method] [-H Headers {key1:value1,...}] [-D Data | --json JsonData | --file FilePath] [--timeout timeout]",
		Short: "gurl is a simple curl rip off",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gurl.URL = args[0]

			for _, h := range *headers {
				record := strings.ToLower(h)
				header := strings.Split(record, ":")

				if len(header) != 2 {
					return fmt.Errorf("invalid header format. see -h for correct format")
				}

				gurl.Headers[header[0]] = header[1]
			}

			if err := main(gurl); err != nil {
				return err
			}
			return nil
		},
	}

	gurlCommand.Flags().StringVarP(&gurl.Method, "method", "M", "GET", methodHelp)
	headers = gurlCommand.Flags().StringSliceP("headers", "H", nil, headersHelp)
	gurlCommand.Flags().StringVarP(&gurl.Data, "data", "D", "", bodyHelp)
	gurlCommand.Flags().StringVar(&gurl.FilePath, "file", "", fileHelp)
	gurlCommand.MarkPersistentFlagDirname("file")
	gurlCommand.Flags().StringVar(&gurl.JSONMessage, "json", "", jsonHelp)
	gurlCommand.Flags().IntVar(&gurl.ClientTimeout, "timeout", math.MaxInt32, timeoutHelp)

	return gurlCommand
}
