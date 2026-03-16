package usecase

import (
	"context"
	"errors"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type GetDateSpotsInputPort interface {
	Execute(context.Context) (GetDateSpotsOutput, error)
}

type GetDateSpotsOutput struct {
	DateSpots []*model.DateSpot
	Addresses []*model.Address
}

type GetDateSpotsInteractor struct {
	DateSpotRepository repository.DateSpotRepository
	AddressRepository  repository.AddressRepository
}

func NewGetDateSpotsUsecase(
	dateSpotRepository repository.DateSpotRepository,
	addressRepository repository.AddressRepository,
) GetDateSpotsInputPort {
	return &GetDateSpotsInteractor{
		DateSpotRepository: dateSpotRepository,
		AddressRepository:  addressRepository,
	}
}

func (i *GetDateSpotsInteractor) Execute(ctx context.Context) (GetDateSpotsOutput, error) {
	dateSpots, err := i.DateSpotRepository.FindAll(ctx)
	if err != nil {
		return GetDateSpotsOutput{}, err
	}

	addresses := make([]*model.Address, 0, len(dateSpots))
	for _, ds := range dateSpots {
		address, err := i.AddressRepository.FindByDateSpotID(ctx, ds.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				addresses = append(addresses, nil)
				continue
			}
			return GetDateSpotsOutput{}, err
		}
		addresses = append(addresses, address)
	}

	return GetDateSpotsOutput{
		DateSpots: dateSpots,
		Addresses: addresses,
	}, nil
}
