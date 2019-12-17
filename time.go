package main

const blockIntervalMs = 1 * 1000

func msToIndex(ms int64) int64 {
	return (ms - 1) / blockIntervalMs
}

func msToNextIndex(ms int64) int64 {
	return msToIndex(ms + blockIntervalMs)
}

//
//func main() {
//	ticker := time.NewTicker(time.Duration(1 * time.Second))
//	go func() {
//		for {
//			for now := range ticker.C {
//				fmt.Println(msToNextIndex(now.UnixNano() / 1000000))
//			}
//		}
//	}()
//
//	select {}
//}
