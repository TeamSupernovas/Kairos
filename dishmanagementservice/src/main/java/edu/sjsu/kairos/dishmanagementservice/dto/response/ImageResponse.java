package edu.sjsu.kairos.dishmanagementservice.dto.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import edu.sjsu.kairos.dishmanagementservice.dto.ImageDTO;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class ImageResponse {

    @JsonProperty("images")
    private List<ImageDTO> images;
}
