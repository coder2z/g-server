/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package oss

type options struct {
	Debug           bool   `mapStructure:"debug"`
	Mode            string `mapStructure:"mode"`
	Addr            string `mapStructure:"addr"`
	AccessKeyID     string `mapStructure:"accessKeyId"`
	AccessKeySecret string `mapStructure:"accessKeySecret"`
	CdnName         string `mapStructure:"cdnName"`
	OssBucket       string `mapStructure:"ossBucket"`
	FileBucket      string `mapStructure:"fileBucket"`
	IsDeleteSrcPath bool   `mapStructure:"isDeleteSrcPath"`
}

func newOssOptions() *options {
	return &options{
		Debug:           false,
		Mode:            "file",
		Addr:            "",
		AccessKeyID:     "",
		AccessKeySecret: "",
		CdnName:         "",
		OssBucket:       "",
		FileBucket:      ".",
		IsDeleteSrcPath: false,
	}
}
