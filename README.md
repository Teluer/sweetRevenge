# Sweet Revenge Project #

## Backstory ##
- Ordered an item in an online shop
- Received cheap crap instead of the described product
- Called the shop and requested a refund
- Shop doesn't respond to calls ever since then
- Conclusion: it's a fraud and deserves some punishment

## Project description ##
- Run indefinitely
- Gather escort ladies' phone numbers, popular moldovan first and last names from the Internet
- Submit random fast orders in the online shop for random names and escort ladies' phones
- Generated orders should look legit and be indistinguishable from real orders
- Send orders using TOR to make it impossible to block the server IP
- Able to submit a specific order manually to test if everything is working
- Configurable order rates and schedule
- Keep track of all submitted orders in DB
- Keep logs

## Project purpose ##
- Make the shop operator call escort ladies all day long and ask if they ordered something from the shop.
- Cause lots of butthurt

## Used technologies ##
- Golang
- MySql database (with GORM library)
- TOR
- Docker

## Nice-to-have features ##
- Manual orders submission using RabbitMQ rather than DB table
- Admin panel to change some configs on the fly
- Unit testing