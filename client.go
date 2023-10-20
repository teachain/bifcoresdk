package bifcoresdk

import (
	"errors"
	"github.com/teachain/bifcoresdk/module/account"
	"github.com/teachain/bifcoresdk/module/blockchain"
	"github.com/teachain/bifcoresdk/module/contract"
)

// BIFSDK bif sdk interface
type BIFSDK interface {
	// GetBIFAccountService ...
	GetBIFAccountService() account.BIFAccountService
	// GetBlockService ...
	GetBlockService() blockchain.BIFBlockService
	// GetTransactionService ...
	GetTransactionService() blockchain.BIFTransactionService
	// GetContractService ...
	GetContractService() contract.BIFContractService
}

// SDK ...
type SDK struct {
	url     string
	chainID int
}

// GetInstance initialize the SDK instance
func GetInstance(url string) (*SDK, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	sdk := &SDK{
		url: url,
	}

	return sdk, nil
}

func (sdk *SDK) GetBIFAccountService() account.BIFAccountService {
	return account.GetAccountInstance(sdk.url)
}

func (sdk *SDK) GetBlockService() blockchain.BIFBlockService {
	return blockchain.GetBlockInstance(sdk.url)
}

func (sdk *SDK) GetTransactionService() blockchain.BIFTransactionService {
	return blockchain.GetTransactionInstance(sdk.url)
}

func (sdk *SDK) GetContractService() contract.BIFContractService {
	return contract.GetContractInstance(sdk.url)
}
