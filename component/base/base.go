// Package base defines the base that a robot uses to move around.
package base

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	viamutils "go.viam.com/utils"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rlog"
)

// SubtypeName is a constant that identifies the component resource subtype string "arm".
const SubtypeName = resource.SubtypeName("base")

// Subtype is a constant that identifies the component resource subtype.
var Subtype = resource.NewSubtype(
	resource.ResourceNamespaceRDK,
	resource.ResourceTypeComponent,
	SubtypeName,
)

// Named is a helper for getting the named Base's typed resource name.
func Named(name string) resource.Name {
	return resource.NameFromSubtype(Subtype, name)
}

// A Base represents a physical base of a robot.
type Base interface {
	// MoveStraight moves the robot straight a given distance at a given speed. The method
	// can be requested to block until the move is complete. If a distance or speed of zero is given,
	// the base will stop.
	MoveStraight(ctx context.Context, distanceMillis int, millisPerSec float64, block bool) error

	// MoveArc moves the robot in an arc a given distance at a given speed and degs per second of movement.
	// The degs per sec represents the angular velocity the robot has during its movement. This function
	// can be requested to block until move is complete. If a distance of 0 is given the resultant motion
	// is a spin and if speed of 0 is given the base will stop.
	// Note: ramping affects when and how arc is performed, further improvements may be needed
	MoveArc(ctx context.Context, distanceMillis int, millisPerSec float64, degsPerSec float64, block bool) error

	// Spin spins the robot by a given angle in degrees at a given speed. The method can be requested
	// to block until the move is complete. If a speed of 0 the base will stop.
	Spin(ctx context.Context, angleDeg float64, degsPerSec float64, block bool) error

	// Stop stops the base. It is assumed the base stops immediately.
	Stop(ctx context.Context) error

	// WidthGet returns the width of the base in millimeters.
	WidthGet(ctx context.Context) (int, error)
}

var (
	_ = Base(&reconfigurableBase{})
	_ = resource.Reconfigurable(&reconfigurableBase{})
)

type reconfigurableBase struct {
	mu     sync.RWMutex
	actual Base
}

func (r *reconfigurableBase) ProxyFor() interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual
}

func (r *reconfigurableBase) MoveStraight(
	ctx context.Context, distanceMillis int, millisPerSec float64, block bool,
) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual.MoveStraight(ctx, distanceMillis, millisPerSec, block)
}

func (r *reconfigurableBase) MoveArc(
	ctx context.Context, distanceMillis int, millisPerSec float64, degsPerSec float64, block bool,
) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual.MoveArc(ctx, distanceMillis, millisPerSec, degsPerSec, block)
}

func (r *reconfigurableBase) Spin(ctx context.Context, angleDeg float64, degsPerSec float64, block bool) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual.Spin(ctx, angleDeg, degsPerSec, block)
}

func (r *reconfigurableBase) Stop(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual.Stop(ctx)
}

func (r *reconfigurableBase) WidthGet(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actual.WidthGet(ctx)
}

func (r *reconfigurableBase) Reconfigure(ctx context.Context, newBase resource.Reconfigurable) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	actual, ok := newBase.(*reconfigurableBase)
	if !ok {
		return errors.Errorf("expected new arm to be %T but got %T", r, newBase)
	}
	if err := viamutils.TryClose(ctx, r.actual); err != nil {
		rlog.Logger.Errorw("error closing old", "error", err)
	}
	r.actual = actual.actual
	return nil
}

// WrapWithReconfigurable converts a regular Base implementation to a reconfigurableBase.
// If base is already a reconfigurableBase, then nothing is done.
func WrapWithReconfigurable(r interface{}) (resource.Reconfigurable, error) {
	base, ok := r.(Base)
	if !ok {
		return nil, errors.Errorf("expected resource to be Base but got %T", r)
	}
	if reconfigurable, ok := base.(*reconfigurableBase); ok {
		return reconfigurable, nil
	}
	return &reconfigurableBase{actual: base}, nil
}

// A Move describes instructions for a robot to spin followed by moving straight.
type Move struct {
	DistanceMillis int
	MillisPerSec   float64
	AngleDeg       float64
	DegsPerSec     float64
	Block          bool
}

// DoMove performs the given move on the given base.
func DoMove(ctx context.Context, move Move, base Base) error {
	if move.AngleDeg != 0 {
		err := base.Spin(ctx, move.AngleDeg, move.DegsPerSec, move.Block)
		if err != nil {
			return err
		}
	}

	if move.DistanceMillis != 0 {
		err := base.MoveStraight(ctx, move.DistanceMillis, move.MillisPerSec, move.Block)
		if err != nil {
			return err
		}
	}

	return nil
}