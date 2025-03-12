package edu.sjsu.kairos.dishmanagementservice.service;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
@Log4j2
public class GeoLocationService {

    @Value("${google.maps.api.key}")
    private String googleApiKey;

    @Autowired
    private RestTemplate restTemplate;

    @Autowired
    private Environment environment;


    public double[] getLatLongFromAddress(String address) {
        try {
            if (address == null || address.trim().isEmpty()) {
                log.error("Address cannot be null or empty");
                throw new IllegalArgumentException("Address cannot be null or empty");
            }

            String normalizedAddress = address.replace(" ", "_").replace(",", "");
            log.info("normalizedAddress : {}", normalizedAddress);
            String cachedCoordinates = environment.getProperty("geo.cache." + normalizedAddress);
            log.info("cachedCoordinates : {}", cachedCoordinates);
            if (cachedCoordinates != null) {
                log.info("Fetching cached coordinates for address - {}", address);
                String[] parts = cachedCoordinates.split(",");
                return new double[]{Double.parseDouble(parts[0]), Double.parseDouble(parts[1])};
            }

            String url = "https://maps.googleapis.com/maps/api/geocode/json?address=" +
                    address.replace(" ", "+") + "&key=" + googleApiKey;

            log.info("Calling Geocoding API with address : {}", address);
            String response = restTemplate.getForObject(url, String.class);

            ObjectMapper objectMapper = new ObjectMapper();
            JsonNode jsonNode = objectMapper.readTree(response);

            JsonNode results = jsonNode.path("results");
            if (results.isEmpty()) {
                log.error("No results found for address: {}", address);
                throw new RuntimeException("No results found for address: ");
            }

            JsonNode location = results.get(0).path("geometry").path("location");
            double latitude = location.path("lat").asDouble();
            double longitude = location.path("lng").asDouble();

            return new double[]{latitude, longitude};
        } catch (Exception e) {
            log.error("Failed to get coordinates for address: {}", address);
            throw new RuntimeException("Failed to get coordinates for address", e);
        }
    }
}
