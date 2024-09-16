package orderbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(10000, sellOrderA)
	ob.PlaceLimitOrder(9000, sellOrderB)

	require.Equal(t, 2, len(ob.Orders))
	require.Equal(t, ob.Orders[sellOrderA.ID], sellOrderA)
	require.Equal(t, ob.Orders[sellOrderB.ID], sellOrderB)
	require.Equal(t, 2, len(ob.asks))
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrder := NewOrder(false, 20)
	ob.PlaceLimitOrder(10000, sellOrder)

	buyOrder := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 1, len(ob.asks))
	assert.Equal(t, 10.0, ob.AskTotalVolume())
	assert.Equal(t, sellOrder, matches[0].Ask)
	assert.Equal(t, buyOrder, matches[0].Bid)
	assert.Equal(t, 10.0, matches[0].SizeFilled)
	assert.Equal(t, 10000.0, matches[0].Price)
	assert.True(t, buyOrder.IsFilled())
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)
	buyOrderD := NewOrder(true, 1)

	ob.PlaceLimitOrder(5000, buyOrderC)
	ob.PlaceLimitOrder(5000, buyOrderD)
	ob.PlaceLimitOrder(9000, buyOrderB)
	ob.PlaceLimitOrder(10000, buyOrderA)

	assert.Equal(t, 24.0, ob.BidTotalVolume())

	sellOrder := NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert.Equal(t, 4.0, ob.BidTotalVolume())
	assert.Equal(t, 3, len(matches))
	assert.Equal(t, 1, len(ob.bids))

	fmt.Println(matches)
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderbook()
	buyOrder := NewOrder(true, 4)
	ob.PlaceLimitOrder(10000.0, buyOrder)

	assert.Equal(t, 1, len(ob.bids))
	assert.Equal(t, 4.0, ob.BidTotalVolume())

	ob.CancelOrder(buyOrder)
	assert.Equal(t, 0.0, ob.BidTotalVolume())

	_, ok := ob.Orders[buyOrder.ID]
	assert.Equal(t, false, ok)
}
