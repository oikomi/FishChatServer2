package org.miaohong.jobs.msgjob.dal.kafka;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;

/**
 * Created by miaohong on 17/1/10.
 */
@Component
public class KafkaConsumer {
    private org.apache.kafka.clients.consumer.KafkaConsumer<String, String> consumer = null;
    private List<String> topics = null;
    private int id = 0;
    public void init(int id, String groupId,  List<String> topics) {
        this.id = id;
        this.topics = topics;
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("group.id", groupId);
        props.put("key.deserializer", StringDeserializer.class.getName());
        props.put("value.deserializer", StringDeserializer.class.getName());
        this.consumer = new org.apache.kafka.clients.consumer.KafkaConsumer<>(props);
        this.consumer.subscribe(topics);
    }

    public void consume() {
        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(Long.MAX_VALUE);
            for (ConsumerRecord<String, String> record : records) {
                Map<String, Object> data = new HashMap<>();
                data.put("partition", record.partition());
                data.put("offset", record.offset());
                data.put("value", record.value());
                System.out.println(this.id + ": " + data);
            }

        }
    }

}
