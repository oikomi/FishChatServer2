package org.miaohong.jobs.msgjob.dal;

import org.junit.Test;
import org.miaohong.jobs.msgjob.AbstractTest;
import org.miaohong.jobs.msgjob.service.kafka.ConsumerService;
import org.springframework.beans.factory.annotation.Autowired;


/**
 * Created by miaohong on 17/1/10.
 */
public class KafkaConsumerTest extends AbstractTest {
    @Autowired
    ConsumerService consumerService;

    @Test
    public void testConsume() {
        consumerService.start();
    }
}
