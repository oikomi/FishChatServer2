package org.miaohong.jobs.msgjob.dal.kafka;


import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Properties;


/**
 * Created by miaohong on 17/1/10.
 */
@Component
public class KafkaManager {
    private KafkaConsumer<String, String> consumer = null;
    private  List<String> topics = null;
    private  int id = 0;
    public KafkaManager(int id, String groupId,  List<String> topics) {
        this.id = id;
        this.topics = topics;
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("group.id", groupId);
        props.put("key.deserializer", StringDeserializer.class.getName());
        props.put("value.deserializer", StringDeserializer.class.getName());
        this.consumer = new KafkaConsumer<>(props);
    }

}
