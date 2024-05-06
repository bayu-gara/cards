# Cards

The application built using Go 1.19 and consist 3 API's to simulate playing with a standard 52-card deck of French cards. It includes all thirteen ranks in each of the four suits: clubs, diamonds, hearts and spades.


# How to Run
Run this command in your terminal
```
make gorun
```

# API's

## Create a New Deck
Request
```
curl --location --request POST 'http://localhost:8080/v1/deck/create?shuffle=true&cards=AS%2C2H%2CKH%2CJC%2CQC'
```
Response
```json
{
	"deck_id": "b822befb-c760-4ed0-984a-399ab640df45",
	"shuffled": true,
	"remaining": 5
}
```

## Open a Deck

Request
```
curl --location 'http://localhost:8080/v1/deck/open?deck_id=b822befb-c760-4ed0-984a-399ab640df45'
```
Response
```json
{
	"deck_id": "b822befb-c760-4ed0-984a-399ab640df45",
	"shuffled": true,
	"remaining": 5,
	"cards": [
		{
			"code": "QC",
			"suit": "CLUBS",
			"value": "QUEEN"
		},
		{
			"code": "2H",
			"suit": "HEARTS",
			"value": "2"
		},
		{
			"code": "AS",
			"suit": "SPADES",
			"value": "ACE"
		},
		{
			"code": "JC",
			"suit": "CLUBS",
			"value": "JACK"
		},
		{
			"code": "KH",
			"suit": "HEARTS",
			"value": "KING"
		}
	]
}
```

## Draw a Card
Request
```
curl --location 'http://localhost:8080/v1/deck/draw?deck_id=b822befb-c760-4ed0-984a-399ab640df45&count=2'
```
Response
```json
{
	"cards": [
		{
			"code": "QC",
			"suit": "CLUBS",
			"value": "QUEEN"
		},
		{
			"code": "2H",
			"suit": "HEARTS",
			"value": "2"
		}
	]
}
```