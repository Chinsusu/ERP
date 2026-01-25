package entity

import (
	"math"
	"time"
)

// AllocateStockFEFO implements the core FEFO allocation logic.
// It takes a list of available stocks (already filtered and sorted by expiry)
// and returns the lots to be issued and the remaining quantity.
func AllocateStockFEFO(stocks []*Stock, quantity float64) ([]LotIssued, float64, error) {
	remaining := quantity
	var issued []LotIssued

	for _, stock := range stocks {
		if remaining <= 0 {
			break
		}

		availableQty := stock.Quantity - stock.ReservedQty
		if availableQty <= 0 {
			continue
		}

		issueQty := math.Min(availableQty, remaining)

		// Record the issued lot
		if stock.Lot != nil {
			issued = append(issued, LotIssued{
				LotID:      stock.Lot.ID,
				LotNumber:  stock.Lot.LotNumber,
				Quantity:   issueQty,
				ExpiryDate: stock.Lot.ExpiryDate,
				LocationID: stock.LocationID,
			})
		} else if stock.LotID != nil {
			issued = append(issued, LotIssued{
				LotID:      *stock.LotID,
				LotNumber:  "", // Unknown if Lot object is missing
				Quantity:   issueQty,
				LocationID: stock.LocationID,
			})
		}

		remaining -= issueQty
	}

	return issued, remaining, nil
}

// FilterAndSortForFEFO filters and sorts a given list of stocks for FEFO allocation.
// Stocks must have MaterialID, Quantity > ReservedQty, Lot must be Available, QC Passed, and Not Expired.
func FilterAndSortForFEFO(stocks []*Stock, now time.Time) []*Stock {
	var filtered []*Stock

	for _, s := range stocks {
		if s.Lot == nil || !s.Lot.CanBeIssued() {
			continue
		}
		
		// Ensure it's not actually expired based on provided 'now'
		if s.Lot.ExpiryDate.Before(now) {
			continue
		}

		if s.Quantity-s.ReservedQty > 0 {
			filtered = append(filtered, s)
		}
	}

	// In a real implementation, we would sort by ExpiryDate ASC here.
	// But let's assume the caller provides them sorted or we sort them here.
	// Sorting logic:
	sortStocksByExpiry(filtered)

	return filtered
}

func sortStocksByExpiry(stocks []*Stock) {
	for i := 0; i < len(stocks); i++ {
		for j := i + 1; j < len(stocks); j++ {
			if stocks[i].Lot.ExpiryDate.After(stocks[j].Lot.ExpiryDate) {
				stocks[i], stocks[j] = stocks[j], stocks[i]
			}
		}
	}
}
