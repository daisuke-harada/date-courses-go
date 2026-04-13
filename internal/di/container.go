package di

// BuildContainer はリポジトリ・ドメインサービス・ユースケースなど、アプリケーション全体の依存関係を Container に登録します。
func BuildContainer(ct *Container) {
	ProvideRepositories(ct)
	ProvideServices(ct)
	ProvideUsecases(ct)
}
