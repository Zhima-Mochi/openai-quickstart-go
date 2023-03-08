# openai-quickstart-go
## OpenAI API Quickstart - Golang example app  
This is an example pet name generator app used in the OpenAI API [quickstart tutorial](https://platform.openai.com/docs/quickstart). It uses the [Go](https://go.dev/) framework with [Gin](https://github.com/gin-gonic/gin). Check out the tutorial or follow the instructions below to get set up.

![Text box that says name my pet with an icon of a dog](https://user-images.githubusercontent.com/10623307/213887080-b2bc4645-7fdb-4dbd-ae42-efce00d0dc29.png)


## Setup

1. install go (version >= 1.18)

2. Clone this repository

3. Navigate into the project directory

   ```bash
   $ cd openai-quickstart-go
   ```

4. Install the requirements

   ```bash
   $ go mod tidy
   ```

5. Make a copy of the example environment variables file

   On Linux systems: 
   ```bash
   $ cp .env.example .env
   ```
   On Windows:
   ```powershell
   $ copy .env.example .env
   ```
6. Add your [API key](https://platform.openai.com/account/api-keys) to the newly created `.env` file

7. Run the app

   ```bash
   $ go run main.go
   ```

You should now be able to access the app at [http://localhost:8080](http://localhost:8080)! For the full context behind this example app, check out the [tutorial](https://platform.openai.com/docs/quickstart).