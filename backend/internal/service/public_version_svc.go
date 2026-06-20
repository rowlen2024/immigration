package service

import (
	"sort"
	"strings"

	"mygo-immigration/backend/internal/repository"
)

type PublicVersionResolver func(key string) (repository.PublicVersion, error)

type PublicVersionRegistry struct {
	resolvers map[string]PublicVersionResolver
}

func NewPublicVersionRegistry() *PublicVersionRegistry {
	return &PublicVersionRegistry{resolvers: make(map[string]PublicVersionResolver)}
}

func (r *PublicVersionRegistry) Register(prefix string, resolver PublicVersionResolver) {
	if prefix == "" || resolver == nil {
		return
	}
	r.resolvers[prefix] = resolver
}

func (r *PublicVersionRegistry) Resolve(keys []string) (map[string]string, error) {
	result := make(map[string]string)
	prefixes := make([]string, 0, len(r.resolvers))
	for prefix := range r.resolvers {
		prefixes = append(prefixes, prefix)
	}
	sort.Slice(prefixes, func(i, j int) bool {
		return len(prefixes[i]) > len(prefixes[j])
	})

	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				version, err := r.resolvers[prefix](key)
				if err != nil {
					return nil, err
				}
				result[key] = version.String()
				break
			}
		}
	}
	return result, nil
}

type PublicVersionRegistrar interface {
	RegisterPublicVersions(reg *PublicVersionRegistry)
}

type PublicVersionService struct {
	registry *PublicVersionRegistry
}

func NewPublicVersionService() *PublicVersionService {
	return &PublicVersionService{registry: NewPublicVersionRegistry()}
}

func (s *PublicVersionService) Resolve(keys []string) (map[string]string, error) {
	return s.registry.Resolve(keys)
}
