package main

//I ommited the meta data under Response.Ledger.Transactions.Tx
type Response struct {
	Result string `json:"result"`
	Ledger struct {
		AccountHash     string `json:"account_hash"`
		CloseFlags      string `json:"close_flags"`
		CloseTime       int    `json:"close_time"`
		CloseTimeHuman  string `json:"close_time_human"`
		Closed          string `json:"closed"`
		Hash            string `json:"hash"`
		LedgerHash      string `json:"ledger_hash"`
		LedgerIndex     int    `json:"ledger_index"`
		ParentCloseTime int    `json:"parent_close_time"`
		ParentHash      string `json:"parent_hash"`
		SeqNum          string `json:"seqNum"`
		TotalCoins      string `json:"totalCoins"`
		TransactionHash string `json:"transaction_hash"`
		Transactions    []struct {
			Hash        string `json:"hash"`
			LedgerIndex int    `json:"ledger_index"`
			Date        string `json:"date"`
			Tx          struct {
				TransactionType string `json:"TransactionType"`
				Flags           int    `json:"Flags"`
				Sequence        int    `json:"Sequence"`
				Amount          string `json:"Amount"`
				Fee             string `json:"Fee"`
				SigningPubKey   string `json:"SigningPubKey"`
				TxnSignature    string `json:"TxnSignature"`
				Account         string `json:"Account"`
				Destination     string `json:"Destination"`
			} `json:"tx"`
		} `json:"transactions"`
	} `json:"ledger"`
}
