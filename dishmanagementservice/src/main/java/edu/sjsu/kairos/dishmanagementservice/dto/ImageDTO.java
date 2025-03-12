package edu.sjsu.kairos.dishmanagementservice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.net.URL;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class ImageDTO {

	@JsonProperty("image_id")
	private UUID imageId;

	@JsonProperty("image_url")
	private URL imageurl;

	@JsonProperty("image_name")
	private String imageName;
}
