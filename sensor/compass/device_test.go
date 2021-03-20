package compass_test

import (
	"context"
	"errors"
	"testing"

	"go.viam.com/robotcore/sensor/compass"
	"go.viam.com/robotcore/testutils/inject"

	"github.com/edaniels/test"
)

func TestMedianHeading(t *testing.T) {
	dev := &inject.Compass{}
	err1 := errors.New("whoops")
	dev.HeadingFunc = func(ctx context.Context) (float64, error) {
		return 0, err1
	}
	_, err := compass.MedianHeading(context.Background(), dev)
	test.That(t, err, test.ShouldEqual, err1)

	readings := []float64{1, 2, 3, 4, 4, 2, 4, 4, 1, 1, 2}
	readCount := 0
	dev.HeadingFunc = func(ctx context.Context) (float64, error) {
		reading := readings[readCount]
		readCount++
		return reading, nil
	}
	med, err := compass.MedianHeading(context.Background(), dev)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, med, test.ShouldEqual, 3)
}
