# iplocation-service
A simple service to get the location and external IP address of a client. The application downloads the latest GEO IP database from 
db-ip.com and uses it to resolve the location of a client. Every Monday the database is updated automatically.
## Example Response
```json
{
   "city_en":"Frankfurt am Main",
   "country_en":"Germany",
   "country_de":"Deutschland",
   "country_code":"DE",
   "continent_en":"Europe",
   "continent_de":"Europa",
   "longitude":"18.682130",
   "latitude":"53.31313",
   "subdivision":"Hesse",
   "ip":"222.229.169.30"
}
```

## Usage with Docker
The docker container is available on [Docker Hub](https://hub.docker.com/r/kickbeak/iplocation-service/).

### Run the container
```bash
docker run -d -p 3000:3000 kickbeak/iplocation-service
```