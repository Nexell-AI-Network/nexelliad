package blockrelay

import (
	"testing"

	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/netadapter/router"
)

func TestIBDBatchSizeLessThanRouteCapacity(t *testing.T) {
	// The `ibdBatchSize` constant must be equal at both syncer and syncee. Therefore, we do not want
	// to set it to `router.DefaultMaxMessages` to avoid confusion and human errors.
	// However, nonetheless we must enforce that it does not exceed `router.DefaultMaxMessages`
	if ibdBatchSize >= router.DefaultMaxMessages {
		t.Fatalf("IBD batch size (%d) must be smaller than router.DefaultMaxMessages (%d)",
			ibdBatchSize, router.DefaultMaxMessages)
	}
}
