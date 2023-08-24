package util_test

import (
	. "github.com/hanindo/util/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("FileNotExist",
	func(name string, x bool) {
		Expect(FileNotExist(name)).To(Equal(x))
	},
	Entry("not exist", "/proc/0", true),
	Entry("exist", "/proc/1/fd", false),
	Entry("permission denied", "/proc/1/fd/0", false),
)
