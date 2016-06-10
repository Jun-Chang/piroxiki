package piroxiki

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("conf.go", func() {
	It("LoadConf is parsed toml file correctly.", func() {
		cnf, err := LoadConf("./example.toml")
		fmt.Println(err)
		Ω(err).Should(BeNil())
		Ω(cnf).To(Equal(Conf{
			"mailToSlack": {
				In: InConf{
					Mail: &Mail{
						Server:   "imap@mail.com:993",
						User:     "user@xxx.xxx",
						Password: "xxxx",
						MailBox:  "",
					},
					Policy: Policy{
						Count:           10,
						IntervalSeconds: 100,
					},
					Filter: Filter{
						"subject": "subjects",
						"from":    "@system.info",
						"body":    "need to notification",
					},
				},
				Out: OutConf{
					Http: &Http{
						Method: "GET",
						URL:    "https://slack.com/api/chat.postMessage?aaa=bb",
					},
					Policy: Policy{
						Count:           10,
						IntervalSeconds: 100,
					},
				},
			},
		}))
	})
})
