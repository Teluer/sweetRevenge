package main

import (
	"sweetRevenge/fetcher"
)

/*TODO:
check data in DB: select least used phone and least used name
if no luck with phones, go to 999 and update
if no luck with names, go to names sources and update

connect to TOR via Socks
go to the shop and make an order

set timer to sleep
*/

func main() {
	//fetcher.GetLadiesPhones()
	//db.Connect()

	fetcher.UpdateLastNames()
	fetcher.UpdateFirstNames()
}
