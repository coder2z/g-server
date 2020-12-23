/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 16:26
 **/
package oss

import (
	"errors"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/oss/alioss"
	"github.com/myxy99/component/oss/file"
	"github.com/myxy99/component/oss/standard"
)

func (i *ossInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range config.GetStringMap(prefix) {
		cfg := config.UnmarshalWithExpect(prefix+"."+name, newOssOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}

func (i *ossInvoker) new(o *options) (client standard.Oss) {
	var err error
	switch o.Mode {
	case "aliOss":
		client, err = alioss.NewOss(o.Addr, o.AccessKeyID, o.AccessKeySecret, o.OssBucket, o.IsDeleteSrcPath)
	case "file":
		client, err = file.NewOss(o.CdnName, o.FileBucket, o.IsDeleteSrcPath)
	default:
		err = errors.New("oss mode not exist")
	}
	if err != nil {
		panic(err)
	}
	return
}
