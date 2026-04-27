// Package types defines the common capability interfaces used by set
// implementations and helper packages.
//
// The package separates read-only, mutable, clearable, and lock-coordinated
// behaviours into small composable interfaces. Concrete set types may
// implement any combination of these contracts.
//
// Other packages in the module use these interfaces to accept arbitrary set
// implementations without depending on a specific storage strategy.
package types
