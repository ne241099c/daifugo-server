// Package ai は、Python(scikit-learn)で学習した線形SVMの「着手判断」を
// Go 側で推論するためのパッケージ。
//
// 線形SVMの推論は「標準化した特徴ベクトル x と重み w の内積 + b」だけなので、
// 学習結果(w, b と StandardScaler の mean/std)を model.json として持ち込み、
// Go 上で合法手を1つずつ採点して、スコア最大の手を選ぶ。
//
// model.json は Jupyter ノートブックの「モデル書き出し」セルで生成する。
package ai

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed model.json
var modelJSON []byte

// Model は学習済み線形SVMのパラメータ。
type Model struct {
	Features []string  `json:"features"` // 特徴量の名前（順序の確認用）
	W        []float64 `json:"w"`        // 重みベクトル（次元 = len(Features)）
	B        float64   `json:"b"`        // バイアス
	Mean     []float64 `json:"mean"`     // StandardScaler の平均
	Std      []float64 `json:"std"`      // StandardScaler の標準偏差
}

// Default は埋め込んだ model.json からモデルを読み込む。
func Default() (*Model, error) {
	var m Model
	if err := json.Unmarshal(modelJSON, &m); err != nil {
		return nil, fmt.Errorf("ai: model.json の読み込みに失敗: %w", err)
	}
	d := len(m.W)
	if d == 0 || len(m.Mean) != d || len(m.Std) != d {
		return nil, fmt.Errorf("ai: model.json の次元が不正 (w=%d mean=%d std=%d)",
			len(m.W), len(m.Mean), len(m.Std))
	}
	if len(m.Features) != 0 && len(m.Features) != d {
		return nil, fmt.Errorf("ai: features と w の次元が不一致 (%d != %d)", len(m.Features), d)
	}
	// 標準偏差が 0 の特徴で 0 割りを防ぐ。
	for i := range m.Std {
		if m.Std[i] == 0 {
			m.Std[i] = 1
		}
	}
	return &m, nil
}

// score は特徴ベクトルに対する SVM の決定関数 w·standardize(x) + b を返す。
// 値が大きいほど「良い手」。
func (m *Model) score(f []float64) float64 {
	s := m.B
	for j, v := range f {
		s += m.W[j] * ((v - m.Mean[j]) / m.Std[j])
	}
	return s
}
