//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package graphql

import (
	"reflect"
	"testing"
)

func Test_extractNearTextFn(t *testing.T) {
	type args struct {
		source map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *NearTextParams
	}{
		{
			"Extract with concepts",
			args{
				source: map[string]interface{}{
					"concepts": []interface{}{"c1", "c2", "c3"},
				},
			},
			&NearTextParams{
				Values: []string{"c1", "c2", "c3"},
			},
		},
		{
			"Extract with concepts, certainty, limit and network",
			args{
				source: map[string]interface{}{
					"concepts":  []interface{}{"c1", "c2", "c3"},
					"certainty": float64(0.4),
					"limit":     100,
					"network":   true,
				},
			},
			&NearTextParams{
				Values:    []string{"c1", "c2", "c3"},
				Certainty: 0.4,
				Limit:     100,
				Network:   true,
			},
		},
		{
			"Extract with moveTo and moveAwayFrom",
			args{
				source: map[string]interface{}{
					"concepts":  []interface{}{"c1", "c2", "c3"},
					"certainty": float64(0.89),
					"limit":     500,
					"network":   false,
					"moveTo": map[string]interface{}{
						"concepts": []interface{}{"positive"},
						"force":    float64(0.5),
					},
					"moveAwayFrom": map[string]interface{}{
						"concepts": []interface{}{"epic"},
						"force":    float64(0.25),
					},
				},
			},
			&NearTextParams{
				Values:    []string{"c1", "c2", "c3"},
				Certainty: 0.89,
				Limit:     500,
				Network:   false,
				MoveTo: ExploreMove{
					Values: []string{"positive"},
					Force:  0.5,
				},
				MoveAwayFrom: ExploreMove{
					Values: []string{"epic"},
					Force:  0.25,
				},
			},
		},
		{
			"Extract with moveTo and moveAwayFrom (and objects)",
			args{
				source: map[string]interface{}{
					"concepts":  []interface{}{"c1", "c2", "c3"},
					"certainty": float64(0.89),
					"limit":     500,
					"network":   false,
					"moveTo": map[string]interface{}{
						"concepts": []interface{}{"positive"},
						"force":    float64(0.5),
						"objects": []interface{}{
							map[string]interface{}{
								"id": "moveTo-uuid1",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveTo-uuid2",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveTo-uuid3",
							},
						},
					},
					"moveAwayFrom": map[string]interface{}{
						"concepts": []interface{}{"epic"},
						"force":    float64(0.25),
						"objects": []interface{}{
							map[string]interface{}{
								"id": "moveAwayFrom-uuid1",
							},
							map[string]interface{}{
								"id": "moveAwayFrom-uuid2",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveAwayFrom-uuid3",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveAwayFrom-uuid4",
							},
						},
					},
				},
			},
			&NearTextParams{
				Values:    []string{"c1", "c2", "c3"},
				Certainty: 0.89,
				Limit:     500,
				Network:   false,
				MoveTo: ExploreMove{
					Values: []string{"positive"},
					Force:  0.5,
					Objects: []ObjectMove{
						{ID: "moveTo-uuid1"},
						{Beacon: "weaviate://localhost/moveTo-uuid2"},
						{Beacon: "weaviate://localhost/moveTo-uuid3"},
					},
				},
				MoveAwayFrom: ExploreMove{
					Values: []string{"epic"},
					Force:  0.25,
					Objects: []ObjectMove{
						{ID: "moveAwayFrom-uuid1"},
						{ID: "moveAwayFrom-uuid2"},
						{Beacon: "weaviate://localhost/moveAwayFrom-uuid3"},
						{Beacon: "weaviate://localhost/moveAwayFrom-uuid4"},
					},
				},
			},
		},
		{
			"Extract with moveTo and moveAwayFrom (and doubled objects)",
			args{
				source: map[string]interface{}{
					"concepts":  []interface{}{"c1", "c2", "c3"},
					"certainty": float64(0.89),
					"limit":     500,
					"network":   false,
					"moveTo": map[string]interface{}{
						"concepts": []interface{}{"positive"},
						"force":    float64(0.5),
						"objects": []interface{}{
							map[string]interface{}{
								"id":     "moveTo-uuid1",
								"beacon": "weaviate://localhost/moveTo-uuid2",
							},
							map[string]interface{}{
								"id":     "moveTo-uuid1",
								"beacon": "weaviate://localhost/moveTo-uuid2",
							},
						},
					},
					"moveAwayFrom": map[string]interface{}{
						"concepts": []interface{}{"epic"},
						"force":    float64(0.25),
						"objects": []interface{}{
							map[string]interface{}{
								"id":     "moveAwayFrom-uuid1",
								"beacon": "weaviate://localhost/moveAwayFrom-uuid1",
							},
							map[string]interface{}{
								"id":     "moveAwayFrom-uuid2",
								"beacon": "weaviate://localhost/moveAwayFrom-uuid2",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveAwayFrom-uuid3",
							},
							map[string]interface{}{
								"beacon": "weaviate://localhost/moveAwayFrom-uuid4",
							},
						},
					},
				},
			},
			&NearTextParams{
				Values:    []string{"c1", "c2", "c3"},
				Certainty: 0.89,
				Limit:     500,
				Network:   false,
				MoveTo: ExploreMove{
					Values: []string{"positive"},
					Force:  0.5,
					Objects: []ObjectMove{
						{ID: "moveTo-uuid1", Beacon: "weaviate://localhost/moveTo-uuid2"},
						{ID: "moveTo-uuid1", Beacon: "weaviate://localhost/moveTo-uuid2"},
					},
				},
				MoveAwayFrom: ExploreMove{
					Values: []string{"epic"},
					Force:  0.25,
					Objects: []ObjectMove{
						{ID: "moveAwayFrom-uuid1", Beacon: "weaviate://localhost/moveAwayFrom-uuid1"},
						{ID: "moveAwayFrom-uuid2", Beacon: "weaviate://localhost/moveAwayFrom-uuid2"},
						{Beacon: "weaviate://localhost/moveAwayFrom-uuid3"},
						{Beacon: "weaviate://localhost/moveAwayFrom-uuid4"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractNearTextFn(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractNearTextFn() = %v, want %v", got, tt.want)
			}
		})
	}
}
