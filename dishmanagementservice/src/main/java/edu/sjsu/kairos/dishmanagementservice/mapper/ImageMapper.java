package edu.sjsu.kairos.dishmanagementservice.mapper;

import edu.sjsu.kairos.dishmanagementservice.dto.ImageDTO;
import edu.sjsu.kairos.dishmanagementservice.model.Image;
import org.springframework.stereotype.Component;

@Component
public class ImageMapper {

    public ImageDTO toImageDTO(Image image) {
        return ImageDTO.builder()
                .imageId(image.getImageId())
                .imageName(image.getImageRef())
                .build();

    }

    public Image toImage(ImageDTO imageDTO) {
        return Image.builder()
                .imageId(imageDTO.getImageId())
                .imageRef(imageDTO.getImageName())
                .build();
    }

}
