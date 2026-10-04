package separation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// 允许的音频类型，做一层粗校验
var allowedExt = map[string]bool{
	".mp3": true, ".wav": true, ".flac": true, ".m4a": true, ".ogg": true,
}

// Separate 是分离服务的业务内核。
// 现在返回 mock 结果；将来换成调用 demucs/spleeter 时，只改这个函数，HTTP 层一行不动。
//
// 参数：
//
//	r       —— 音频字节流（由 HTTP 层从 multipart 中取出）
//	taskID  —— 由调用方（业务服务）生成并传入的任务 ID，便于全链路日志追踪
//	ext     —— 音频扩展名，如 ".wav"，决定 mock 出来的音轨后缀
func Separate(r io.Reader, taskID, ext string) ([]string, error) {
	// 真实分离很慢，模拟一下耗时，方便你在业务服务侧观察到 "running" 状态
	time.Sleep(2 * time.Second)

	// 读一下字节数，证明"音频确实传过来了"，同时把流消费掉避免连接挂起
	n, err := io.Copy(io.Discard, r)
	if err != nil {
		return nil, fmt.Errorf("读取音频流失败: %w", err)
	}
	if n == 0 {
		return nil, fmt.Errorf("音频内容为空")
	}
	log.Printf("[task %s] 收到音频 %d 字节，开始 mock 分离", taskID, n)

	// mock：返回四条常见音轨（人声/鼓/贝斯/其他），真实模型也是这样分
	stems := []string{
		fmt.Sprintf("/stems/%s/vocals%s", taskID, ext),
		fmt.Sprintf("/stems/%s/drums%s", taskID, ext),
		fmt.Sprintf("/stems/%s/bass%s", taskID, ext),
		fmt.Sprintf("/stems/%s/other%s", taskID, ext),
	}
	return stems, nil
}

func RegisterHandlers() {
	http.HandleFunc("/separate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// 32MB 留在内存，超过的自动落到临时文件，避免大音频把内存打爆
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("解析 multipart 失败: %v", err)
			writeErr(w, http.StatusBadRequest, "解析表单失败: "+err.Error())
			return
		}

		taskID := r.FormValue("taskID")
		if taskID == "" {
			writeErr(w, http.StatusBadRequest, "缺少 taskId")
			return
		}

		file, header, err := r.FormFile("audio")
		if err != nil {
			writeErr(w, http.StatusBadRequest, "缺少 audio 文件字段")
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedExt[ext] {
			writeErr(w, http.StatusBadRequest, "不支持的音频格式: "+ext)
			return
		}

		log.Printf("[task %s] 收到文件 %s (%d 字节)", taskID, header.Filename, header.Size)

		// 真正的分离动作委托给业务内核
		stems, err := Separate(file, taskID, ext)
		if err != nil {
			log.Printf("[task %s] 分离失败: %v", taskID, err)
			// 202 表示"收到了但没成功"，业务服务据此把任务标记为 failed
			writeErr(w, http.StatusAccepted, err.Error())
			return
		}

		// 组装响应：直接复用共享契约 SeperationTask
		task := SeperationTask{
			ID:     taskID,
			Status: StatusDone,
			Stems:  stems,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		log.Printf("[task %s] 分离完成，产出 %d 条音轨", taskID, len(stems))
	})
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
