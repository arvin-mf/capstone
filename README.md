After cloning this repository, make sure to create a `.env` file and copying the keys provided in `.env.example`. You can leave the `INFLUXDB_TOKEN` value empty for now.

Start the required services using Docker Compose.
> Make sure your Docker Engine is running before proceeding

Navigate to the project directory in your terminal and run the following command:
```
docker compose up
```

Once the services are running, open your browser and go to `localhost:8086` to log in to your InfluxDB instance. Copy your token and paste it into the `.env` file as the value of the `INFLUXDB_TOKEN=` key. 

Ensure that all services have started successfully before proceeding to the next steps.

Next, open the seeder file located at: `/internal/seeder/caps_seeder.go`.
Fill in the `deviceClients` array with the client IDs of the microcontrollers used in your area.

Then run the following command:
```
go run internal/seeder/caps_seeder.go
```
> Make sure Go is installed on your machine and properly configured in your system's environment variables

Finally, start your server with the following command:
```
go run cmd/main.go
```