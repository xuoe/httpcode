package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

var mdnLinksFlag = flag.Bool("m", false, "include MDN links")

func init() {
	bin := os.Args[0]
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", bin)
		fmt.Fprintf(os.Stderr, " %s [code|text|pattern]...\n\n", bin)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	var (
		args  = flag.Args()
		codes []int
		rc    int
	)
	defer func() { os.Exit(rc) }()

	if len(args) == 0 {
		codes, rc = statusCodes[:], 0
	} else {
		codes, rc = findCodes(args)
	}
	out := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, code := range codes {
		fmt.Fprintln(out, statusLine(code, *mdnLinksFlag))
	}
	out.Flush()
}

func statusLine(code int, includeLinks bool) string {
	text := http.StatusText(code)
	s, args := "%d\t%s", []interface{}{code, text}
	if includeLinks {
		s, args = s+"\t%s", append(args, fmt.Sprintf("%s/%d", mdnURL, code))
	}
	return fmt.Sprintf(s, args...)
}

const mdnURL = "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"

func findCodes(terms []string) (codes []int, notFound int) {
	seen := make(map[int]struct{})
	for _, term := range terms {
		var found int
		for _, code := range expandTerm(term) {
			found++
			if _, ok := seen[code]; ok {
				continue
			}
			seen[code] = struct{}{}
			codes = append(codes, code)
		}
		if found == 0 {
			notFound++
		}
	}
	sort.Ints(codes)
	return
}

func expandTerm(term string) []int {
	term = strings.ReplaceAll(term, "_", "?")
	switch {
	case strings.ContainsAny(term, "*?"):
		return expandNumericPattern(term)
	case strings.ContainsAny(term, "0123456789"):
		code, err := strconv.Atoi(term)
		if err != nil || http.StatusText(code) == "" {
			return nil
		}
		return []int{code}
	default:
		return expandTextPattern(term)
	}
}

func expandNumericPattern(pattern string) (res []int) {
	for _, code := range statusCodes {
		if ok, _ := path.Match(pattern, strconv.Itoa(code)); ok {
			res = append(res, code)
		}
	}
	return
}

func expandTextPattern(pattern string) (res []int) {
	for _, code := range statusCodes {
		if stringContainsAnyOf(http.StatusText(code), strings.Fields(pattern)) {
			res = append(res, code)
		}
	}
	return
}

var statusCodes = [...]int{
	// 100s
	http.StatusContinue,
	http.StatusSwitchingProtocols,
	http.StatusProcessing,
	http.StatusEarlyHints,

	// 200s
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNonAuthoritativeInfo,
	http.StatusNoContent,
	http.StatusResetContent,
	http.StatusPartialContent,
	http.StatusMultiStatus,
	http.StatusAlreadyReported,
	http.StatusIMUsed,

	// 300s
	http.StatusMultipleChoices,
	http.StatusMovedPermanently,
	http.StatusFound,
	http.StatusSeeOther,
	http.StatusNotModified,
	http.StatusUseProxy,
	http.StatusTemporaryRedirect,
	http.StatusPermanentRedirect,

	// 400s
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusPaymentRequired,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
	http.StatusNotAcceptable,
	http.StatusProxyAuthRequired,
	http.StatusRequestTimeout,
	http.StatusConflict,
	http.StatusGone,
	http.StatusLengthRequired,
	http.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed,
	http.StatusTeapot,
	http.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity,
	http.StatusLocked,
	http.StatusFailedDependency,
	http.StatusTooEarly,
	http.StatusUpgradeRequired,
	http.StatusPreconditionRequired,
	http.StatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons,

	// 500s
	http.StatusInternalServerError,
	http.StatusNotImplemented,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
	http.StatusHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates,
	http.StatusInsufficientStorage,
	http.StatusLoopDetected,
	http.StatusNotExtended,
	http.StatusNetworkAuthenticationRequired,
}

func stringContainsAnyOf(haystack string, needle []string) bool {
	haystack = strings.ToLower(haystack)
	for _, needle := range needle {
		if strings.Contains(haystack, strings.ToLower(needle)) {
			return true
		}
	}
	return false
}
