package org.miaohong.jobs.msgjob.manager;

import org.miaohong.jobs.msgjob.service.kafka.ConsumerService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;

/**
 * Created by miaohong on 17/1/14.
 */
@Component("startUpEvent")
public class InitManager implements ApplicationListener {
    private static boolean isInit = false;
    @Autowired
    ConsumerService consumerService;
    public void onApplicationEvent(ApplicationEvent applicationEvent) {
        if (!isInit) {
            consumerService.start();
            isInit = true;
        }
    }
}

