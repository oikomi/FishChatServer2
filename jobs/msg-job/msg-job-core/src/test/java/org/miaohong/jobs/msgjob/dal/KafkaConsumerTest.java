package org.miaohong.jobs.msgjob.dal;

import org.junit.Test;
import org.miaohong.jobs.msgjob.AbstractTest;
import org.miaohong.jobs.msgjob.dal.kafka.KafkaConsumer;

import javax.annotation.Resource;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by miaohong on 17/1/10.
 */
public class KafkaConsumerTest extends AbstractTest {
    @Resource
    KafkaConsumer kafkaConsumer;

    @Test
    public void testConsume() {
        List<String> topics = new ArrayList<>();
        topics.add("hello");
        kafkaConsumer.init(0,"111", topics);
        kafkaConsumer.consume();
    }
}
