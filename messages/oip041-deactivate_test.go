package messages

import "testing"

func TestVerifyOIP041Deactivate(t *testing.T) {
	// ToDo: ...
}

var deactivateJSON string = `{
  "oip-041": {
    "deactivateArtifact": {
      "txid": "96bad8e17f908da4c695c58b0f843a03928e338b361b3035eda16a864eafc3a2",
      "timestamp": 1481697196
    },
    "signature": "H8wPKDHSTrrJIY4RVzoWuWlt5Fta4PaWbWNQvrGt9hRyBFrB2YDoVnhctC4V08KnGMD/CZKbA4cPKysStVq8jyE="
  }
}`

// 'string' artifact that needs deactivation, need to find key to FNRCCaR7Y4T4oY5KmUMgjRULsMp7uh6uZY
//`{
//  "oip-041": {
//    "deactivateArtifact": {
//      "txid": "b1fbdd9e2d36696fb35d05434c449a468a75f7f0b161b4bfaac7379a453d5c12",
//      "timestamp": 1523642725
//    },
//    "signature": ""
//  }
//}`
// b1fbdd9e2d36696fb35d05434c449a468a75f7f0b161b4bfaac7379a453d5c12-FMBMKcV9fY8G8iahwTAYzMg4jqrbAEc7ED-1523642725
