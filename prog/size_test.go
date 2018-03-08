// Copyright 2016 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package prog

import (
	"bytes"
	"strings"
	"testing"
)

func TestAssignSizeRandom(t *testing.T) {
	target, rs, iters := initTest(t)
	for i := 0; i < iters; i++ {
		p := target.Generate(rs, 10, nil)
		data0 := p.Serialize()
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
		if data1 := p.Serialize(); !bytes.Equal(data0, data1) {
			t.Fatalf("different lens assigned, initial:\n%s\nnew:\n%s\n", data0, data1)
		}
		p.Mutate(rs, 10, nil, nil)
		data0 = p.Serialize()
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
	}
}

func TestAssignSize(t *testing.T) {
	target := initTargetTest(t, "test", "64")
	tests := []struct {
		unsizedProg string
		sizedProg   string
	}{
		{
			"syz_test$length0(&(0x7f0000000000)={0xff, 0x0})",
			"syz_test$length0(&(0x7f0000000000)={0xff, 0x2})",
		},
		{
			"syz_test$length1(&(0x7f0000001000)={0xff, 0x0})",
			"syz_test$length1(&(0x7f0000001000)={0xff, 0x4})",
		},
		{
			"syz_test$length2(&(0x7f0000001000)={0xff, 0x0})",
			"syz_test$length2(&(0x7f0000001000)={0xff, 0x8})",
		},
		{
			"syz_test$length3(&(0x7f0000005000)={0xff, 0x0, 0x0})",
			"syz_test$length3(&(0x7f0000005000)={0xff, 0x4, 0x2})",
		},
		{
			"syz_test$length4(&(0x7f0000003000)={0x0, 0x0})",
			"syz_test$length4(&(0x7f0000003000)={0x2, 0x2})",
		},
		{
			"syz_test$length5(&(0x7f0000002000)={0xff, 0x0})",
			"syz_test$length5(&(0x7f0000002000)={0xff, 0x4})",
		},
		{
			"syz_test$length6(&(0x7f0000002000)={[0xff, 0xff, 0xff, 0xff], 0x0})",
			"syz_test$length6(&(0x7f0000002000)={[0xff, 0xff, 0xff, 0xff], 0x4})",
		},
		{
			"syz_test$length7(&(0x7f0000003000)={[0xff, 0xff, 0xff, 0xff], 0x0})",
			"syz_test$length7(&(0x7f0000003000)={[0xff, 0xff, 0xff, 0xff], 0x8})",
		},
		{
			"syz_test$length8(&(0x7f000001f000)={0x00, {0xff, 0x0, 0x00, [0xff, 0xff, 0xff]}, [{0xff, 0x0, 0x00, [0xff, 0xff, 0xff]}], 0x00, 0x0, [0xff, 0xff]})",
			"syz_test$length8(&(0x7f000001f000)={0x32, {0xff, 0x1, 0x10, [0xff, 0xff, 0xff]}, [{0xff, 0x1, 0x10, [0xff, 0xff, 0xff]}], 0x10, 0x1, [0xff, 0xff]})",
		},
		{
			"syz_test$length9(&(0x7f000001f000)={&(0x7f0000000000/0x5000)=nil, 0x0000})",
			"syz_test$length9(&(0x7f000001f000)={&(0x7f0000000000/0x5000)=nil, 0x5000})",
		},
		{
			"syz_test$length10(&(0x7f0000000000/0x5000)=nil, 0x0000, 0x0000, 0x0000, 0x0000)",
			"syz_test$length10(&(0x7f0000000000/0x5000)=nil, 0x5000, 0x5000, 0x2800, 0x1400)",
		},
		{
			"syz_test$length11(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x00)",
			"syz_test$length11(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x30)",
		},
		{
			"syz_test$length12(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x00)",
			"syz_test$length12(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x30)",
		},
		{
			"syz_test$length13(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x00)",
			"syz_test$length13(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x30)",
		},
		{
			"syz_test$length14(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x00)",
			"syz_test$length14(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x30)",
		},
		{
			"syz_test$length15(0xff, 0x0)",
			"syz_test$length15(0xff, 0x2)",
		},
		{
			"syz_test$length16(&(0x7f0000000000)={[0x42, 0x42], 0xff, 0xff, 0xff, 0xff, 0xff})",
			"syz_test$length16(&(0x7f0000000000)={[0x42, 0x42], 0x2, 0x10, 0x8, 0x4, 0x2})",
		},
		{
			"syz_test$length17(&(0x7f0000000000)={0x42, 0xff, 0xff, 0xff, 0xff})",
			"syz_test$length17(&(0x7f0000000000)={0x42, 0x8, 0x4, 0x2, 0x1})",
		},
		{
			"syz_test$length18(&(0x7f0000000000)={0x42, 0xff, 0xff, 0xff, 0xff})",
			"syz_test$length18(&(0x7f0000000000)={0x42, 0x8, 0x4, 0x2, 0x1})",
		},
		{
			"syz_test$length19(&(0x7f0000000000)={{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0xff}, 0xff, 0xff, 0xff})",
			"syz_test$length19(&(0x7f0000000000)={{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x14}, 0x14, 0x14, 0x5})",
		},
		{
			"syz_test$length20(&(0x7f0000000000)={{{0xff, 0xff, 0xff, 0xff}, 0xff, 0xff, 0xff}, 0xff, 0xff})",
			"syz_test$length20(&(0x7f0000000000)={{{0x4, 0x4, 0x7, 0x9}, 0x7, 0x7, 0x9}, 0x9, 0x9})",
		},
		{
			"syz_test$length21(&(0x7f0000000000)=0x0, 0x0)",
			"syz_test$length21(&(0x7f0000000000), 0x40)",
		},
		{
			"syz_test$length22(&(0x7f0000000000)='12345', 0x0)",
			"syz_test$length22(&(0x7f0000000000)='12345', 0x28)",
		},
		{
			"syz_test$length23(&(0x7f0000000000)={0x1, {0x2, 0x0}})",
			"syz_test$length23(&(0x7f0000000000)={0x1, {0x2, 0x6}})",
		},
		{
			"syz_test$length24(&(0x7f0000000000)={{0x0, {0x0}}, {0x0, {0x0}}})",
			"syz_test$length24(&(0x7f0000000000)={{0x0, {0x8}}, {0x0, {0x10}}})",
		},
		{
			"syz_test$length26(&(0x7f0000000000), 0x0)",
			"syz_test$length26(&(0x7f0000000000), 0x8)",
		},
		{
			"syz_test$length27(&(0x7f0000000000), 0x0)",
			"syz_test$length27(&(0x7f0000000000), 0x2a)",
		},
		{
			"syz_test$length28(&(0x7f0000000000), 0x0)",
			"syz_test$length28(&(0x7f0000000000), 0x2a)",
		},
		{
			"syz_test$length29(&(0x7f0000000000)={'./a\\x00', './b/c\\x00', 0x0, 0x0, 0x0})",
			"syz_test$length29(&(0x7f0000000000)={'./a\\x00', './b/c\\x00', 0xa, 0x14, 0x21})",
		},
	}

	for i, test := range tests {
		p, err := target.Deserialize([]byte(test.unsizedProg))
		if err != nil {
			t.Fatalf("failed to deserialize prog %v: %v", i, err)
		}
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
		p1 := strings.TrimSpace(string(p.Serialize()))
		if p1 != test.sizedProg {
			t.Fatalf("failed to assign sizes in prog %v\ngot  %v\nwant %v", i, p1, test.sizedProg)
		}
	}
}
