package bifcore

import (
	"errors"
	"github.com/teachain/bifcoresdk"
	"github.com/teachain/bifcoresdk/module/account"
	"github.com/teachain/bifcoresdk/module/blockchain"
	"github.com/teachain/bifcoresdk/module/contract"
	"github.com/teachain/bifcoresdk/types/request"
	"github.com/teachain/bifcoresdk/types/response"
)

type Client struct {
	client *bifcoresdk.SDK
}

func NewClient(url string) (*Client, error) {
	client, err := bifcoresdk.GetInstance(url)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// GetBIFAccountService 获取账户服务
func (c *Client) GetBIFAccountService() account.BIFAccountService {
	return c.client.GetBIFAccountService()
}

// GetBlockService 获取区块服务
func (c *Client) GetBlockService() blockchain.BIFBlockService {
	return c.client.GetBlockService()
}

// GetTransactionService 获取交易服务
func (c *Client) GetTransactionService() blockchain.BIFTransactionService {
	return c.client.GetTransactionService()
}

// GetContractService  获取合约服务
func (c *Client) GetContractService() contract.BIFContractService {
	return c.client.GetContractService()
}

func (c *Client) GetBlockNumber() (int64, error) {
	res := c.GetBlockService().GetBlockNumber()
	if res.ErrorCode != 0 {
		return 0, errors.New("unknown")
	}
	return res.Result.Header.BlockNumber, nil
}
func (c *Client) GetBlockHeaderByNumber(blockNumber int64) (response.BIFBlockHeader, error) {
	res := c.GetBlockService().GetBlockInfo(request.BIFBlockGetInfoRequest{BlockNumber: blockNumber})
	if res.ErrorCode != 0 {
		return response.BIFBlockHeader{}, errors.New("unknown")
	}
	return res.Result.Header, nil
}
func (c *Client) GetTransactionsByNumber(blockNumber int64) (response.BIFBlockGetTransactionsResult, error) {
	res := c.GetBlockService().GetTransactions(request.BIFBlockGetTransactionsRequest{BlockNumber: blockNumber})
	if res.ErrorCode != 0 {
		return response.BIFBlockGetTransactionsResult{}, errors.New("unknown")
	}
	return res.Result, nil
}

type T struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		TotalCount   int `json:"total_count"`
		Transactions []struct {
			ActualFee        int      `json:"actual_fee"`
			CloseTime        int64    `json:"close_time"`
			ContractTxHashes []string `json:"contract_tx_hashes,omitempty"`
			ErrorCode        int      `json:"error_code"`
			ErrorDesc        string   `json:"error_desc"`
			Hash             string   `json:"hash"`
			LedgerSeq        int      `json:"ledger_seq"`
			Signatures       []struct {
				PublicKey string `json:"public_key"`
				SignData  string `json:"sign_data"`
			} `json:"signatures,omitempty"`
			Transaction struct {
				FeeLimit   int `json:"fee_limit,omitempty"`
				GasPrice   int `json:"gas_price,omitempty"`
				Nonce      int `json:"nonce"`
				Operations []struct {
					PayCoin struct {
						DestAddress string `json:"dest_address"`
						Input       string `json:"input"`
					} `json:"pay_coin,omitempty"`
					SourceAddress string `json:"source_address,omitempty"`
					Type          int    `json:"type"`
					Log           struct {
						Datas []string `json:"datas"`
						Topic string   `json:"topic"`
					} `json:"log,omitempty"`
				} `json:"operations"`
				SourceAddress string `json:"source_address"`
			} `json:"transaction"`
			TxSize  int `json:"tx_size"`
			Trigger struct {
				Transaction struct {
					Hash  string `json:"hash"`
					Index int    `json:"index,omitempty"`
				} `json:"transaction"`
			} `json:"trigger,omitempty"`
		} `json:"transactions"`
	} `json:"result"`
}
