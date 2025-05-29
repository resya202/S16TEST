package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/resya202/S16TEST/internal/config"
	"github.com/resya202/S16TEST/internal/model"
)

// FetchAndStoreAll iterates a list of validators and stores their delegation snapshots

func FetchAndStoreAll(db *gorm.DB) error {
	validators := []string{
		"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
	}
	for _, val := range validators {
		if err := fetchAndStore(db, val); err != nil {
			return err
		}
	}
	return nil
}

// fetchAndStore pulls the delegations for one validator and writes hourly records
func fetchAndStore(db *gorm.DB, validatorAddr string) error {
	base := config.CosmosAPIBase
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations", base, validatorAddr)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse the JSON response
	var parsed struct {
		DelegationResponses []struct {
			Delegation struct {
				DelegatorAddress string `json:"delegator_address"`
			} `json:"delegation"`
			Balance struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"balance"`
		} `json:"delegation_responses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, entry := range parsed.DelegationResponses {
		amt, err := strconv.ParseInt(entry.Balance.Amount, 10, 64)
		if err != nil {
			return err
		}

		var last model.HourlyDelegation
		db.Where("validator_addr = ? AND delegator_addr = ?", validatorAddr, entry.Delegation.DelegatorAddress).
			Order("timestamp DESC").Limit(1).Find(&last)

		snap := model.HourlyDelegation{
			ValidatorAddr: validatorAddr,
			DelegatorAddr: entry.Delegation.DelegatorAddress,
			Timestamp:     now,
			AmountUAtom:   amt,
			ChangeUAtom:   amt - last.AmountUAtom,
		}
		if err := db.Create(&snap).Error; err != nil {
			return err
		}
	}

	return nil
}
