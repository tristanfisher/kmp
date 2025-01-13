package kmp

import (
	"reflect"
	"testing"
)

func TestKMPSearch(t *testing.T) {
	type args struct {
		body        string
		searchToken string
	}
	tests := []struct {
		name      string
		args      args
		want      []int
		wantCount int
	}{
		{
			name: "noBody",
			args: args{
				body:        "",
				searchToken: "hello",
			},
			want:      []int{},
			wantCount: 0,
		},
		{
			name: "noSearchToken",
			args: args{
				body:        "a",
				searchToken: "",
			},
			want:      []int{},
			wantCount: 0,
		},
		{
			name: "noBodyOrSearchToken",
			args: args{
				body:        "",
				searchToken: "",
			},
			want:      []int{},
			wantCount: 0,
		},
		{
			name: "findSingle",
			args: args{
				body:        "a",
				searchToken: "a",
			},
			// zeroth index in body, 1 count
			want:      []int{0},
			wantCount: 1,
		},
		{
			name: "findDouble",
			args: args{
				body:        "aa",
				searchToken: "a",
			},
			// zeroth index in body, 1 count
			want:      []int{0, 1},
			wantCount: 2,
		},
		{
			name: "findWithSpace",
			args: args{
				body:        `<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>`,
				searchToken: "Relationship Id",
			},
			want:      []int{1},
			wantCount: 1,
		},
		{
			name: "findWithSpaceAndRepeat",
			args: args{
				body:        `<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>`,
				searchToken: "Relationship Id",
			},
			want:      []int{1, 146},
			wantCount: 2,
		},
		{
			name: "multipleFind",
			args: args{
				body:        "Unto the pure all things are pure: but unto them that are defiled and unbelieving is nothing pure; but even their mind and conscience is defiled.",
				searchToken: "pure",
			},
			want:      []int{9, 29, 93},
			wantCount: 3,
		},
		{
			name: "miss",
			args: args{
				body:        "Unto the pure all things are pure: but unto them that are defiled and unbelieving is nothing pure; but even their mind and conscience is defiled.",
				searchToken: "unpure",
			},
			want:      []int{},
			wantCount: 0,
		},
		{
			name: "multipleMatch",
			args: args{
				body:        "Unto the pure all things are pure: but unto them that are defiled and unbelieving is nothing pure; but even their mind and conscience is defiled.",
				searchToken: "pure",
			},
			want:      []int{9, 29, 93},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Search(tt.args.body, tt.args.searchToken)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantCount {
				t.Errorf("Search() got1 = %v, want %v", got1, tt.wantCount)
			}
		})
	}
}

func Test_kmpPreprocessTable(t *testing.T) {
	type args struct {
		searchToken string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			// by definition, we have a zeroth entry of value 0
			name: "zeroLength",
			args: args{
				searchToken: "",
			},
			want: []int{},
		},
		{
			name: "singleChar",
			args: args{
				searchToken: "a",
			},
			want: []int{0},
		},
		{
			name: "repeatChar2",
			args: args{
				searchToken: "aa",
			},
			want: []int{0, 1},
		},
		{
			name: "repeatChar3",
			args: args{
				searchToken: "aaa",
			},
			want: []int{0, 1, 2},
		},
		{
			name: "adjacentRepeat",
			args: args{
				searchToken: "baba",
			},
			// note the 5 char hello repeat with len + 1
			want: []int{0, 0, 1, 2},
		},
		{
			name: "helloWorldWithRepeat",
			args: args{
				searchToken: "helloWorldhello",
			},
			// note the 5 char hello repeat with len + 1
			want: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := preprocessTable(tt.args.searchToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preprocessTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
