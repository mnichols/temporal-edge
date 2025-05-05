package internal

import (
	"encoding/base64"
	"encoding/json"

	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}
func init() {
	proxywasm.SetHttpContext(newContext)
}

func newContext(contextID uint32) types.HttpContext {
	return &httpContext{}
}

type httpContext struct {
	types.DefaultHttpContext
}

func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if !endOfStream {
		return types.ActionContinue
	}

	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogCriticalf("failed to read request body: %v", err)
		return types.ActionContinue
	}

	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		proxywasm.LogCriticalf("failed to parse JSON: %v", err)
		return types.ActionContinue
	}

	// Navigate to input.payloads[*].data
	input, ok := jsonBody["input"].(map[string]interface{})
	if !ok {
		proxywasm.LogWarn("missing or malformed 'input' field")
		return types.ActionContinue
	}

	payloads, ok := input["payloads"].([]interface{})
	if !ok {
		proxywasm.LogWarn("missing or malformed 'payloads' field")
		return types.ActionContinue
	}

	for _, item := range payloads {
		payloadMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		dataField, ok := payloadMap["data"].(string)
		if !ok {
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(dataField)
		if err != nil {
			proxywasm.LogWarnf("invalid base64 in data field: %v", err)
			continue
		}

		encrypted := dummyEncrypt(string(decoded))
		payloadMap["data"] = base64.StdEncoding.EncodeToString([]byte(encrypted))
	}

	// Marshal and set the new request body
	newBody, err := json.Marshal(jsonBody)
	if err != nil {
		proxywasm.LogCriticalf("failed to re-marshal JSON: %v", err)
		return types.ActionContinue
	}

	if err := proxywasm.ReplaceHttpRequestBody(newBody); err != nil {
		proxywasm.LogCriticalf("failed to set new body: %v", err)
		return types.ActionContinue
	}

	return types.ActionContinue
}

// dummyEncrypt simply reverses the input string — replace with real logic
func dummyEncrypt(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
