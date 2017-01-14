package org.miaohong.jobs.msgjob.service.kafka;

import org.miaohong.jobs.msgjob.dal.kafka.KafkaConsumer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

/**
 * Created by miaohong on 17/1/13.
 */
@Component
public class ConsumerService {
    @Autowired
    KafkaConsumer kafkaConsumer;
    private ExecutorService executor = null;

    public void start() {
        ExecutorService executor = Executors.newFixedThreadPool(1);
        executor.execute(kafkaConsumer);
    }

    public void stop() {
        try {
            System.out.println("attempt to shutdown executor");
            executor.shutdown();
            executor.awaitTermination(5, TimeUnit.SECONDS);
        }
        catch (InterruptedException e) {
            System.err.println("tasks interrupted");
        }
        finally {
            if (!executor.isTerminated()) {
                System.err.println("cancel non-finished tasks");
                executor.shutdownNow();
            }
            //executor.shutdownNow();
            System.out.println("shutdown finished");
        }
    }
}
