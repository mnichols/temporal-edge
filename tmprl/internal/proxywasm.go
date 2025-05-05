package internal

//
//import (
//	"encoding/json"
//	"strings"
//	"unsafe"
//
//	proxywasm "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
//)
//
//func main() {
//	proxywasm.SetNewHttpContext(newContext)
//}
//
//func newContext(contextID uint32) proxywasm.HttpContext {
//	return &httpCtx{}
//}
//
//type httpCtx struct {
//	proxywasm.DefaultHttpContext
//}
//
//func (ctx *httpCtx) OnHttpRequestBody(bodySize int, endOfStream bool) proxywasm.Action {
//	if !endOfStream {
//		return proxywasm.ActionContinue
//	}
//
//	// Get the request body
//	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
//	if err != nil {
//		proxywasm.LogCritical("failed to get request body: " + err.Error())
//		return proxywasm.ActionContinue
//	}
//
//	// Parse JSON body
//	var jsonData map[string]interface{}
//	if err := json.Unmarshal(body, &jsonData); err != nil {
//		proxywasm.LogCritical("failed to parse JSON: " + err.Error())
//		return proxywasm.ActionContinue
//	}
//
//	// Encrypt the field (dummy encryption)
//	if val, ok := jsonData["sensitive_field"].(string); ok {
//		encrypted := dummyEncrypt(val)
//		jsonData["sensitive_field"] = encrypted
//	}
//
//	// Marshal back to JSON
//	newBody, err := json.Marshal(jsonData)
//	if err != nil {
//		proxywasm.LogCritical("failed to re-marshal JSON: " + err.Error())
//		return proxywasm.ActionContinue
//	}
//
//	// Replace request body
//	if err := proxywasm.SetHttpRequestBody(newBody); err != nil {
//		proxywasm.LogCritical("failed to set new body: " + err.Error())
//	}
//
//	return proxywasm.ActionContinue
//}
//
//// Dummy encryption: reverses the string
//func dummyEncrypt(input string) string {
//	runes := []rune(input)
//	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
//		runes[i], runes[j] = runes[j], runes[i]
//	}
//	return string(runes)
//}
