package wails

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	domainPlugin "github.com/cli-tool-center/tool-center/internal/domain/plugin"
	"github.com/cli-tool-center/tool-center/internal/infrastructure/filesystem"
)

// callbackMessage 复刻 Wails dispatcher 用来回传结果的结构：
// internal/frontend/dispatcher/calls.go 里 CallbackMessage.Result 是 interface{}，
// 把 *PluginDto 装进 interface{} 再 json.Marshal —— 这就是前端拿到的 JSON。
// 该测试验证 json.RawMessage 在这条路径上保持原始 JSON 对象（而非 base64），
// 且 schema properties 字段顺序与源文件一致。
type callbackMessage struct {
	CallbackID string      `json:"callbackID"`
	Result     interface{} `json:"result"`
	Err        string      `json:"err,omitempty"`
}

// loadTiaExport 加载真实 tia-export 插件，返回完整 domain Plugin。
func loadTiaExport(t *testing.T) *domainPlugin.Plugin {
	t.Helper()
	// tia-export 目录相对仓库根。go test 的工作目录是包目录。
	dir, err := filepath.Abs(filepath.Join("..", "..", "..", "plugins"))
	if err != nil {
		t.Fatal(err)
	}
	loader := filesystem.NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("load plugins: %v", err)
	}
	for _, p := range plugins {
		if p.Metadata.ID == "tia-export" {
			return p
		}
	}
	t.Fatal("tia-export plugin not found in plugins/")
	return nil
}

// TestPluginDto_RawMessageSurvivesWailsMarshal 验证端到端序列化：
// pluginToFull → *PluginDto 装进 callbackMessage.Result(interface{}) → json.Marshal，
// 结果 JSON 里 form.schema 必须是嵌套对象（不是 base64 字符串），
// 且 properties 字段顺序必须是 indir, outdir, exportMode, ...（与 schema.json 一致）。
func TestPluginDto_RawMessageSurvivesWailsMarshal(t *testing.T) {
	p := loadTiaExport(t)
	dto := pluginToFull(p)

	// 复刻 Wails dispatcher 的回传序列化路径。
	msg := callbackMessage{CallbackID: "cb1", Result: &dto}
	out, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal callbackMessage: %v", err)
	}
	jsonStr := string(out)

	// 1) schema 必须是对象，不是 base64 字符串。base64 只含 A-Za-z0-9+/=，
	//    对象会以 {"$schema":... 或 {"type":... 开头。schema.json 带 $schema 字段，
	//    所以这里找 properties 对象头。
	schemaObjIdx := strings.Index(jsonStr, `"schema":{"$schema":`)
	if schemaObjIdx < 0 {
		schemaObjIdx = strings.Index(jsonStr, `"schema":{"type":`)
	}
	if schemaObjIdx < 0 {
		// 可能是 base64 —— 截取 schema 字段值前后 120 字符帮助诊断。
		i := strings.Index(jsonStr, `"schema":"`)
		if i >= 0 {
			end := i + 120
			if end > len(jsonStr) {
				end = len(jsonStr)
			}
			t.Fatalf("schema 透传失败：前端会收到 base64 而非对象。片段: %s", jsonStr[i:end])
		}
		t.Fatalf("在结果 JSON 里找不到 schema 对象。完整片段前 400 字符: %s", trunc(jsonStr, 400))
	}

	// 2) properties 字段顺序必须保留 schema.json 定义顺序。
	//    用 json.Decoder Token 流从 schema 子对象提取 properties 键序，
	//    不经过 map（map 会字母序）。
	//    先把外层 callbackMessage.result.form.schema 抽出来。
	var top struct {
		Result struct {
			Form struct {
				Schema json.RawMessage `json:"schema"`
			} `json:"form"`
		} `json:"result"`
	}
	if err := json.Unmarshal(out, &top); err != nil {
		t.Fatalf("unmarshal back: %v", err)
	}
	order := propertyOrderFromBytes(t, top.Result.Form.Schema)

	want := []string{
		"indir", "outdir", "exportMode", "keepFolderStructure", "maxWorkers",
		"projectFilter", "sclFormat", "stlFormat", "ladFormat", "dbFormat",
		"udtFormat", "safetyDbFormat", "safetyUdtFormat", "umacUser", "umacPassword",
		"portalMode", "logfile", "loglevel",
	}
	if len(order) != len(want) {
		t.Fatalf("property count = %d, want %d; got %v", len(order), len(want), order)
	}
	for i, w := range want {
		if order[i] != w {
			t.Errorf("property[%d] = %q, want %q (full order: %v)", i, order[i], w, order)
		}
	}

	// 3) enum 字段（sclFormat）必须以对象形式存在，确保前端 select widget 能拿到 enumOptions。
	if !strings.Contains(jsonStr, `"enum":["ExternalSource","SimaticML","SimaticSD"]`) {
		// 字段顺序可能被 JSON.Marshal 保留（RawMessage 原样），enum 数组也应原样。
		t.Errorf("sclFormat enum 数组未在序列化结果中按原样出现，可能影响前端 select 渲染")
	}
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// propertyOrderFromBytes 用 json.Decoder Token 流从 schema 字节里按出现顺序
// 提取 properties 下的字段名，不经过 map（避免字母序）。
func propertyOrderFromBytes(t *testing.T, schemaJSON []byte) []string {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(string(schemaJSON)))
	tok, err := dec.Token()
	if err != nil {
		t.Fatalf("first token: %v", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		t.Fatalf("schema 顶层不是对象: %v", tok)
	}
	for {
		key, err := dec.Token()
		if err != nil {
			t.Fatalf("next key: %v", err)
		}
		if d, ok := key.(json.Delim); ok && d == '}' {
			return nil
		}
		keyStr, _ := key.(string)
		if keyStr == "properties" {
			tt, err := dec.Token()
			if err != nil {
				t.Fatalf("properties open: %v", err)
			}
			d, ok := tt.(json.Delim)
			if !ok || d != '{' {
				t.Fatalf("properties 不是对象: %v", tt)
			}
			var order []string
			for {
				name, err := dec.Token()
				if err != nil {
					t.Fatalf("property name: %v", err)
				}
				if d, ok := name.(json.Delim); ok && d == '}' {
					return order
				}
				ns, _ := name.(string)
				order = append(order, ns)
				var v any
				if err := dec.Decode(&v); err != nil {
					t.Fatalf("skip property value: %v", err)
				}
			}
		}
		// 跳过非 properties 字段的值。
		var v any
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("skip value: %v", err)
		}
	}
}

// 避免未使用导入报错。
