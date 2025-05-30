package edu.sjsu.kairos.dishmanagementservice.config;

import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.regions.Regions;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;


@Configuration
public class AmazonS3Config {

    @Value("${cloud.aws.auth.credentials.accessKey}")
    private String accessKey;

    @Value("${cloud.aws.auth.credentials.secretKey}")
    private String secretKey;

    @Value("${cloud.aws.region}")
    private String region;

    @Bean
    public AmazonS3 amazonS3client() {
        BasicAWSCredentials awsCreds = new BasicAWSCredentials(accessKey, secretKey);
        return AmazonS3ClientBuilder.standard()
                .withRegion(Regions.fromName(region))
                .withCredentials(new AWSStaticCredentialsProvider(awsCreds))
                .build();
    }
}
