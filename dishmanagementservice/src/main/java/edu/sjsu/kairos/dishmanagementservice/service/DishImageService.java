package edu.sjsu.kairos.dishmanagementservice.service;

import java.io.IOException;
import java.net.URL;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import edu.sjsu.kairos.dishmanagementservice.dto.ImageDTO;
import edu.sjsu.kairos.dishmanagementservice.model.Image;
import edu.sjsu.kairos.dishmanagementservice.repository.ImageRepository;
import edu.sjsu.kairos.dishmanagementservice.util.DishManagementServiceConstants;
import org.apache.kafka.common.errors.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;


@Service
public class DishImageService {

    @Autowired
    private AmazonS3Service amazonS3Service;

    @Autowired
    private ImageRepository imageRepository;

    @Autowired
    private MessageService messageService;

    public Image createImage(MultipartFile multipartFile) throws IOException {
        String filename = amazonS3Service.uploadFile(multipartFile);
        return Image.builder()
                .imageRef(filename).build();
    }

    public ImageDTO getImageByImageId(UUID imageId) throws ResourceNotFoundException {
        return imageRepository.findById(imageId)
                .map(image -> {
                    URL imageURL = amazonS3Service.getPublicImageUrl(image.getImageRef());
                    return ImageDTO.builder().imageurl(imageURL).build();
                })
                .orElseThrow(() -> new ResourceNotFoundException(
                        messageService.getMessage(DishManagementServiceConstants.IMAGE_NOT_FOUND_ERROR_KEY, imageId)));
    }

    public List<ImageDTO> getImagesByDishId(UUID dishId) {
        List<Image> images = imageRepository.findByDish_DishIdAndDeletedAtIsNull(dishId);

        return images.stream()
                .map(image -> {
                    URL imageURL = amazonS3Service.getPublicImageUrl(image.getImageRef());
                    return ImageDTO.builder()
                            .imageId(image.getImageId())
                            .imageurl(imageURL)
                            .imageName(image.getImageRef())
                            .build();
                })
                .collect(Collectors.toList());
    }

    public void deleteImageByImageId(UUID imageId) {
        imageRepository.findById(imageId)
                .ifPresentOrElse(image -> {
                    image.setDeletedAt(LocalDateTime.now());
                    amazonS3Service.deleteFile(image.getImageRef());
                }, () -> new ResourceNotFoundException(
                        messageService.getMessage(DishManagementServiceConstants.IMAGE_NOT_FOUND_ERROR_KEY, imageId)));
    }

}
