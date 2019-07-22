package trafficshape

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Converts a sorted slice of Throttles to their ChangeBandwidth actions. In adddition, checks for
// overlapping throttle ranges. Returns a slice of actions and an error specifying if the throttles
// passed the non-overlapping verification.
// Idea: For every throttle, add two ChangeBandwidth actions (one for start and one for end), unless
// the ending byte of one throttle is the same as the starting byte of the next throttle, in which
// case we do not add the end ChangeBandwidth for the first throttle, or if the end of a throttle
// is -1 (representing till the end of file), in which case we do not add the end ChangeBandwidth
// action for the throttle. Note, we only allow the last throttle in the sorted list to have an end
// of -1, since otherwise there would be an overlap.
func getActionsFromThrottles(throttles []*Throttle, defaultBandwidth int64) ([]Action, error) {

	lenThr := len(throttles)
	var actions []Action
	for index, throttle := range throttles {
		start := throttle.ByteStart
		end := throttle.ByteEnd

		if index == lenThr-1 {
			if end == -1 {
				actions = append(actions,
					Action(&ChangeBandwidth{
						Byte:      start,
						Bandwidth: throttle.Bandwidth,
					}))
			} else {
				actions = append(actions,
					Action(&ChangeBandwidth{
						Byte:      start,
						Bandwidth: throttle.Bandwidth,
					}),
					Action(&ChangeBandwidth{
						Byte:      end,
						Bandwidth: defaultBandwidth,
					}))
			}
			break
		}

		if end > throttles[index+1].ByteStart || end == -1 {
			return actions, errors.New("overlapping throttle intervals found")
		}

		if end == throttles[index+1].ByteStart {
			actions = append(actions,
				Action(&ChangeBandwidth{
					Byte:      start,
					Bandwidth: throttle.Bandwidth,
				}))
		} else {
			actions = append(actions,
				Action(&ChangeBandwidth{
					Byte:      start,
					Bandwidth: throttle.Bandwidth,
				}),
				Action(&ChangeBandwidth{
					Byte:      end,
					Bandwidth: defaultBandwidth,
				}))
		}
	}
	return actions, nil
}

// Parses the Trafficshape object and populates/updates Traffficshape.Shapes,
// while performing verifications. Returns an error in case a verification check fails.
func parseShapes(ts *Trafficshape) error {
	var err error
	for shapeIndex, shape := range ts.Shapes {
		if shape == nil {
			return fmt.Errorf("nil shape at index: %d", shapeIndex)
		}
		if shape.URLRegex == "" {
			return fmt.Errorf("no url_regex for shape at index: %d", shapeIndex)
		}

		if _, err = regexp.Compile(shape.URLRegex); err != nil {
			return fmt.Errorf("url_regex for shape at index doesn't compile: %d", shapeIndex)
		}

		if shape.MaxBandwidth < 0 {
			return fmt.Errorf("max_bandwidth cannot be negative for shape at index: %d", shapeIndex)
		}

		if shape.MaxBandwidth == 0 {
			shape.MaxBandwidth = DefaultBitrate / 8
		}

		shape.WriteBucket = NewBucket(shape.MaxBandwidth, time.Second)

		// Verify and process the throttles, filling in their ByteStart and ByteEnd.
		for throttleIndex, throttle := range shape.Throttles {
			if throttle == nil {
				return fmt.Errorf("nil throttle at index %d in shape index %d", throttleIndex, shapeIndex)
			}

			if throttle.Bandwidth <= 0 {
				return fmt.Errorf("invalid bandwidth: %d at throttle index %d in shape index %d",
					throttle.Bandwidth, throttleIndex, shapeIndex)
			}
			sl := strings.Split(throttle.Bytes, "-")

			if len(sl) != 2 {
				return fmt.Errorf("invalid bytes: %s at throttle index %d in shape index %d",
					throttle.Bytes, throttleIndex, shapeIndex)
			}

			start := sl[0]
			end := sl[1]

			if start == "" {
				throttle.ByteStart = 0
			} else {
				throttle.ByteStart, err = strconv.ParseInt(start, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid bytes: %s at throttle index %d in shape index %d",
						throttle.Bytes, throttleIndex, shapeIndex)
				}
			}

			if end == "" {
				throttle.ByteEnd = -1
			} else {
				throttle.ByteEnd, err = strconv.ParseInt(end, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid bytes: %s at throttle index %d in shape index %d",
						throttle.Bytes, throttleIndex, shapeIndex)
				}
				if throttle.ByteEnd < throttle.ByteStart {
					return fmt.Errorf("invalid bytes: %s at throttle index %d in shape index %d",
						throttle.Bytes, throttleIndex, shapeIndex)
				}
			}

			if throttle.ByteStart == throttle.ByteEnd {
				return fmt.Errorf("invalid bytes: %s at throttle index %d in shape index %d",
					throttle.Bytes, throttleIndex, shapeIndex)
			}
		}
		// Fill in the actions, while performing verification.
		shape.Actions = make([]Action, len(shape.Halts)+len(shape.CloseConnections))

		for index, value := range shape.Halts {
			if value == nil {
				return fmt.Errorf("nil halt at index %d in shape index %d", index, shapeIndex)
			}
			if value.Duration < 0 || value.Byte < 0 {
				return fmt.Errorf("invalid halt at index %d in shape index %d", index, shapeIndex)
			}
			if value.Count == 0 {
				return fmt.Errorf(" 0 count for halt at index %d in shape index %d", index, shapeIndex)
			}
			shape.Actions[index] = Action(value)
		}
		offset := len(shape.Halts)
		for index, value := range shape.CloseConnections {
			if value == nil {
				return fmt.Errorf("nil close_connection at index %d in shape index %d",
					index, shapeIndex)
			}
			if value.Byte < 0 {
				return fmt.Errorf("invalid close_connection at index %d in shape index %d",
					index, shapeIndex)
			}
			if value.Count == 0 {
				return fmt.Errorf("0 count for close_connection at index %d in shape index %d",
				 index, shapeIndex)
			}
			shape.Actions[offset+index] = Action(value)
		}

		sort.SliceStable(shape.Throttles, func(i, j int) bool { 
			return shape.Throttles[i].ByteStart < shape.Throttles[j].ByteStart 
			})

		defaultBandwidth := DefaultBitrate / 8
		if shape.MaxBandwidth > 0 {
			defaultBandwidth = shape.MaxBandwidth
		}

		throttleActions, err := getActionsFromThrottles(shape.Throttles, defaultBandwidth)
		if err != nil {
			return fmt.Errorf("err: %s in shape index %d", err.Error(), shapeIndex)
		}
		shape.Actions = append(shape.Actions, throttleActions...)

		// Sort the actions according to their byte offset.
		sort.SliceStable(shape.Actions, func(i, j int) bool { return shape.Actions[i].getByte() < shape.Actions[j].getByte() })
	}
	return nil
}

