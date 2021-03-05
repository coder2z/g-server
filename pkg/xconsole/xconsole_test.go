/**
 * @Author: yangon
 * @Description
 * @Date: 2021/3/5 17:26
 **/
package xconsole

import "testing"

func TestBlue(t *testing.T) {
	ResetDebug(true)
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"1", args{
			"test",
		},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Blue(tt.args.msg)
			Yellow(tt.args.msg)
			Red(tt.args.msg)
			Green(tt.args.msg)

			Bluef(tt.args.msg, "data")
			Yellowf(tt.args.msg, "data")
			Redf(tt.args.msg, "data")
			Greenf(tt.args.msg, "data")

		})
	}
}
