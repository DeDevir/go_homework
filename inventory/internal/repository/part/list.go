package part

import (
	"context"
	"github.com/DeDevir/go_homework/inventory/internal/model"
	"github.com/DeDevir/go_homework/inventory/internal/repository/converter"
	"log"
)

func (r *repository) List(_ context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := converter.PartsFilterModelToInfo(filter)

	filteredParts := make([]*model.Part, 0, len(r.data))

	log.Printf("filter - %v", repoFilter)
	for _, repoPart := range r.data {
		if len(repoFilter.UUIDs) > 0 {
			log.Printf("filter not nil")
			if _, ok := repoFilter.UUIDs[repoPart.Uuid]; !ok {
				continue
			}
		}

		//if repoFilter.Names != nil {
		//	if _, ok := repoFilter.Names[repoPart]
		//}
		log.Printf("adding %v\n", repoPart.Uuid)
		filteredParts = append(filteredParts, converter.PartInfoToModel(repoPart))
	}

	log.Printf("all parts %v\n", r.data)
	log.Printf("finded parts after for %v\n", filteredParts)

	return filteredParts, nil
}
