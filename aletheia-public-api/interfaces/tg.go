//go:generate tg client -go --services . --outPath ../pkg/clients/stories
//go:generate tg transport --services . --out ../internal/transport --outSwagger ../api/openapi.yaml
//go:generate goimports -l -w ../internal/transport ../pkg/clients/stories
package interfaces
