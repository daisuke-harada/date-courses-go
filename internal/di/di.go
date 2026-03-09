package di

import (
	"log/slog"

	"go.uber.org/dig"
)

type Container struct {
	*dig.Container
}

// NewContainer は新しい Container を生成します。
func NewContainer() *Container {
	return &Container{Container: dig.New()}
}

// MustProvide は dig.Container.Provide を呼び出し、失敗した場合は panic します。
// DI の設定ミスはプログラマエラーであるため、起動時に即座にクラッシュさせます。
func (ct *Container) MustProvide(constructor any, opts ...dig.ProvideOption) {
	if err := ct.Provide(constructor, opts...); err != nil {
		slog.Error("di.MustProvide failed", "err", err)
		panic(err)
	}
}
