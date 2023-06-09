# Sweet Revenge Project #

## Backstory ##
- Ordered an item in an online shop
- Received cheap crap instead of the described product
- Called the shop and requested a refund
- Shop doesn't answer my calls ever since then
- Conclusion: it's a fraud and deserves some punishment
- Fortunately, the shop doesn't have user registration or captcha, and it's possible to submit fast orders using customer name and phone only

## Project purpose ##
- Make the shop operator call escort ladies all day long and ask if they ordered something from the shop.
- Cause lots of butthurt

## Project description ##
- Run indefinitely
- Gather escort ladies' phone numbers, popular moldovan first and last names from the Internet
- Submit random fast orders in the online shop for random names and escort ladies' phones
- Generated orders should look legit and be indistinguishable from real orders (see [Order Concealing](#order-concealing) section)
- Send orders using TOR to make it impossible to block the server IP or trace a source
- Configurable order rates and schedule
- Keep track of all submitted orders in DB
- A server with an HTML Control Panel to submit orders manually and change some configs on the fly
- Admin panel to change some configs on the fly
- Keep logs

## Used technologies ##
- Golang (with gorm, gocron, goquery, logrus, testify etc)
- MySql database
- RabbitMq (with connection pooling)
- TOR
- Selenium
- Docker

## Order Concealing ##
- Send orders at random intervals
- Keep mean order rate low enough (around 30-60 minutes) to avoid causing panic
- Send orders from random IP addresses
- Use random user-agent header
- Use random phone and name formats
- Avoid reusing phones and names unless there are no new values available
- Gather cookies on each website call and attach them to the subsequent requests
- Visit success page
- For selenium flow, make random user actions while filling the order form

## Project story ##
- First working instance was built within a week and sent several dozens of orders.
- After order #45, the web shop enabled Google ReCaptcha. Ability to bypass captcha was partially implemented but couldn't be tested, because...
- The webshop disabled captcha a day later. Either the admin assumed captcha didn't help, because some fake orders were sent manually, or the admin decided captcha would spoil user experience.
- The project was rolled back to non-captcha logic, which is fortunate, because current ReCaptcha solution requires a captcha-solving browser extension and must use VPN instead of TOR.
- After order #165, a control order on a friend's number was submitted to check if the operator still calls, and it worked. This is puzzling, because it would be so easy to check source and ignore foreign IP addresses. Looks like human stupidity is the most powerful weapon after all.
- After order #235, the website enabled another (surprisingly easy) captcha. Selenium flow and captcha-solving logic were implemented a day later.
- After order #249, the website removed captcha again, because it didn't make any difference. The program didn't need any changes and continued working normally. Order interval temporally reduced to 4 minutes 30 seconds as a punishment.
- After order #362, it starts getting boring. Running out of phones again, adding new phone categories.
- After order #366, ReCaptcha came back. This one is hard to solve. Maybe it's time to wrap things up. Final project stage: improve order sending mechanism, wait for captcha to be removed (if this happens), fire as many orders as possible.
- Parallel order sending implemented. The flow successfully reaches recaptcha challenge (and then fails).
