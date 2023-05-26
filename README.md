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
- Send orders using TOR to make it impossible to block the server IP or trace a source
- Configurable order rates and schedule
- Keep track of all submitted orders in DB
- A server with an HTML Control Panel to submit orders manually and change some configs on the fly
- Admin panel to change some configs on the fly
- Keep logs

## Project purpose ##
- Make the shop operator call escort ladies all day long and ask if they ordered something from the shop.
- Cause lots of butthurt

## Project story ##
- First working instance was built within a week and sent several dozens of orders
- The web shop enabled captcha after about two days
- Ability to bypass captcha was implemented but couldn't be tested, because...
- The webshop disabled captcha. Either the admin assumed captcha didn't help, because some fake orders were sent manually, or the admin decided captcha would spoil user experience. It's possible the admin started checking sourse IP address and serving only those orders that originated in moldova.
- The project was rolled back to non-captcha logic, which is fortunate, because captcha solving requires a real browser with selenium and must use VPN instead of TOR.
- 

## Used technologies ##
- Golang (with gorm, logrus, testify)
- MySql database (with GORM library)
- RabbitMq (with connection pooling)
- TOR
- Docker

