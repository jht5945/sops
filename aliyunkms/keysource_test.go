package aliyunkms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseRegionIdFromArn(t *testing.T) {
	regionId := parseRegionIdFromArn("")
	assert.EqualValues(t, "", regionId)

	regionId2 := parseRegionIdFromArn("acs:kms:cn-hangzhou:1192853035118460:key/key-hzz64a3dbd1prbfsnnvpe")
	assert.EqualValues(t, "cn-hangzhou", regionId2)
}
