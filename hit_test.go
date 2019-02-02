package main

import (
	"log"
	"testing"
)

var result *ClientInfo

func BenchmarkBuildHit(b *testing.B) {
	mongo, err := ConnectMongo()
	if err != nil {
		log.Fatal("Unable to connect to MongoDB")
	}

	kafka, err := ConnectKafka()
	if err != nil {
		log.Fatal("Unable to connect to Kafka")
	}

	log.Println("[Message Broker] Kafka connected")

	defer func() {
		mongo.Close()

		if err := kafka.Close(); err != nil {
			panic(err)
		}
	}()

	ids := make(chan *ClientInfo)

	type args struct {
		event   []byte
		reqInfo *RequestInfo
		xcred   []byte
		ids     chan *ClientInfo
	}
	tests := &args{
		event:   []byte(`{"type":"click","screen":"https://www.santander.com.br//","lastScreen":"https://www.santander.com.br/","context":"/","category":"/","pageTitle":"Santander","session":{"referal":"","lat":"","lon":""},"user":{"uid":"9C91187FE9EC1E74583DDF093B3AFAFDF6EC204373850C375250B6A22B25BCD2"},"data":{"title":"Abra a sua conta","link":"/conversa/conta/passo2"},"channelId":"SWS8001","windowWidth":"379","windowHeight":"680","createdAt":"2018-09-02T18:04:58.692Z","deviceOtherIds":{"_gaid":"ASDFA111234"}}`),
		reqInfo: &RequestInfo{IP: "127.0.0.1", Channel: "TST8000", UserAgent: "TEST"},
		// xcred:   []byte(`eyJzaWQiOnsidmFsdWUiOiI1YjczN2IzMTc5OTkwNDdmZDYzMmZmYjEiLCJleHBpcmVzQXQiOiIyMDE4LTA4LTE0IDIyOjMwOjMzLjM5OTQyMzcyMiAtMDMwMCAtMDMifSwiZGlkIjp7InZhbHVlIjoiNWI3MzdiMzE3OTk5MDQ3ZmQ2MzJmZmIyIn19`),
		ids: ids,
	}

	b.Run("Bench", func(b *testing.B) {
		var res *ClientInfo
		for n := 0; n < b.N; n++ {
			// always record the result of res to prevent
			// the compiler eliminating the function call.
			go BuildHit(tests.event, tests.reqInfo, tests.xcred, tests.ids)
			res = <-ids
		}
		// always store the result to a package level variable
		// so the compiler cannot eliminate the Benchmark itself.
		result = res
	})
}
