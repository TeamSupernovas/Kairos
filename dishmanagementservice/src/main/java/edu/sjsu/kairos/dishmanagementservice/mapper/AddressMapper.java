package edu.sjsu.kairos.dishmanagementservice.mapper;

import edu.sjsu.kairos.dishmanagementservice.dto.AddressDTO;
import edu.sjsu.kairos.dishmanagementservice.model.Address;
import org.springframework.stereotype.Component;

@Component
public class AddressMapper {
    public Address toAddress(AddressDTO addressDTO) {
        return Address.builder()
                .street(addressDTO.getStreet())
                .city(addressDTO.getCity())
                .state(addressDTO.getState())
                .country(addressDTO.getCountry())
                .postalCode(addressDTO.getPostalCode())
                .country(addressDTO.getCountry())
                .build();
    }

    public AddressDTO toAddressDTO(Address address) {
        return AddressDTO.builder()
                .street(address.getStreet())
                .city(address.getCity())
                .state(address.getState())
                .country(address.getCountry())
                .postalCode(address.getPostalCode())
                .country(address.getCountry())
                .build();
    }
}
