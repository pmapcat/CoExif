package main

import (
	"github.com/jmcvetta/napping"
	"github.com/satori/go.uuid"
	. "gopkg.in/check.v1"
	"log"
	"net/url"
	"os"
	"path"
	"testing"
)

const PORT = "9999"

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	R *RESTHandler
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *C) {
	// start server in goroutine
	s.R = &RESTHandler{}
	go s.R.Run(PORT, "admin", "admin", "/")
}

func (s *TestSuite) TearDownSuite(c *C) {
	// delete all exiftool processes
	s.R.Dispatch.Exit()
}

func (s *TestSuite) TestAllRight(c *C) {
	BASE_URL := "http://127.0.0.1:9999"
	cur_dir, err := os.Getwd()
	if err != nil {
		c.Error(err)
	}
	test_file := path.Join(cur_dir, "temp/a.png")
	finite_url := BASE_URL + test_file
	log.Println(finite_url)

	// ============ authenticate ==================
	session := napping.Session{
		Userinfo: url.UserPassword("admin", "admin"),
	}
	//  ==========================================
	for i := 1; i <= 10; i++ {
		log.Println("Working ", i)
		result := IdioticJSON{}
		payload := StandartJSON{}
		payload["Artist"] = uuid.NewV4().String()
		payload["Author"] = uuid.NewV4().String()
		payload["Comment"] = uuid.NewV4().String()
		payload["Copyright"] = uuid.NewV4().String()
		// wait for server to be up ..
		_, err = session.Post(finite_url, &payload, &result, nil)
		for err != nil {
			_, err := session.Post(finite_url, &payload, &result, nil)
			if err == nil {
				break
			}
		}
		c.Assert(result.Items[0]["Artist"], Equals, payload["Artist"])
		c.Assert(result.Items[0]["Author"], Equals, payload["Author"])
		c.Assert(result.Items[0]["Comment"], Equals, payload["Comment"])
		c.Assert(result.Items[0]["Copyright"], Equals, payload["Copyright"])
	}
}
