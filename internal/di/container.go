package di

// BuildContainer はリポジトリ・ユースケースなど、アプリケーション全体の依存関係を Container に登録します。
func BuildContainer(ct *Container) {
	ProvideRepositories(ct)
	ProvideUsecases(ct)
}
