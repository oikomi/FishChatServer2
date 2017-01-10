package org.miaohong.jobs.msgjob;

import org.junit.Before;
import org.junit.runner.RunWith;
import org.kubek2k.springockito.annotations.SpringockitoContextLoader;
import org.mockito.MockitoAnnotations;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.transaction.TransactionConfiguration;

/**
 * Created by miaohong on 17/1/10.
 */
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(loader = SpringockitoContextLoader.class,locations = {
        "classpath:META-INF/spring/application-service.xml"   })
@TransactionConfiguration(transactionManager = "transactionManager" ,defaultRollback = true)
public abstract class AbstractTest {
    @Before
    public void init() {
        MockitoAnnotations.initMocks(this);
    }
}
