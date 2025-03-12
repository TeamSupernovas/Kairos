package edu.sjsu.kairos.dishmanagementservice.config;

import edu.sjsu.kairos.dishmanagementservice.messaging.OrderEvent;
import lombok.extern.log4j.Log4j2;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.StringSerializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.*;
import org.springframework.kafka.support.serializer.JsonDeserializer;
import org.springframework.kafka.support.serializer.JsonSerializer;

import java.util.HashMap;
import java.util.Map;
import java.util.stream.Collectors;

@Configuration
@Log4j2
public class KafkaConfig {

    @Value("${kairos.dishmanagementservice.kafka.server}")
    private String kafkaServer;

   @Value("${kairos.dishmanagementservice.kafka.username}")
   private String kafkaUsername;

   @Value("${kairos.dishmanagementservice.kafka.password}")
   private String kafkaPassword;

    @Bean
    public ConsumerFactory<String, OrderEvent> orderEventConsumerFactory() {
        return createConsumerFactory(OrderEvent.class);
    }

    private <T> ConsumerFactory<String, T> createConsumerFactory(Class<T> eventType) {
        Map<String, Object> props = new HashMap<>();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, kafkaServer);
        props.put(ConsumerConfig.GROUP_ID_CONFIG, "order-group");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, JsonDeserializer.class);
        props.put("security.protocol", "SASL_SSL");
        props.put("sasl.mechanism", "PLAIN");
        String jaasTemplate = "org.apache.kafka.common.security.plain.PlainLoginModule required username=\'%s\' password=\'%s\';";
        String jaasCfg = String.format(jaasTemplate, kafkaUsername, kafkaPassword);
        props.put("sasl.jaas.config", jaasCfg);

        log.info("props::::::::::::::: {}" , props.values().stream().map(v -> v.toString()).collect(Collectors.joining("\n")));
        return new DefaultKafkaConsumerFactory<>(props, new StringDeserializer(), new JsonDeserializer<>(eventType));
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, OrderEvent> orderEventKafkaListenerContainerFactory() {
        return createKafkaListenerContainerFactory(orderEventConsumerFactory());
    }

    private <T> ConcurrentKafkaListenerContainerFactory<String, T> createKafkaListenerContainerFactory(ConsumerFactory<String, T> consumerFactory) {
        ConcurrentKafkaListenerContainerFactory<String, T> factory = new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        return factory;
    }

    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> configProps = new HashMap<>();
        configProps.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, kafkaServer);
        configProps.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        configProps.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JsonSerializer.class);
        configProps.put("security.protocol", "SASL_SSL");
        configProps.put("sasl.mechanism", "PLAIN");
        String jaasTemplate = "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"%s\" password=\"%s\";";
        String jaasCfg = String.format(jaasTemplate, kafkaUsername, kafkaPassword);
        configProps.put("sasl.jaas.config", jaasCfg);
        configProps.put(JsonSerializer.ADD_TYPE_INFO_HEADERS, false);

        return new DefaultKafkaProducerFactory<>(configProps);
    }

    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }
}
