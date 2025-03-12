package edu.sjsu.kairos.dishmanagementservice.mapper;

import edu.sjsu.kairos.dishmanagementservice.dto.LocationDTO;
import edu.sjsu.kairos.dishmanagementservice.model.Location;
import org.springframework.stereotype.Component;

@Component
public class LocationMapper {

    public LocationDTO toLocationDTO(Location location) {
        return LocationDTO.builder()
                .latitude(location.getLatitude())
                .longitude(location.getLongitude())
                .build();
    }

    public Location toLocation(LocationDTO locationDTO) {
        return Location.builder()
                .latitude(locationDTO.getLatitude())
                .longitude(locationDTO.getLongitude())
                .build();
    }
}
