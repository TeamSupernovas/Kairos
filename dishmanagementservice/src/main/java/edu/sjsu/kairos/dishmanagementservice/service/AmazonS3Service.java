package edu.sjsu.kairos.dishmanagementservice.service;

import com.amazonaws.HttpMethod;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.model.DeleteObjectRequest;
import com.amazonaws.services.s3.model.GeneratePresignedUrlRequest;
import com.amazonaws.services.s3.model.ObjectMetadata;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.net.URL;
import java.util.Date;

@Service
public class AmazonS3Service {

    @Autowired
    private AmazonS3 amazonS3Client;

    @Value("${cloud.aws.s3.bucketname}")
    private String bucketName;

    public String uploadFile(MultipartFile file) throws IOException {
        String fileName = System.currentTimeMillis() + "_" + file.getOriginalFilename();
        amazonS3Client.putObject(bucketName, fileName, file.getInputStream(), new ObjectMetadata());
        return fileName ;
    }
    
    public URL generatePresignedUrl(String fileName, int expirationInMinutes) {
        Date expiration = new Date(System.currentTimeMillis() + expirationInMinutes * 60 * 1000);
        GeneratePresignedUrlRequest generatePresignedUrlRequest =
                new GeneratePresignedUrlRequest(bucketName, fileName)
                        .withMethod(HttpMethod.GET)
                        .withExpiration(expiration);

        return amazonS3Client.generatePresignedUrl(generatePresignedUrlRequest);
    }
    
    public void deleteFile(String fileName) {
    	DeleteObjectRequest deleteObjectRequest = new DeleteObjectRequest(bucketName, fileName);
    	amazonS3Client.deleteObject(deleteObjectRequest);
    }
}
